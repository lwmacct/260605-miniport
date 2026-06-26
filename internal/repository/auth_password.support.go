package repository

import "time"

type AuthPasswordRecord struct {
	IdentityUserID int64
	PasswordHash   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
