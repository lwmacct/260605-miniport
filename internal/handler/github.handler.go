package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/identity"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/requestctx"

	"github.com/lwmacct/260605-miniport/internal/service"
)

type githubHandler struct {
	config  Config
	service *service.GithubService
}

func RegisterGithub(api huma.API, config Config, githubService *service.GithubService) {
	handler := githubHandler{config: config, service: githubService}
	huma.Register(api, huma.Operation{OperationID: "get-github-status", Method: http.MethodGet, Path: "/github/status", Tags: []string{"github"}}, handler.status)
	huma.Register(api, huma.Operation{OperationID: "begin-github-connection", Method: http.MethodPost, Path: "/github/connections", Tags: []string{"github"}}, handler.beginConnection)
	huma.Register(api, huma.Operation{OperationID: "complete-github-connection", Method: http.MethodGet, Path: "/github/setup", Tags: []string{"github"}}, handler.completeConnection)
	huma.Register(api, huma.Operation{OperationID: "list-github-installations", Method: http.MethodGet, Path: "/github/installations", Tags: []string{"github"}}, handler.listInstallations)
	huma.Register(api, huma.Operation{OperationID: "disconnect-github-installation", Method: http.MethodDelete, Path: "/github/installations/{id}/connection", Tags: []string{"github"}}, handler.disconnect)
	huma.Register(api, huma.Operation{OperationID: "sync-github-installation", Method: http.MethodPost, Path: "/github/installations/{id}/sync", Tags: []string{"github"}}, handler.syncInstallation)
	huma.Register(api, huma.Operation{OperationID: "list-github-repositories", Method: http.MethodGet, Path: "/github/repositories", Tags: []string{"github"}}, handler.listRepositories)
	huma.Register(api, huma.Operation{OperationID: "receive-github-webhook", Method: http.MethodPost, Path: "/github/webhooks", Tags: []string{"github"}, MaxBodyBytes: 2 << 20}, handler.webhook)
}

func (h githubHandler) status(ctx context.Context, input *GithubSessionInputDTO) (*BodyDTO[GithubStatusDTO], error) {
	if _, err := h.actor(ctx, input.Session); err != nil {
		return nil, err
	}
	status := h.service.Status()
	return &BodyDTO[GithubStatusDTO]{Body: GithubStatusDTO{Enabled: status.Enabled, AppSlug: status.AppSlug, InstallURL: status.InstallURL}}, nil
}

func (h githubHandler) beginConnection(ctx context.Context, input *GithubSessionInputDTO) (*BodyDTO[GithubConnectionStartDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	start, err := h.service.BeginConnection(ctx, actor.OwnerSubject)
	if err != nil {
		return nil, utilGithubAPIError(err)
	}
	return &BodyDTO[GithubConnectionStartDTO]{Body: GithubConnectionStartDTO{URL: start.URL}}, nil
}

func (h githubHandler) completeConnection(ctx context.Context, input *GithubSetupInputDTO) (*RedirectDTO, error) {
	if err := h.service.CompleteConnection(ctx, input.State, input.InstallationID); err != nil {
		return nil, utilGithubAPIError(err)
	}
	return &RedirectDTO{Status: http.StatusSeeOther, Location: h.service.SetupReturnURL()}, nil
}

func (h githubHandler) listInstallations(ctx context.Context, input *GithubSessionInputDTO) (*BodyDTO[[]GithubInstallationDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	items, err := h.service.ListInstallations(ctx, actor.OwnerSubject)
	if err != nil {
		return nil, utilGithubAPIError(err)
	}
	return &BodyDTO[[]GithubInstallationDTO]{Body: ToGithubInstallationDTOs(items)}, nil
}

func (h githubHandler) disconnect(ctx context.Context, input *GithubInstallationInputDTO) (*BodyDTO[DeleteDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.service.Disconnect(ctx, actor.OwnerSubject, input.ID); err != nil {
		return nil, utilGithubAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h githubHandler) syncInstallation(ctx context.Context, input *GithubInstallationInputDTO) (*BodyDTO[GithubInstallationDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if syncErr := h.service.SyncInstallation(ctx, actor.OwnerSubject, input.ID); syncErr != nil {
		return nil, utilGithubAPIError(syncErr)
	}
	items, err := h.service.ListInstallations(ctx, actor.OwnerSubject)
	if err != nil {
		return nil, utilGithubAPIError(err)
	}
	for _, item := range items {
		if item.ID == input.ID {
			return &BodyDTO[GithubInstallationDTO]{Body: ToGithubInstallationDTO(item)}, nil
		}
	}
	return nil, huma.Error404NotFound("github installation not found")
}

func (h githubHandler) listRepositories(ctx context.Context, input *GithubRepositoryListInputDTO) (*BodyDTO[[]GithubRepositoryDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	items, err := h.service.ListRepositories(ctx, actor.OwnerSubject, input.Query, input.State)
	if err != nil {
		return nil, utilGithubAPIError(err)
	}
	return &BodyDTO[[]GithubRepositoryDTO]{Body: ToGithubRepositoryDTOs(items)}, nil
}

func (h githubHandler) webhook(ctx context.Context, input *GithubWebhookInputDTO) (*BodyDTO[DeleteDTO], error) {
	if err := h.service.HandleWebhook(ctx, input.Delivery, input.Event, input.Signature, input.RawBody); err != nil {
		return nil, utilGithubAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: false}}, nil
}

func (h githubHandler) actor(ctx context.Context, sessionID string) (service.PortsvcActor, error) {
	if sessionID == "" || h.config.Identity == nil {
		return service.PortsvcActor{}, huma.Error401Unauthorized("unauthorized")
	}
	request, ok := requestctx.RequestFromContext(ctx)
	if !ok {
		return service.PortsvcActor{}, huma.Error401Unauthorized("unauthorized")
	}
	principal, err := h.config.Identity.CurrentPrincipal(ctx, sessionID, request)
	if err != nil || principal == nil || !principal.Active() || principal.Status == identity.StatusDisabled {
		return service.PortsvcActor{}, huma.Error401Unauthorized("unauthorized")
	}
	actor := ToPortsvcActor(principal)
	if actor.OwnerSubject == "" {
		return service.PortsvcActor{}, huma.Error401Unauthorized("unauthorized")
	}
	return actor, nil
}
