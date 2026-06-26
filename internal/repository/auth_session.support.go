package repository

import "time"

type AuthSessionCreateInput struct {
	IDHash        string
	UserID        int64
	LoginIP       string
	LastIP        string
	UserAgentHash string
	ExpiresAt     time.Time
	CreatedAt     time.Time
	LastSeenAt    time.Time
}

type AuthSessionRecord struct {
	IDHash        string
	UserID        int64
	LoginIP       string
	LastIP        string
	UserAgentHash string
	ExpiresAt     time.Time
	RevokedAt     *time.Time
	CreatedAt     time.Time
	LastSeenAt    time.Time
}

func utilAuthSessionRecordFromModel(model *AuthSessionModel) *AuthSessionRecord {
	if model == nil {
		return nil
	}
	return &AuthSessionRecord{
		IDHash:        model.IDHash,
		UserID:        model.UserID,
		LoginIP:       model.LoginIP,
		LastIP:        model.LastIP,
		UserAgentHash: model.UserAgentHash,
		ExpiresAt:     model.ExpiresAt,
		RevokedAt:     model.RevokedAt,
		CreatedAt:     model.CreatedAt,
		LastSeenAt:    model.LastSeenAt,
	}
}
