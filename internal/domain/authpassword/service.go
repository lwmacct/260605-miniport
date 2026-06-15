package authpassword

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"

	"github.com/lwmacct/260605-miniport/internal/domain/identityuser"
)

type Service struct {
	db bun.IDB
}

func NewService(db bun.IDB) *Service {
	if db == nil {
		panic("authpassword.NewService: db is nil")
	}
	return &Service{db: db}
}

func (svc *Service) CheckStrength(username string, password string) error {
	return validatePassword(username, password)
}

func (svc *Service) Set(ctx context.Context, username string, userID int64, password string) error {
	if err := validatePassword(username, password); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	row := &Password{
		UserID:       userID,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	_, err = svc.db.NewInsert().Model(row).Exec(ctx)
	return err
}

func (svc *Service) Authenticate(ctx context.Context, username string, password string, users *identityuser.Service) (*identityuser.User, error) {
	if strings.TrimSpace(username) == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := users.ByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := users.EnsureActive(user); err != nil {
		return nil, err
	}

	row := new(Password)
	if err := svc.db.NewSelect().Model(row).Where("user_id = ?", user.ID).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(row.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
