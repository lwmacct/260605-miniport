package service

import (
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

const (
	IdentityUserStatusActive   = repository.IdentityUserModelStatusActive
	IdentityUserStatusDisabled = repository.IdentityUserModelStatusDisabled
)

type IdentityUser struct {
	ID          int64
	Username    string
	DisplayName string
	Status      string
	DisabledAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateIdentityUserInput struct {
	Username    string
	DisplayName string
}

var (
	ErrIdentityUserUsernameTaken      = errors.New("username taken")
	ErrIdentityUserInvalidCredentials = errors.New("invalid credentials")
	ErrIdentityUserDisabled           = errors.New("user disabled")
)

func utilIdentityUser(row *repository.IdentityUserModel) *IdentityUser {
	if row == nil {
		return nil
	}
	return &IdentityUser{
		ID:          row.ID,
		Username:    row.Username,
		DisplayName: row.DisplayName,
		Status:      row.Status,
		DisabledAt:  row.DisabledAt,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func IsNotFound(err error) bool {
	return errors.Is(err, repository.ErrNotFound)
}

func utilNormalizeUsername(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
