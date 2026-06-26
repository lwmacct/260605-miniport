package repository

import (
	"context"
	"time"
)

func (s *Store) CreateAuthPassword(ctx context.Context, userID int64, hash string) (*AuthPasswordModel, error) {
	now := time.Now().UTC()
	row := &AuthPasswordModel{
		UserID:       userID,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Store) FetchAuthPasswordByUserID(ctx context.Context, userID int64) (*AuthPasswordModel, error) {
	row := new(AuthPasswordModel)
	err := s.db.NewSelect().Model(row).Where("user_id = ?", userID).Scan(ctx)
	return row, WrapNotFound(err)
}
