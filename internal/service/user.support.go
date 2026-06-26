package service

import (
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

const (
	UserStatusActive   = repository.UserModelStatusActive
	UserStatusDisabled = repository.UserModelStatusDisabled
)

type User struct {
	ID          int64
	Username    string
	DisplayName string
	Status      string
	DisabledAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateUserInput struct {
	Username    string
	DisplayName string
}

var (
	ErrUserUsernameTaken      = errors.New("username taken")
	ErrUserInvalidCredentials = errors.New("invalid credentials")
	ErrUserDisabled           = errors.New("user disabled")
)

func utilUser(row *repository.UserModel) *User {
	if row == nil {
		return nil
	}
	return &User{
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
