package identityuser

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type repository struct {
	db bun.IDB
}

func newRepository(db bun.IDB) *repository {
	if db == nil {
		panic("identityuser.newRepository: db is nil")
	}
	return &repository{db: db}
}

func (repo *repository) create(ctx context.Context, input CreateInput) (*User, error) {
	username := normalizeUsername(input.Username)
	if username == "" {
		return nil, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	user := &User{
		Username:    username,
		DisplayName: defaultDisplayName(username, input.DisplayName),
		Status:      StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if _, err := repo.db.NewInsert().Model(user).Exec(ctx); err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("%w: %w", ErrUsernameTaken, err)
		}
		return nil, err
	}
	return user, nil
}

func (repo *repository) byID(ctx context.Context, id int64) (*User, error) {
	user := new(User)
	if err := repo.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *repository) byUsername(ctx context.Context, username string) (*User, error) {
	user := new(User)
	if err := repo.db.NewSelect().Model(user).Where("username = ?", normalizeUsername(username)).Scan(ctx); err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *repository) list(ctx context.Context) ([]User, error) {
	var users []User
	if err := repo.db.NewSelect().Model(&users).Order("username ASC").Scan(ctx); err != nil {
		return nil, err
	}
	return users, nil
}

func ensureActive(user *User) error {
	if user == nil {
		return sql.ErrNoRows
	}
	if user.Status == StatusDisabled || user.DisabledAt != nil {
		return ErrDisabled
	}
	return nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	text := strings.ToLower(err.Error())
	return strings.Contains(text, "unique") || strings.Contains(text, "duplicate")
}
