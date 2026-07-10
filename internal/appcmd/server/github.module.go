package server

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/appmodule"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/identity"
	"github.com/uptrace/bun"

	"github.com/lwmacct/260605-miniport/internal/config"
	"github.com/lwmacct/260605-miniport/internal/handler"
	"github.com/lwmacct/260605-miniport/internal/infra/dbschema"
	"github.com/lwmacct/260605-miniport/internal/repository"
	"github.com/lwmacct/260605-miniport/internal/service"
)

type GithubModule struct {
	identity identity.SessionResolver
	value    *service.GithubService
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

var _ appmodule.Module = (*GithubModule)(nil)
var _ appmodule.Closer = (*GithubModule)(nil)

func NewGithubSpec(parent context.Context, cfg *config.Config) appmodule.Spec {
	module := &GithubModule{}
	return appmodule.Spec{
		Name: module.Name(), Requires: []string{"auth"}, Schema: applyGithubSchema,
		Build: func(moduleCtx *appmodule.Context) (appmodule.Module, error) {
			authModule := appmodule.MustContextGet[*AuthModule](moduleCtx, "auth")
			githubService, err := service.NewGithubService(repository.NewStore(moduleCtx.DB()), githubServiceConfig(cfg.Server.GitHub))
			if err != nil {
				return nil, err
			}
			built := &GithubModule{identity: authModule, value: githubService}
			built.startReconciler(parent, cfg.Server.GitHub)
			return built, nil
		},
	}
}

func (m *GithubModule) Name() string { return "github" }

func (m *GithubModule) Register(api huma.API) {
	handler.RegisterGithub(api, handler.Config{Identity: m.identity}, m.value)
}

func (m *GithubModule) Close() error {
	if m.cancel != nil {
		m.cancel()
		m.wg.Wait()
	}
	return nil
}

func (m *GithubModule) startReconciler(parent context.Context, cfg config.ServerGitHub) {
	if !cfg.Enabled || cfg.ReconcileInterval <= 0 {
		return
	}
	ctx, cancel := context.WithCancel(parent)
	m.cancel = cancel
	m.wg.Go(func() {
		ticker := time.NewTicker(cfg.ReconcileInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := m.value.Reconcile(ctx); err != nil {
					slog.Error("github reconciliation failed", "error", err)
				}
			}
		}
	})
}

func applyGithubSchema(ctx context.Context, db *bun.DB) error {
	if err := dbschema.Apply(ctx, db, repository.GithubSchema(), repository.GithubIndexesSchema()); err != nil {
		return fmt.Errorf("apply github schema: %w", err)
	}
	return nil
}

func githubServiceConfig(cfg config.ServerGitHub) service.GithubConfig {
	return service.GithubConfig{
		Enabled: cfg.Enabled, AppID: cfg.AppID, AppSlug: cfg.AppSlug, PrivateKeyFile: cfg.PrivateKeyFile,
		WebhookSecret: cfg.WebhookSecret, APIURL: cfg.APIURL, WebURL: cfg.WebURL,
		SetupReturnURL: cfg.SetupReturnURL, ReconcileInterval: cfg.ReconcileInterval,
	}
}
