package repository

import "time"

type AuthPasswordRecord struct {
	UserID       int64
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
