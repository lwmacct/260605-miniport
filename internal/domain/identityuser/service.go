package identityuser

import (
	"context"

	"github.com/uptrace/bun"
)

type Service struct {
	repo *repository
}

func NewService(db bun.IDB) *Service {
	return &Service{repo: newRepository(db)}
}

func (svc *Service) Create(ctx context.Context, input CreateInput) (*User, error) {
	return svc.repo.create(ctx, input)
}

func (svc *Service) ByID(ctx context.Context, id int64) (*User, error) {
	return svc.repo.byID(ctx, id)
}

func (svc *Service) ByUsername(ctx context.Context, username string) (*User, error) {
	return svc.repo.byUsername(ctx, username)
}

func (svc *Service) List(ctx context.Context) ([]User, error) {
	return svc.repo.list(ctx)
}

func (svc *Service) EnsureActive(user *User) error {
	return ensureActive(user)
}
