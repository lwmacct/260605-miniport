package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260605-miniport/internal/service"
)

type authHandler struct {
	config   Config
	services Services
}

func RegisterAuth(api huma.API, config Config, services Services) {
	handler := authHandler{config: config, services: services}
	auth := huma.NewGroup(api, "/auth")
	huma.Register(auth, huma.Operation{OperationID: "get-auth-config", Method: http.MethodGet, Path: "/config", Tags: []string{"Auth"}}, handler.configOutput)
	huma.Register(auth, huma.Operation{OperationID: "get-auth-state", Method: http.MethodGet, Path: "/state", Tags: []string{"Auth"}}, handler.state)
	huma.Register(auth, huma.Operation{OperationID: "create-auth-challenge", Method: http.MethodPost, Path: "/challenges", Tags: []string{"Auth"}}, handler.createChallenge)
	huma.Register(auth, huma.Operation{OperationID: "register-password-user", Method: http.MethodPost, Path: "/password/register", DefaultStatus: http.StatusCreated, Tags: []string{"Auth"}}, handler.passwordRegister)
	huma.Register(auth, huma.Operation{OperationID: "login-password", Method: http.MethodPost, Path: "/password/login", Tags: []string{"Auth"}}, handler.passwordLogin)
	huma.Register(auth, huma.Operation{OperationID: "logout", Method: http.MethodPost, Path: "/logout", Tags: []string{"Auth"}}, handler.logout)
	huma.Register(auth, huma.Operation{OperationID: "get-current-user", Method: http.MethodGet, Path: "/me", Tags: []string{"Auth"}}, handler.me)
}

func (h authHandler) configOutput(_ context.Context, _ *struct{}) (*BodyDTO[AuthConfigDTO], error) {
	return &BodyDTO[AuthConfigDTO]{Body: h.configBody()}, nil
}

func (h authHandler) configBody() AuthConfigDTO {
	challenge := h.services.Challenges.PublicConfig()
	body := AuthConfigDTO{}
	body.Local.LoginEnabled = h.config.LocalLoginEnabled
	body.Local.RegistrationEnabled = h.config.LocalRegistrationEnabled
	body.Challenge.Provider = challenge.Provider
	body.Challenge.SiteKey = challenge.SiteKey
	return body
}

func (h authHandler) state(ctx context.Context, input *AuthSessionInputDTO) (*BodyDTO[AuthStateDTO], error) {
	body := AuthStateDTO{
		Config:  h.configBody(),
		Session: AuthSessionDTO{Authenticated: false},
	}
	if input.Session == "" {
		return &BodyDTO[AuthStateDTO]{Body: body}, nil
	}
	session, err := h.session(ctx, input.Session)
	if err != nil {
		return &BodyDTO[AuthStateDTO]{Body: body}, nil //nolint:nilerr // Invalid sessions are represented as unauthenticated state.
	}
	body.Session = *session
	return &BodyDTO[AuthStateDTO]{Body: body}, nil
}

func (h authHandler) createChallenge(ctx context.Context, _ *struct{}) (*BodyDTO[AuthChallengeCreateDTO], error) {
	request, err := utilRequest(ctx, h.config)
	if err != nil {
		return nil, err
	}
	challenge, err := h.services.Challenges.Create(ctx, ToAuthChallengeInput(request))
	if err != nil {
		if errors.Is(err, service.ErrAuthChallengeLimitExceeded) {
			return nil, huma.Error429TooManyRequests("too many challenges")
		}
		return nil, huma.Error400BadRequest("challenge creation unsupported")
	}
	return &BodyDTO[AuthChallengeCreateDTO]{Body: ToAuthChallengeCreateDTO(challenge)}, nil
}

func (h authHandler) passwordRegister(ctx context.Context, input *BodyInputDTO[AuthCredentialsDTO]) (*AuthSessionResponseDTO, error) {
	if !h.config.LocalLoginEnabled || !h.config.LocalRegistrationEnabled {
		return nil, huma.Error403Forbidden("password registration disabled")
	}
	request, err := utilRequest(ctx, h.config)
	if err != nil {
		return nil, err
	}
	challengeErr := h.verifyChallenge(ctx, input.Body.Challenge, request)
	if challengeErr != nil {
		return nil, huma.Error401Unauthorized("invalid challenge")
	}
	user, err := h.services.Passwords.Register(ctx, service.AuthPasswordRegisterInput{Username: input.Body.Username, Password: input.Body.Password})
	if err != nil {
		if errors.Is(err, service.ErrIdentityUserUsernameTaken) || errors.Is(err, service.ErrAuthPasswordWeakPassword) || errors.Is(err, service.ErrIdentityUserInvalidCredentials) {
			return nil, huma.Error400BadRequest(utilRegisterErrorMessage(err))
		}
		return nil, huma.Error500InternalServerError("register failed")
	}
	return h.createSessionResponse(ctx, user.ID, request)
}

func (h authHandler) passwordLogin(ctx context.Context, input *BodyInputDTO[AuthCredentialsDTO]) (*AuthSessionResponseDTO, error) {
	if !h.config.LocalLoginEnabled {
		return nil, huma.Error403Forbidden("password login disabled")
	}
	request, err := utilRequest(ctx, h.config)
	if err != nil {
		return nil, err
	}
	challengeErr := h.verifyChallenge(ctx, input.Body.Challenge, request)
	if challengeErr != nil {
		return nil, huma.Error401Unauthorized("invalid challenge")
	}
	user, err := h.services.Passwords.Authenticate(ctx, input.Body.Username, input.Body.Password, h.services.Users)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	return h.createSessionResponse(ctx, user.ID, request)
}

func (h authHandler) logout(ctx context.Context, input *AuthSessionInputDTO) (*AuthSessionResponseDTO, error) {
	if input.Session != "" {
		_ = h.services.Sessions.Delete(ctx, input.Session)
	}
	return &AuthSessionResponseDTO{
		SetCookie: utilClearSessionCookie(h.config.SecureCookies),
		Body:      AuthSessionDTO{Authenticated: false},
	}, nil
}

func (h authHandler) me(ctx context.Context, input *AuthSessionInputDTO) (*BodyDTO[AuthSessionDTO], error) {
	session, err := h.session(ctx, input.Session)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	return &BodyDTO[AuthSessionDTO]{Body: *session}, nil
}

func (h authHandler) session(ctx context.Context, sessionID string) (*AuthSessionDTO, error) {
	if sessionID == "" {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	request, err := utilRequest(ctx, h.config)
	if err != nil {
		return nil, err
	}
	if request.InvalidSessionCookie {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	sessionUser, err := h.services.Sessions.User(ctx, sessionID, request, h.services.Users)
	if err != nil {
		return nil, err
	}
	user, err := h.services.Users.ByID(ctx, sessionUser.ID)
	if err != nil {
		return nil, err
	}
	out := ToAuthSessionDTO(sessionUser, user, utilIsRuntimeAdmin(user, h.config.RuntimeAdmins))
	return &out, nil
}

func (h authHandler) verifyChallenge(ctx context.Context, challenge AuthChallengeDTO, request service.AuthSessionInput) error {
	return h.services.Challenges.Verify(ctx, ToAuthChallengeAnswer(challenge), ToAuthChallengeInput(request))
}

func (h authHandler) createSessionResponse(ctx context.Context, userID int64, request service.AuthSessionInput) (*AuthSessionResponseDTO, error) {
	sessionID, expiresAt, err := h.services.Sessions.Create(ctx, userID, request)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}
	user, err := h.services.Users.ByID(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}
	return &AuthSessionResponseDTO{
		SetCookie: utilSessionCookieValue(sessionID, expiresAt, h.config.SecureCookies),
		Body: AuthSessionDTO{
			Authenticated: true,
			ExpiresAt:     utilHTTPTime(expiresAt),
			User:          ToAuthUserDTO(user, utilIsRuntimeAdmin(user, h.config.RuntimeAdmins)),
		},
	}, nil
}
