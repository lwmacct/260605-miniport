package service

import (
	"context"

	"github.com/lwmacct/260605-miniport/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthPasswordService struct {
	store *repository.Store
}

func NewAuthPasswordService(store *repository.Store) *AuthPasswordService {
	if store == nil {
		panic("NewAuthPasswordService: store is nil")
	}
	return &AuthPasswordService{store: store}
}

func (s *AuthPasswordService) CheckStrength(username string, password string) error {
	return validatePassword(username, password)
}

func (s *AuthPasswordService) Set(ctx context.Context, username string, identityUserID int64, password string) error {
	if err := validatePassword(username, password); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.store.CreateAuthPassword(ctx, identityUserID, string(hash))
	return err
}

func (s *AuthPasswordService) Register(ctx context.Context, input AuthPasswordRegisterInput) (*IdentityUser, error) {
	var user *IdentityUser
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		users := NewIdentityUserService(txStore)
		passwords := NewAuthPasswordService(txStore)
		created, err := users.Create(ctx, CreateIdentityUserInput{Username: input.Username})
		if err != nil {
			return err
		}
		if err := passwords.Set(ctx, created.Username, created.ID, input.Password); err != nil {
			return err
		}
		user = created
		return nil
	})
	return user, err
}

func (s *AuthPasswordService) Authenticate(ctx context.Context, username string, password string, users *IdentityUserService) (*IdentityUser, error) {
	if utilNormalizeUsername(username) == "" || password == "" {
		return nil, ErrAuthPasswordInvalidCredentials
	}
	user, err := users.ByUsername(ctx, username)
	if err != nil {
		if IsNotFound(err) {
			return nil, ErrAuthPasswordInvalidCredentials
		}
		return nil, err
	}
	activeErr := users.EnsureActive(user)
	if activeErr != nil {
		return nil, activeErr
	}
	row, err := s.store.FetchAuthPasswordByIdentityUserID(ctx, user.ID)
	if err != nil {
		if IsNotFound(err) {
			return nil, ErrAuthPasswordInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(row.PasswordHash), []byte(password)); err != nil {
		return nil, ErrAuthPasswordInvalidCredentials
	}
	return user, nil
}
