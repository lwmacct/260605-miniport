package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

type GithubService struct {
	store *repository.Store
	cfg   GithubConfig
	api   githubAPI
	now   func() time.Time
}

func NewGithubService(store *repository.Store, cfg GithubConfig) (*GithubService, error) {
	service := &GithubService{store: store, cfg: cfg, now: time.Now}
	if !cfg.Enabled {
		return service, nil
	}
	if cfg.AppID <= 0 || strings.TrimSpace(cfg.AppSlug) == "" || strings.TrimSpace(cfg.PrivateKeyFile) == "" {
		return nil, errors.New("github app-id, app-slug and private-key-file are required when github is enabled")
	}
	if strings.TrimSpace(cfg.WebhookSecret) == "" {
		return nil, errors.New("github webhook-secret is required when github is enabled")
	}
	client, err := NewGithubProvider(cfg)
	if err != nil {
		return nil, err
	}
	service.api = client
	return service, nil
}

func (s *GithubService) Status() GithubStatus {
	status := GithubStatus{Enabled: s.cfg.Enabled, AppSlug: s.cfg.AppSlug}
	if s.cfg.Enabled {
		status.InstallURL = utilGithubInstallURL(githubWebURL, s.cfg.AppSlug)
	}
	return status
}

func (s *GithubService) SetupReturnURL() string { return s.cfg.SetupReturnURL }

func (s *GithubService) BeginConnection(ctx context.Context) (*GithubConnectionStart, error) {
	if !s.cfg.Enabled {
		return nil, ErrGithubDisabled
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return nil, fmt.Errorf("generate github connection state: %w", err)
	}
	state := base64.RawURLEncoding.EncodeToString(raw)
	now := s.now().UTC()
	if err := s.store.AddGithubConnectionState(ctx, utilHashGithubState(state), now.Add(githubConnectionStateTTL), now); err != nil {
		return nil, err
	}
	return &GithubConnectionStart{URL: utilGithubConnectionURL(s.Status().InstallURL, state)}, nil
}

func (s *GithubService) CompleteConnection(ctx context.Context, state string, installationID int64) error {
	if !s.cfg.Enabled {
		return ErrGithubDisabled
	}
	if strings.TrimSpace(state) == "" || installationID <= 0 {
		return ErrGithubInvalidState
	}
	now := s.now().UTC()
	stateHash := utilHashGithubState(state)
	_, err := s.store.FetchGithubConnectionStateByHash(ctx, stateHash, now)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrGithubInvalidState
		}
		return err
	}
	consumed, err := s.store.DeleteGithubConnectionState(ctx, stateHash, now)
	if err != nil {
		return err
	}
	if !consumed {
		return ErrGithubInvalidState
	}
	installation, err := s.refreshInstallation(ctx, installationID)
	if err != nil {
		return err
	}
	return s.syncInstallation(ctx, *installation)
}

func (s *GithubService) ListInstallations(ctx context.Context) ([]GithubInstallation, error) {
	return s.store.ListGithubInstallations(ctx)
}

func (s *GithubService) ListRepositories(ctx context.Context, query, state string) ([]GithubRepository, error) {
	return s.store.ListGithubRepositories(ctx, query, state)
}

func (s *GithubService) SyncInstallation(ctx context.Context, installationID string) error {
	if !s.cfg.Enabled {
		return ErrGithubDisabled
	}
	installation, err := s.store.FetchGithubInstallationByID(ctx, installationID)
	if err != nil {
		return utilGithubServiceError(err)
	}
	refreshed, err := s.refreshInstallation(ctx, installation.GithubInstallationID)
	if err != nil {
		return err
	}
	return s.syncInstallation(ctx, *refreshed)
}

func (s *GithubService) Reconcile(ctx context.Context) error {
	if !s.cfg.Enabled {
		return nil
	}
	installations, err := s.store.ListGithubInstallationsByActiveStatus(ctx)
	if err != nil {
		return err
	}
	var result error
	for _, installation := range installations {
		refreshed, refreshErr := s.refreshInstallation(ctx, installation.GithubInstallationID)
		if refreshErr == nil {
			refreshErr = s.syncInstallation(ctx, *refreshed)
		}
		if refreshErr != nil {
			result = errors.Join(result, fmt.Errorf("reconcile github installation %d: %w", installation.GithubInstallationID, refreshErr))
		}
	}
	return result
}

func (s *GithubService) HandleWebhook(ctx context.Context, deliveryID, event, signature string, body []byte) error {
	if !s.cfg.Enabled {
		return ErrGithubDisabled
	}
	if !validateGithubWebhookSignature(s.cfg.WebhookSecret, signature, body) {
		return ErrGithubInvalidSignature
	}
	var envelope struct {
		Action       string `json:"action"`
		Installation struct {
			ID          int64      `json:"id"`
			SuspendedAt *time.Time `json:"suspended_at"`
		} `json:"installation"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("decode github webhook: %w", err)
	}
	if strings.TrimSpace(deliveryID) == "" || strings.TrimSpace(event) == "" {
		return errors.New("github webhook headers are required")
	}
	now := s.now().UTC()
	err := s.store.AddGithubWebhookDelivery(ctx, repository.GithubWebhookDeliveryRecord{
		DeliveryID: deliveryID, Event: event, Action: envelope.Action,
		InstallationID: envelope.Installation.ID, Status: "processing", ReceivedAt: now,
	})
	if errors.Is(err, repository.ErrGithubWebhookDeliveryExists) {
		return nil
	}
	if err != nil {
		return err
	}

	processErr := s.processWebhook(ctx, event, envelope.Action, envelope.Installation.ID, envelope.Installation.SuspendedAt)
	status := "processed"
	errorText := ""
	if processErr != nil {
		status = "failed"
		errorText = processErr.Error()
	}
	if finishErr := s.store.ReplaceGithubWebhookDeliveryResult(ctx, deliveryID, status, errorText, s.now().UTC()); finishErr != nil {
		return errors.Join(processErr, finishErr)
	}
	return processErr
}

func (s *GithubService) processWebhook(ctx context.Context, event, action string, installationID int64, suspendedAt *time.Time) error {
	if event == "ping" {
		return nil
	}
	if installationID <= 0 {
		return nil
	}
	if event == "installation" && action == "deleted" {
		return s.processDeletedInstallation(ctx, installationID)
	}
	if event == "installation" && action == "suspend" {
		return s.processSuspendedInstallation(ctx, installationID, suspendedAt)
	}
	if event != "installation" && event != "installation_repositories" && event != "repository" {
		return nil
	}
	installation, err := s.refreshInstallation(ctx, installationID)
	if err != nil {
		return err
	}
	return s.syncInstallation(ctx, *installation)
}

func (s *GithubService) processDeletedInstallation(ctx context.Context, installationID int64) error {
	installation, err := s.store.FetchGithubInstallationByExternalID(ctx, installationID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	now := s.now().UTC()
	if _, err = s.store.UpdateGithubInstallationStatus(ctx, installationID, "deleted", time.Time{}, now); err != nil {
		return err
	}
	return s.store.ReplaceGithubRepositoriesStateWithUnavailable(ctx, installation.ID, now)
}

func (s *GithubService) processSuspendedInstallation(ctx context.Context, installationID int64, suspendedAt *time.Time) error {
	when := s.now().UTC()
	if suspendedAt != nil {
		when = suspendedAt.UTC()
	}
	installation, err := s.store.UpdateGithubInstallationStatus(ctx, installationID, "suspended", when, s.now().UTC())
	if errors.Is(err, repository.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return s.store.ReplaceGithubRepositoriesStateWithUnavailable(ctx, installation.ID, s.now().UTC())
}

func (s *GithubService) refreshInstallation(ctx context.Context, externalID int64) (*GithubInstallation, error) {
	remote, err := s.api.Installation(ctx, externalID)
	if err != nil {
		return nil, err
	}
	permissions, err := json.Marshal(remote.Permissions)
	if err != nil {
		return nil, err
	}
	now := s.now().UTC()
	status := "active"
	var suspendedAt time.Time
	if remote.SuspendedAt != nil {
		status = "suspended"
		suspendedAt = remote.SuspendedAt.UTC()
	}
	return s.store.UpsertGithubInstallation(ctx, repository.GithubInstallationRecord{
		GithubInstallationID: remote.ID, AccountID: remote.Account.ID, AccountLogin: remote.Account.Login,
		AccountType: remote.Account.Type, AvatarURL: remote.Account.AvatarURL,
		RepositorySelection: remote.RepositorySelection, Permissions: string(permissions), Status: status,
		SuspendedAt: suspendedAt, CreatedAt: now, UpdatedAt: now,
	})
}

func (s *GithubService) syncInstallation(ctx context.Context, installation GithubInstallation) error {
	if installation.Status != "active" {
		return fmt.Errorf("github installation is %s", installation.Status)
	}
	remoteRepositories, err := s.api.Repositories(ctx, installation.GithubInstallationID)
	now := s.now().UTC()
	if err != nil {
		_, _ = s.store.UpdateGithubInstallationSync(ctx, installation.ID, now, err.Error())
		return err
	}
	repositories := make([]repository.GithubRepositoryRecord, 0, len(remoteRepositories))
	for _, item := range remoteRepositories {
		visibility := item.Visibility
		if visibility == "" {
			if item.Private {
				visibility = "private"
			} else {
				visibility = "public"
			}
		}
		repositories = append(repositories, repository.GithubRepositoryRecord{
			InstallationID: installation.ID, GithubRepositoryID: item.ID, NodeID: item.NodeID,
			OwnerLogin: item.Owner.Login, Name: item.Name, FullName: item.FullName, HTMLURL: item.HTMLURL,
			Description: item.Description, DefaultBranch: item.DefaultBranch, Visibility: visibility,
			Private: item.Private, Fork: item.Fork, Archived: item.Archived, Disabled: item.Disabled,
			PushedAt: item.PushedAt, RemoteUpdatedAt: item.UpdatedAt,
		})
	}
	if replaceErr := s.store.ReplaceGithubRepositories(ctx, installation.ID, repositories, now); replaceErr != nil {
		_, _ = s.store.UpdateGithubInstallationSync(ctx, installation.ID, now, replaceErr.Error())
		return replaceErr
	}
	_, err = s.store.UpdateGithubInstallationSync(ctx, installation.ID, now, "")
	return err
}
