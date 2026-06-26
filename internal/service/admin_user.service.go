package service

import "context"

type AdminUserService struct {
	users          *UserService
	runtimeAdminFn func(string) bool
}

func NewAdminUserService(users *UserService, runtimeAdminFn func(string) bool) *AdminUserService {
	if users == nil {
		panic("NewAdminUserService: users is nil")
	}
	if runtimeAdminFn == nil {
		runtimeAdminFn = func(string) bool { return false }
	}
	return &AdminUserService{users: users, runtimeAdminFn: runtimeAdminFn}
}

func (s *AdminUserService) List(ctx context.Context) ([]AdminUser, error) {
	users, err := s.users.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]AdminUser, 0, len(users))
	for _, user := range users {
		out = append(out, utilAdminUser(user, s.runtimeAdminFn(user.Username)))
	}
	return out, nil
}
