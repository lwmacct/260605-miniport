package server

import (
	"context"
	"fmt"

	"github.com/lwmacct/260605-miniport/internal/config"
	"github.com/lwmacct/260605-miniport/internal/infra/captcha"
	"github.com/lwmacct/260605-miniport/internal/infra/database"
	"github.com/lwmacct/260605-miniport/internal/infra/dbschema"
	"github.com/lwmacct/260605-miniport/internal/repository"
	"github.com/lwmacct/260605-miniport/internal/service"
)

func (app *App) bootstrap(ctx context.Context) error {
	if err := app.cfg.Server.HTTP.TLS.Validate(); err != nil {
		return err
	}

	db, err := database.Open(ctx, app.cfg.Server.Database)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	models := append([]any{}, repository.InventorySchema()...)
	models = append(models, repository.IdentityUserSchema()...)
	models = append(models, repository.AuthPasswordSchema()...)
	models = append(models, repository.AuthSessionSchema()...)
	if err := dbschema.Apply(ctx, db, models, repository.InventoryIndexesSchema()); err != nil {
		_ = db.Close()
		return fmt.Errorf("apply database schema: %w", err)
	}

	app.container.db = db
	app.container.store = repository.NewStore(db)
	app.container.inventory = service.NewInventoryService(app.container.store)
	app.container.users = service.NewIdentityUserService(app.container.store)
	app.container.passwords = service.NewAuthPasswordService(app.container.store)
	app.container.sessions = service.NewAuthSessionService(app.container.store, app.cfg.Server.HTTP.SessionTTL)
	app.container.challenges = app.newChallengeService()
	app.container.adminUsers = service.NewAdminUserService(app.container.users, app.isRuntimeAdminUsername)
	app.container.requests = newRequestContextMiddleware(app.cfg.Server.HTTP.TrustedProxies)
	return nil
}

func (app *App) newChallengeService() *service.AuthChallengeService {
	cfg := app.cfg.Server.Auth.Challenge
	switch cfg.Provider {
	case service.AuthChallengeProviderHCaptcha:
		provider, err := captcha.NewRemoteTokenProvider(
			service.AuthChallengeProviderHCaptcha,
			cfg.HCaptcha.SiteKey,
			cfg.HCaptcha.Secret,
			cfg.HCaptcha.VerifyURL,
		)
		if err == nil {
			return service.NewAuthChallengeService(provider)
		}
	case service.AuthChallengeProviderTurnstile:
		provider, err := captcha.NewRemoteTokenProvider(
			service.AuthChallengeProviderTurnstile,
			cfg.Turnstile.SiteKey,
			cfg.Turnstile.Secret,
			cfg.Turnstile.VerifyURL,
		)
		if err == nil {
			return service.NewAuthChallengeService(provider)
		}
	}
	return service.NewAuthChallengeService(captcha.NewImageProvider(cfg.Image.MaxChallenges))
}

func (app *App) closeDatabase() {
	if app.container.db == nil {
		return
	}
	_ = app.container.db.Close()
	app.container.db = nil
	app.container.store = nil
}

func databaseDisplay(cfg config.ServerDatabase) string {
	if cfg.Type == "pgsql" {
		host := cfg.PGSQL.Host
		if host == "" {
			host = "localhost"
		}
		name := cfg.PGSQL.Database
		if name == "" {
			name = "postgres"
		}
		return "pgsql://" + host + "/" + name
	}
	return cfg.SQLite
}
