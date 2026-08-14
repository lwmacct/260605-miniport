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
	h := githubHandler{service: githubService}
	console := huma.NewGroup(api, "/console")
	tags := []string{"Console", "GitHub"}
	huma.Register(console, huma.Operation{OperationID: "console-get-github-status", Method: http.MethodGet, Path: "/github/status", Summary: "Get GitHub status", Tags: tags}, h.status)
	huma.Register(console, huma.Operation{OperationID: "console-begin-github-connection", Method: http.MethodPost, Path: "/github/connections", Summary: "Begin GitHub connection", Tags: tags}, h.beginConnection)
	huma.Register(console, huma.Operation{OperationID: "console-list-github-installations", Method: http.MethodGet, Path: "/github/installations", Summary: "List GitHub installations", Tags: tags}, h.listInstallations)
	huma.Register(console, huma.Operation{OperationID: "console-sync-github-installations", Method: http.MethodPost, Path: "/github/installations/sync", Summary: "Sync GitHub installations", Tags: tags}, h.syncInstallations)
	huma.Register(console, huma.Operation{OperationID: "console-list-github-repositories", Method: http.MethodGet, Path: "/github/repositories", Summary: "List GitHub repositories", Tags: tags}, h.listRepositories)

	integrations := huma.NewGroup(api, "/integrations/github")
	integrationTags := []string{"GitHub Integration"}
	huma.Register(integrations, huma.Operation{OperationID: "complete-github-connection", Method: http.MethodGet, Path: "/setup", DefaultStatus: http.StatusSeeOther, Summary: "Complete GitHub connection", Tags: integrationTags}, h.completeConnection)
	huma.Register(integrations, huma.Operation{OperationID: "receive-github-webhook", Method: http.MethodPost, Path: "/webhooks", Summary: "Receive GitHub webhook", Tags: integrationTags, MaxBodyBytes: 2 << 20}, h.webhook)
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

func (h githubHandler) listInstallations(ctx context.Context, _ *struct{}) (*BodyDTO[GithubInstallationListDTO], error) {
	items, err := h.service.ListInstallations(ctx)
	if err != nil {
		return nil, utilGithubAPIError(err)
	}
	return &BodyDTO[GithubInstallationListDTO]{Body: GithubInstallationListDTO{Items: ToGithubInstallationDTOs(items)}}, nil
}

func (h githubHandler) syncInstallations(ctx context.Context, input *GithubInstallationSyncInputDTO) (*BodyDTO[GithubInstallationBatchDTO], error) {
	items, err := h.service.SyncInstallations(ctx, input.Body.IDs)
	if err != nil {
		return nil, utilGithubAPIError(err)
	}
	return &BodyDTO[GithubInstallationBatchDTO]{Body: GithubInstallationBatchDTO{Items: ToGithubInstallationDTOs(items)}}, nil
}

func (h githubHandler) listRepositories(ctx context.Context, input *GithubRepositoryListInputDTO) (*BodyDTO[GithubRepositoryListDTO], error) {
	items, err := h.service.ListRepositories(ctx, input.Query, input.State)
	if err != nil {
		return nil, utilGithubAPIError(err)
	}
	return &BodyDTO[GithubRepositoryListDTO]{Body: GithubRepositoryListDTO{Items: ToGithubRepositoryDTOs(items)}}, nil
}

func (h githubHandler) webhook(ctx context.Context, input *GithubWebhookInputDTO) (*ActionOutputDTO, error) {
	if err := h.service.HandleWebhook(ctx, input.Delivery, input.Event, input.Signature, input.RawBody); err != nil {
		return nil, utilGithubAPIError(err)
	}
	return &ActionOutputDTO{Body: ActionDTO{OK: true}}, nil
}
