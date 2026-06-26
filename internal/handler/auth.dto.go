package handler

import "time"

type AuthConfigDTO struct {
	Local struct {
		LoginEnabled        bool `json:"loginEnabled"`
		RegistrationEnabled bool `json:"registrationEnabled"`
	} `json:"local"`
	Challenge struct {
		Provider string `json:"provider"`
		SiteKey  string `json:"sitekey,omitempty"`
	} `json:"challenge"`
}

type AuthChallengeCreateDTO struct {
	Provider    string    `json:"provider"`
	ChallengeID string    `json:"challengeId,omitempty"`
	Image       string    `json:"image,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
}

type AuthChallengeDTO struct {
	Provider    string `json:"provider"`
	ChallengeID string `json:"challengeId,omitempty"`
	Answer      string `json:"answer,omitempty"`
	Token       string `json:"token,omitempty"`
}

type AuthCredentialsDTO struct {
	Username  string           `json:"username"`
	Password  string           `json:"password"`
	Challenge AuthChallengeDTO `json:"challenge"`
}

type AuthUserDTO struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`
	Admin       bool   `json:"admin"`
}

type AuthSessionDTO struct {
	Authenticated bool         `json:"authenticated"`
	ExpiresAt     string       `json:"expiresAt,omitempty"`
	User          *AuthUserDTO `json:"user,omitempty"`
}

type AuthStateDTO struct {
	Config  AuthConfigDTO  `json:"config"`
	Session AuthSessionDTO `json:"session"`
}

type AuthSessionResponseDTO struct {
	SetCookie string `header:"Set-Cookie"`
	Body      AuthSessionDTO
}

type AuthSessionInputDTO struct {
	Session string `cookie:"web_session"`
}
