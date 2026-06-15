package server

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/uptrace/bun"

	"github.com/lwmacct/260605-miniport/internal/adapter/httpauth"
	"github.com/lwmacct/260605-miniport/internal/domain/authchallenge"
	"github.com/lwmacct/260605-miniport/internal/domain/authpassword"
	"github.com/lwmacct/260605-miniport/internal/domain/authsession"
	"github.com/lwmacct/260605-miniport/internal/domain/identityuser"
)

type authConfigDTO struct {
	Local struct {
		LoginEnabled        bool `json:"loginEnabled"`
		RegistrationEnabled bool `json:"registrationEnabled"`
	} `json:"local"`
	Challenge struct {
		Provider string `json:"provider"`
		SiteKey  string `json:"sitekey,omitempty"`
	} `json:"challenge"`
}

type authChallengeCreateDTO struct {
	Provider    string    `json:"provider"`
	ChallengeID string    `json:"challengeId,omitempty"`
	Image       string    `json:"image,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
}

type authChallengeDTO struct {
	Provider    string `json:"provider"`
	ChallengeID string `json:"challengeId,omitempty"`
	Answer      string `json:"answer,omitempty"`
	Token       string `json:"token,omitempty"`
}

type authCredentialsDTO struct {
	Username  string           `json:"username"`
	Password  string           `json:"password"`
	Challenge authChallengeDTO `json:"challenge"`
}

type authUserDTO struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`
	Admin       bool   `json:"admin"`
}

type authSessionDTO struct {
	Authenticated bool         `json:"authenticated"`
	ExpiresAt     string       `json:"expiresAt,omitempty"`
	User          *authUserDTO `json:"user,omitempty"`
}

type authBody[T any] struct {
	Body T `json:"body"`
}

type authBodyInput[T any] struct {
	Body T
}

type authSessionInput struct {
	Session string `cookie:"web_session"`
}

type authSessionResponse struct {
	SetCookie string `header:"Set-Cookie"`
	Body      authSessionDTO
}

type adminUserDTO struct {
	ID          int64   `json:"id"`
	Username    string  `json:"username"`
	DisplayName string  `json:"displayName"`
	Status      string  `json:"status"`
	Admin       bool    `json:"admin"`
	DisabledAt  *string `json:"disabledAt,omitempty"`
}

func (app *App) registerAuthHTTPAPI(api huma.API) {
	auth := huma.NewGroup(api, "/auth")
	huma.Register(auth, huma.Operation{
		OperationID: "get-auth-config",
		Method:      http.MethodGet,
		Path:        "/config",
		Tags:        []string{"Auth"},
	}, app.authConfig)
	huma.Register(auth, huma.Operation{
		OperationID: "create-auth-challenge",
		Method:      http.MethodPost,
		Path:        "/challenges",
		Tags:        []string{"Auth"},
	}, app.authCreateChallenge)
	huma.Register(auth, huma.Operation{
		OperationID:   "register-password-user",
		Method:        http.MethodPost,
		Path:          "/password/register",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Auth"},
	}, app.authPasswordRegister)
	huma.Register(auth, huma.Operation{
		OperationID: "login-password",
		Method:      http.MethodPost,
		Path:        "/password/login",
		Tags:        []string{"Auth"},
	}, app.authPasswordLogin)
	huma.Register(auth, huma.Operation{
		OperationID: "logout",
		Method:      http.MethodPost,
		Path:        "/logout",
		Tags:        []string{"Auth"},
	}, app.authLogout)
	huma.Register(auth, huma.Operation{
		OperationID: "get-current-user",
		Method:      http.MethodGet,
		Path:        "/me",
		Tags:        []string{"Auth"},
	}, app.authMe)

	admin := huma.NewGroup(api, "/admin")
	huma.Register(admin, huma.Operation{
		OperationID: "list-admin-users",
		Method:      http.MethodGet,
		Path:        "/users",
		Tags:        []string{"Admin"},
	}, app.adminListUsers)
}

func (app *App) authConfig(_ context.Context, _ *struct{}) (*authBody[authConfigDTO], error) {
	challenge := app.challenges.PublicConfig()
	body := authConfigDTO{}
	body.Local.LoginEnabled = app.cfg.Server.Auth.Local.LoginEnabled
	body.Local.RegistrationEnabled = app.cfg.Server.Auth.Local.RegistrationEnabled
	body.Challenge.Provider = challenge.Provider
	body.Challenge.SiteKey = challenge.SiteKey
	return &authBody[authConfigDTO]{Body: body}, nil
}

func (app *App) authCreateChallenge(ctx context.Context, _ *struct{}) (*authBody[authChallengeCreateDTO], error) {
	request, ok := httpauth.RequestFromContext(ctx)
	if !ok {
		return nil, huma.Error400BadRequest("invalid request source")
	}
	challenge, err := app.challenges.Create(ctx, toChallengeRequest(request))
	if err != nil {
		if errors.Is(err, authchallenge.ErrLimitExceeded) {
			return nil, huma.Error429TooManyRequests("too many challenges")
		}
		return nil, huma.Error400BadRequest("challenge creation unsupported")
	}
	return &authBody[authChallengeCreateDTO]{Body: authChallengeCreateDTO{
		Provider:    challenge.Provider,
		ChallengeID: challenge.ChallengeID,
		Image:       challenge.Image,
		ExpiresAt:   challenge.ExpiresAt,
	}}, nil
}

func (app *App) authPasswordRegister(ctx context.Context, input *authBodyInput[authCredentialsDTO]) (*authSessionResponse, error) {
	if !app.cfg.Server.Auth.Local.LoginEnabled || !app.cfg.Server.Auth.Local.RegistrationEnabled {
		return nil, huma.Error403Forbidden("password registration disabled")
	}
	request, ok := httpauth.RequestFromContext(ctx)
	if !ok {
		return nil, huma.Error400BadRequest("invalid request source")
	}
	if err := app.verifyChallenge(ctx, input.Body.Challenge, request); err != nil {
		return nil, huma.Error401Unauthorized("invalid challenge")
	}
	if err := app.passwords.CheckStrength(input.Body.Username, input.Body.Password); err != nil {
		return nil, huma.Error400BadRequest("weak password")
	}

	var user *identityuser.User
	err := app.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		users := identityuser.NewService(tx)
		passwords := authpassword.NewService(tx)
		created, err := users.Create(ctx, identityuser.CreateInput{Username: input.Body.Username})
		if err != nil {
			return err
		}
		if err := passwords.Set(ctx, created.Username, created.ID, input.Body.Password); err != nil {
			return err
		}
		user = created
		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, authpassword.ErrWeakPassword):
			return nil, huma.Error400BadRequest("weak password")
		case errors.Is(err, identityuser.ErrUsernameTaken):
			return nil, huma.Error400BadRequest("username taken")
		case errors.Is(err, identityuser.ErrInvalidCredentials):
			return nil, huma.Error400BadRequest("invalid credentials")
		default:
			return nil, huma.Error500InternalServerError("register failed")
		}
	}
	return app.createSessionResponse(ctx, user.ID, request)
}

func (app *App) authPasswordLogin(ctx context.Context, input *authBodyInput[authCredentialsDTO]) (*authSessionResponse, error) {
	if !app.cfg.Server.Auth.Local.LoginEnabled {
		return nil, huma.Error403Forbidden("password login disabled")
	}
	request, ok := httpauth.RequestFromContext(ctx)
	if !ok {
		return nil, huma.Error400BadRequest("invalid request source")
	}
	if err := app.verifyChallenge(ctx, input.Body.Challenge, request); err != nil {
		return nil, huma.Error401Unauthorized("invalid challenge")
	}
	user, err := app.passwords.Authenticate(ctx, input.Body.Username, input.Body.Password, app.users)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	return app.createSessionResponse(ctx, user.ID, request)
}

func (app *App) authLogout(ctx context.Context, input *authSessionInput) (*authSessionResponse, error) {
	if input.Session != "" {
		_ = app.sessions.Delete(ctx, input.Session)
	}
	return &authSessionResponse{
		SetCookie: httpauth.ClearSessionCookie(app.httpAuth.SecureCookies()).String(),
		Body: authSessionDTO{
			Authenticated: false,
		},
	}, nil
}

func (app *App) authMe(ctx context.Context, input *authSessionInput) (*authBody[authSessionDTO], error) {
	user, sessionUser, err := app.currentSessionUser(ctx, input.Session)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	return &authBody[authSessionDTO]{
		Body: authSessionDTO{
			Authenticated: true,
			ExpiresAt:     sessionUser.ExpiresAt.UTC().Format(http.TimeFormat),
			User:          app.toAuthUserDTO(user),
		},
	}, nil
}

func (app *App) adminListUsers(ctx context.Context, input *authSessionInput) (*authBody[[]adminUserDTO], error) {
	current, _, err := app.currentSessionUser(ctx, input.Session)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	if !app.isRuntimeAdmin(current) {
		return nil, huma.Error403Forbidden("forbidden")
	}

	users, err := app.users.List(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}
	body := make([]adminUserDTO, 0, len(users))
	for _, item := range users {
		var disabledAt *string
		if item.DisabledAt != nil {
			value := item.DisabledAt.UTC().Format(http.TimeFormat)
			disabledAt = &value
		}
		body = append(body, adminUserDTO{
			ID:          item.ID,
			Username:    item.Username,
			DisplayName: item.DisplayName,
			Status:      item.Status,
			Admin:       app.isRuntimeAdminUsername(item.Username),
			DisabledAt:  disabledAt,
		})
	}
	return &authBody[[]adminUserDTO]{Body: body}, nil
}

func (app *App) verifyChallenge(ctx context.Context, challenge authChallengeDTO, request authsession.Request) error {
	return app.challenges.Verify(ctx, authchallenge.ResponseDTO{
		Provider:    challenge.Provider,
		ChallengeID: challenge.ChallengeID,
		Answer:      challenge.Answer,
		Token:       challenge.Token,
	}, toChallengeRequest(request))
}

func toChallengeRequest(request authsession.Request) authchallenge.RequestDTO {
	return authchallenge.RequestDTO{
		IP:         request.IP,
		UserAgent:  request.UserAgent,
		Method:     request.Method,
		Path:       request.Path,
		RemoteAddr: request.RemoteAddr,
	}
}

func (app *App) createSessionResponse(ctx context.Context, userID int64, request authsession.Request) (*authSessionResponse, error) {
	sessionID, expiresAt, err := app.sessions.Create(ctx, userID, request)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}
	user, err := app.users.ByID(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}
	return &authSessionResponse{
		SetCookie: httpauth.SessionCookieValue(sessionID, expiresAt, app.httpAuth.SecureCookies()).String(),
		Body: authSessionDTO{
			Authenticated: true,
			ExpiresAt:     expiresAt.UTC().Format(http.TimeFormat),
			User:          app.toAuthUserDTO(user),
		},
	}, nil
}

func (app *App) currentSessionUser(ctx context.Context, sessionID string) (*identityuser.User, authsession.UserDTO, error) {
	if sessionID == "" {
		return nil, authsession.UserDTO{}, errors.New("missing session")
	}
	request, ok := httpauth.RequestFromContext(ctx)
	if !ok {
		return nil, authsession.UserDTO{}, errors.New("invalid request source")
	}
	sessionUser, err := app.sessions.User(ctx, sessionID, request, app.users)
	if err != nil {
		return nil, authsession.UserDTO{}, err
	}
	user, err := app.users.ByID(ctx, sessionUser.ID)
	if err != nil {
		return nil, authsession.UserDTO{}, err
	}
	return user, *sessionUser, nil
}

func (app *App) toAuthUserDTO(user *identityuser.User) *authUserDTO {
	if user == nil {
		return nil
	}
	return &authUserDTO{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Status:      user.Status,
		Admin:       app.isRuntimeAdmin(user),
	}
}

func (app *App) isRuntimeAdmin(user *identityuser.User) bool {
	if user == nil {
		return false
	}
	return app.isRuntimeAdminUsername(user.Username)
}

func (app *App) isRuntimeAdminUsername(username string) bool {
	username = strings.TrimSpace(username)
	for _, admin := range app.cfg.Server.Auth.Admins {
		if strings.EqualFold(strings.TrimSpace(admin), username) {
			return true
		}
	}
	return false
}
