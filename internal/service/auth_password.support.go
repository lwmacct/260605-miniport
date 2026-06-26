package service

import (
	"errors"
	"strings"
)

var (
	ErrAuthPasswordWeakPassword       = errors.New("weak password")
	ErrAuthPasswordInvalidCredentials = errors.New("invalid credentials")
)

func validatePassword(username string, password string) error {
	password = strings.TrimSpace(password)
	if len(password) < 8 {
		return ErrAuthPasswordWeakPassword
	}
	if strings.EqualFold(strings.TrimSpace(username), password) {
		return ErrAuthPasswordWeakPassword
	}
	return nil
}

type AuthPasswordRegisterInput struct {
	Username string
	Password string
}
