package service

import (
	"context"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

type UserService struct {
	store *repository.Store
}

func NewUserService(store *repository.Store) *UserService {
	if store == nil {
		panic("NewUserService: store is nil")
	}
	return &UserService{store: store}
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	username := utilNormalizeUsername(input.Username)
	if username == "" {
		return nil, ErrUserInvalidCredentials
	}
	user, err := s.store.CreateUser(ctx, username, input.DisplayName)
	if err != nil {
		if repository.IsUniqueViolation(err) {
			return nil, ErrUserUsernameTaken
		}
		return nil, err
	}
	return utilUser(user), nil
}

func (s *UserService) ByID(ctx context.Context, id int64) (*User, error) {
	user, err := s.store.FetchUserByID(ctx, id)
	return utilUser(user), err
}

func (s *UserService) ByUsername(ctx context.Context, username string) (*User, error) {
	user, err := s.store.FetchUserByUsername(ctx, username)
	return utilUser(user), err
}

func (s *UserService) List(ctx context.Context) ([]User, error) {
	rows, err := s.store.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(rows))
	for idx := range rows {
		users = append(users, *utilUser(&rows[idx]))
	}
	return users, nil
}

func (s *UserService) EnsureActive(user *User) error {
	if user == nil {
		return repository.ErrNotFound
	}
	if user.Status == UserStatusDisabled || user.DisabledAt != nil {
		return ErrUserDisabled
	}
	return nil
}
