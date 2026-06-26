package repository

import "time"

type AuthSessionCreateInput struct {
	IDHash         string
	IdentityUserID int64
	LoginIP        string
	LastIP         string
	UserAgentHash  string
	ExpiresAt      time.Time
	CreatedAt      time.Time
	LastSeenAt     time.Time
}

type AuthSessionRecord struct {
	IDHash         string
	IdentityUserID int64
	LoginIP        string
	LastIP         string
	UserAgentHash  string
	ExpiresAt      time.Time
	RevokedAt      *time.Time
	CreatedAt      time.Time
	LastSeenAt     time.Time
}
