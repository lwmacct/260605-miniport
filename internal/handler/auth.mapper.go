package handler

import "github.com/lwmacct/260605-miniport/internal/service"

func ToAuthUserDTO(user *service.User, admin bool) *AuthUserDTO {
	if user == nil {
		return nil
	}
	return &AuthUserDTO{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Status:      user.Status,
		Admin:       admin,
	}
}

func ToAuthChallengeInput(request service.AuthSessionInput) service.AuthChallengeInput {
	return service.AuthChallengeInput{
		IP:         request.IP,
		UserAgent:  request.UserAgent,
		Method:     request.Method,
		Path:       request.Path,
		RemoteAddr: request.RemoteAddr,
	}
}

func ToAuthChallengeAnswer(challenge AuthChallengeDTO) service.AuthChallengeAnswer {
	return service.AuthChallengeAnswer{
		Provider:    challenge.Provider,
		ChallengeID: challenge.ChallengeID,
		Answer:      challenge.Answer,
		Token:       challenge.Token,
	}
}

func ToAuthChallengeCreateDTO(challenge *service.AuthChallenge) AuthChallengeCreateDTO {
	if challenge == nil {
		return AuthChallengeCreateDTO{}
	}
	return AuthChallengeCreateDTO{
		Provider:    challenge.Provider,
		ChallengeID: challenge.ChallengeID,
		Image:       challenge.Image,
		ExpiresAt:   challenge.ExpiresAt,
	}
}

func ToAuthSessionDTO(sessionUser *service.AuthSessionUser, user *service.User, admin bool) AuthSessionDTO {
	return AuthSessionDTO{
		Authenticated: true,
		ExpiresAt:     utilHTTPTime(sessionUser.ExpiresAt),
		User:          ToAuthUserDTO(user, admin),
	}
}
