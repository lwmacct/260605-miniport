package server

import (
	"context"
	"fmt"

	"github.com/lwmacct/260605-miniport/internal/adapter/authchallengeprovider"
	"github.com/lwmacct/260605-miniport/internal/adapter/httpauth"
	"github.com/lwmacct/260605-miniport/internal/config"
	"github.com/lwmacct/260605-miniport/internal/domain/authchallenge"
	"github.com/lwmacct/260605-miniport/internal/domain/authpassword"
	"github.com/lwmacct/260605-miniport/internal/domain/authsession"
	"github.com/lwmacct/260605-miniport/internal/domain/identityuser"
	"github.com/lwmacct/260605-miniport/internal/domain/inventory"
	"github.com/lwmacct/260605-miniport/internal/infra/database"
	"github.com/lwmacct/260605-miniport/internal/infra/dbschema"
)

func (app *App) bootstrap(ctx context.Context) error {
	if err := app.cfg.Server.HTTP.TLS.Validate(); err != nil {
		return err
	}

	db, err := database.Open(ctx, app.cfg.Server.Database)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	models := append([]any{}, inventory.Schema()...)
	models = append(models, identityuser.Schema()...)
	models = append(models, authpassword.Schema()...)
	models = append(models, authsession.Schema()...)
	if err := dbschema.Apply(ctx, db, models, inventory.IndexStatements()); err != nil {
		_ = db.Close()
		return fmt.Errorf("apply database schema: %w", err)
	}

	app.db = db
	app.inventory = inventory.NewService(db)
	app.users = identityuser.NewService(db)
	app.passwords = authpassword.NewService(db)
	app.sessions = authsession.NewService(db, app.cfg.Server.HTTP.SessionTTL)
	app.challenges = app.newChallengeService()
	app.httpAuth = httpauth.NewService(app.cfg.Server.HTTP.TLS.Enabled(), app.cfg.Server.HTTP.TrustedProxies)
	return nil
}

func (app *App) newChallengeService() *authchallenge.Service {
	cfg := app.cfg.Server.Auth.Challenge
	switch cfg.Provider {
	case authchallenge.ProviderHCaptcha:
		provider, err := authchallengeprovider.NewRemoteTokenProvider(
			authchallenge.ProviderHCaptcha,
			cfg.HCaptcha.SiteKey,
			cfg.HCaptcha.Secret,
			cfg.HCaptcha.VerifyURL,
		)
		if err == nil {
			return authchallenge.NewService(provider)
		}
	case authchallenge.ProviderTurnstile:
		provider, err := authchallengeprovider.NewRemoteTokenProvider(
			authchallenge.ProviderTurnstile,
			cfg.Turnstile.SiteKey,
			cfg.Turnstile.Secret,
			cfg.Turnstile.VerifyURL,
		)
		if err == nil {
			return authchallenge.NewService(provider)
		}
	}
	return authchallenge.NewService(authchallengeprovider.NewImageProvider(cfg.Image.MaxChallenges))
}

func (app *App) closeDatabase() {
	if app.db == nil {
		return
	}
	_ = app.db.Close()
	app.db = nil
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
