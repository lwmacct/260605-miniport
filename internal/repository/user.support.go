package repository

import "time"

const (
	UserModelStatusActive   = "active"
	UserModelStatusDisabled = "disabled"
)

type UserRecord struct {
	ID          int64
	Username    string
	DisplayName string
	Status      string
	DisabledAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
