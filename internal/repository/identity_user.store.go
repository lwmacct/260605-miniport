package repository

import (
	"context"
	"strings"
	"time"
)

func (s *Store) CreateIdentityUser(ctx context.Context, username string, displayName string) (*IdentityUserModel, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		displayName = username
	}
	now := time.Now().UTC()
	user := &IdentityUserModel{
		Username:    username,
		DisplayName: displayName,
		Status:      IdentityUserModelStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if _, err := s.db.NewInsert().Model(user).Exec(ctx); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) FetchIdentityUserByID(ctx context.Context, id int64) (*IdentityUserModel, error) {
	user := new(IdentityUserModel)
	err := s.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	return user, WrapNotFound(err)
}

func (s *Store) FetchIdentityUserByUsername(ctx context.Context, username string) (*IdentityUserModel, error) {
	user := new(IdentityUserModel)
	err := s.db.NewSelect().Model(user).Where("username = ?", strings.ToLower(strings.TrimSpace(username))).Scan(ctx)
	return user, WrapNotFound(err)
}

func (s *Store) ListIdentityUsers(ctx context.Context) ([]IdentityUserModel, error) {
	var users []IdentityUserModel
	if err := s.db.NewSelect().Model(&users).Order("username ASC").Scan(ctx); err != nil {
		return nil, err
	}
	return users, nil
}
