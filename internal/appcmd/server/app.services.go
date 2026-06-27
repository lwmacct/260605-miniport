package server

import (
	"context"
	"fmt"

	"github.com/lwmacct/260605-miniport/internal/infra/captcha"
	"github.com/lwmacct/260605-miniport/internal/infra/database"
	"github.com/lwmacct/260605-miniport/internal/infra/dbschema"
	"github.com/lwmacct/260605-miniport/internal/repository"
	"github.com/lwmacct/260605-miniport/internal/service"
)

func (app *App) bootstrap(ctx context.Context) error {
	db, err := database.Open(ctx, app.cfg.Server.Database)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	models := append([]any{}, repository.UserSchema()...)
	models = append(models, repository.InventorySchema()...)
	models = append(models, repository.AuthPasswordSchema()...)
	models = append(models, repository.AuthSessionSchema()...)
	if err := dbschema.Apply(ctx, db, models, repository.InventoryIndexesSchema()); err != nil {
		_ = db.Close()
		return fmt.Errorf("apply database schema: %w", err)
	}

	app.container.db = db
	app.container.store = repository.NewStore(db)
	app.container.inventory = service.NewInventoryService(app.container.store)
	app.container.users = service.NewUserService(app.container.store)
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
