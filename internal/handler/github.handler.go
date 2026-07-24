package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260605-miniport/internal/service"
)

type githubHandler struct {
	service *service.GithubService
}

func RegisterGithub(api huma.API, githubService *service.GithubService) {
	handler := githubHandler{service: githubService}
	huma.Register(api, huma.Operation{OperationID: "get-github-status", Method: http.MethodGet, Path: "/github/status", Tags: []string{"github"}}, handler.status)
	huma.Register(api, huma.Operation{OperationID: "begin-github-connection", Method: http.MethodPost, Path: "/github/connections", Tags: []string{"github"}}, handler.beginConnection)
	huma.Register(api, huma.Operation{OperationID: "complete-github-connection", Method: http.MethodGet, Path: "/github/setup", Tags: []string{"github"}}, handler.completeConnection)
	huma.Register(api, huma.Operation{OperationID: "list-github-installations", Method: http.MethodGet, Path: "/github/installations", Tags: []string{"github"}}, handler.listInstallations)
	huma.Register(api, huma.Operation{OperationID: "sync-github-installation", Method: http.MethodPost, Path: "/github/installations/{id}/sync", Tags: []string{"github"}}, handler.syncInstallation)
	huma.Register(api, huma.Operation{OperationID: "list-github-repositories", Method: http.MethodGet, Path: "/github/repositories", Tags: []string{"github"}}, handler.listRepositories)
	huma.Register(api, huma.Operation{OperationID: "receive-github-webhook", Method: http.MethodPost, Path: "/github/webhooks", Tags: []string{"github"}, MaxBodyBytes: 2 << 20}, handler.webhook)
}

func (h githubHandler) status(_ context.Context, _ *struct{}) (*BodyDTO[GithubStatusDTO], error) {
	status := h.service.Status()
	return &BodyDTO[GithubStatusDTO]{Body: GithubStatusDTO{Enabled: status.Enabled, AppSlug: status.AppSlug, InstallURL: status.InstallURL}}, nil
}

func (h githubHandler) beginConnection(ctx context.Context, _ *struct{}) (*BodyDTO[GithubConnectionStartDTO], error) {
	start, err := h.service.BeginConnection(ctx)
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

func (h githubHandler) listInstallations(ctx context.Context, _ *struct{}) (*BodyDTO[[]GithubInstallationDTO], error) {
	items, err := h.service.ListInstallations(ctx)
	if err != nil {
		return nil, utilGithubAPIError(err)
	}
	return &BodyDTO[[]GithubInstallationDTO]{Body: ToGithubInstallationDTOs(items)}, nil
}

func (h githubHandler) syncInstallation(ctx context.Context, input *GithubInstallationInputDTO) (*BodyDTO[GithubInstallationDTO], error) {
	if syncErr := h.service.SyncInstallation(ctx, input.ID); syncErr != nil {
		return nil, utilGithubAPIError(syncErr)
	}
	items, err := h.service.ListInstallations(ctx)
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
	items, err := h.service.ListRepositories(ctx, input.Query, input.State)
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
