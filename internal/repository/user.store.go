package repository

import (
	"context"
	"strings"
	"time"
)

func (s *Store) CreateUser(ctx context.Context, username string, displayName string) (*UserModel, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		displayName = username
	}
	now := time.Now().UTC()
	user := &UserModel{
		Username:    username,
		DisplayName: displayName,
		Status:      UserModelStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if _, err := s.db.NewInsert().Model(user).Exec(ctx); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) FetchUserByID(ctx context.Context, id int64) (*UserModel, error) {
	user := new(UserModel)
	err := s.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	return user, WrapNotFound(err)
}

func (s *Store) FetchUserByUsername(ctx context.Context, username string) (*UserModel, error) {
	user := new(UserModel)
	err := s.db.NewSelect().Model(user).Where("username = ?", strings.ToLower(strings.TrimSpace(username))).Scan(ctx)
	return user, WrapNotFound(err)
}

func (s *Store) ListUsers(ctx context.Context) ([]UserModel, error) {
	var users []UserModel
	if err := s.db.NewSelect().Model(&users).Order("username ASC").Scan(ctx); err != nil {
		return nil, err
	}
	return users, nil
}
