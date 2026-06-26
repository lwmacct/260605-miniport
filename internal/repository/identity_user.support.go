package repository

import "time"

const (
	IdentityUserModelStatusActive   = "active"
	IdentityUserModelStatusDisabled = "disabled"
)

type IdentityUserRecord struct {
	ID          int64
	Username    string
	DisplayName string
	Status      string
	DisabledAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
