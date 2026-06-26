package service

import (
	"context"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

type IdentityUserService struct {
	store *repository.Store
}

func NewIdentityUserService(store *repository.Store) *IdentityUserService {
	if store == nil {
		panic("NewIdentityUserService: store is nil")
	}
	return &IdentityUserService{store: store}
}

func (s *IdentityUserService) Create(ctx context.Context, input CreateIdentityUserInput) (*IdentityUser, error) {
	username := utilNormalizeUsername(input.Username)
	if username == "" {
		return nil, ErrIdentityUserInvalidCredentials
	}
	user, err := s.store.CreateIdentityUser(ctx, username, input.DisplayName)
	if err != nil {
		if repository.IsUniqueViolation(err) {
			return nil, ErrIdentityUserUsernameTaken
		}
		return nil, err
	}
	return utilIdentityUser(user), nil
}

func (s *IdentityUserService) ByID(ctx context.Context, id int64) (*IdentityUser, error) {
	user, err := s.store.FetchIdentityUserByID(ctx, id)
	return utilIdentityUser(user), err
}

func (s *IdentityUserService) ByUsername(ctx context.Context, username string) (*IdentityUser, error) {
	user, err := s.store.FetchIdentityUserByUsername(ctx, username)
	return utilIdentityUser(user), err
}

func (s *IdentityUserService) List(ctx context.Context) ([]IdentityUser, error) {
	rows, err := s.store.ListIdentityUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]IdentityUser, 0, len(rows))
	for idx := range rows {
		users = append(users, *utilIdentityUser(&rows[idx]))
	}
	return users, nil
}

func (s *IdentityUserService) EnsureActive(user *IdentityUser) error {
	if user == nil {
		return repository.ErrNotFound
	}
	if user.Status == IdentityUserStatusDisabled || user.DisabledAt != nil {
		return ErrIdentityUserDisabled
	}
	return nil
}
