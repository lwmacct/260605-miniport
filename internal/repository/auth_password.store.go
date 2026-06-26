package repository

import (
	"context"
	"time"
)

func (s *Store) CreateAuthPassword(ctx context.Context, identityUserID int64, hash string) (*AuthPasswordModel, error) {
	now := time.Now().UTC()
	row := &AuthPasswordModel{
		IdentityUserID: identityUserID,
		PasswordHash:   hash,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Store) FetchAuthPasswordByIdentityUserID(ctx context.Context, identityUserID int64) (*AuthPasswordModel, error) {
	row := new(AuthPasswordModel)
	err := s.db.NewSelect().Model(row).Where("identity_user_id = ?", identityUserID).Scan(ctx)
	return row, WrapNotFound(err)
}
