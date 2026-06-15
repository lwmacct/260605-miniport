package authpassword

import (
	"strings"
	"unicode/utf8"
)

func validatePassword(username string, password string) error {
	trimmed := strings.TrimSpace(password)
	if len(trimmed) < 10 {
		return ErrWeakPassword
	}

	normalized := strings.ToLower(trimmed)
	username = strings.ToLower(strings.TrimSpace(username))
	if username != "" && utf8.RuneCountInString(username) >= 3 && strings.Contains(normalized, username) {
		return ErrWeakPassword
	}
	return nil
}
