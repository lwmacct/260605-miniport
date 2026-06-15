package authchallenge

import "time"

type RequestDTO struct {
	IP         string
	UserAgent  string
	Method     string
	Path       string
	RemoteAddr string
}

type PublicConfigDTO struct {
	Provider string `json:"provider"`
	SiteKey  string `json:"sitekey,omitempty"`
}

type ChallengeDTO struct {
	Provider    string    `json:"provider"`
	ChallengeID string    `json:"challengeId,omitempty"`
	Image       string    `json:"image,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
}

type ResponseDTO struct {
	Provider    string `json:"provider"`
	ChallengeID string `json:"challengeId,omitempty"`
	Answer      string `json:"answer,omitempty"`
	Token       string `json:"token,omitempty"`
}
