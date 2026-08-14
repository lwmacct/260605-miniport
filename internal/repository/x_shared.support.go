package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("repository not found")

func WrapNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	text := strings.ToLower(err.Error())
	return strings.Contains(text, "unique") || strings.Contains(text, "duplicate")
}

func IsUUID7(value string) bool {
	id, err := uuid.Parse(value)
	return err == nil && id.Version() == 7
}
