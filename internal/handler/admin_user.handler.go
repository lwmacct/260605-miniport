package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type adminUserHandler struct {
	auth authHandler
}

func RegisterAdminUser(api huma.API, config Config, services Services) {
	handler := adminUserHandler{auth: authHandler{config: config, services: services}}
	admin := huma.NewGroup(api, "/admin")
	huma.Register(admin, huma.Operation{OperationID: "list-admin-users", Method: http.MethodGet, Path: "/users", Tags: []string{"Admin"}}, handler.listUsers)
}

func (h adminUserHandler) listUsers(ctx context.Context, input *AdminUserListInputDTO) (*BodyDTO[[]AdminUserDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	users, err := h.auth.services.AdminUsers.List(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}
	body := make([]AdminUserDTO, 0, len(users))
	for _, user := range users {
		body = append(body, ToAdminUserDTO(user))
	}
	return &BodyDTO[[]AdminUserDTO]{Body: body}, nil
}

func (h adminUserHandler) requireAdmin(ctx context.Context, sessionID string) (*AuthUserDTO, error) {
	session, err := h.auth.session(ctx, sessionID)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	if session.User == nil || !session.User.Admin {
		return nil, huma.Error403Forbidden("forbidden")
	}
	return session.User, nil
}
