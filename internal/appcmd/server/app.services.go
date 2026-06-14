package server

import (
	"context"
	"fmt"

	"github.com/lwmacct/260605-miniport/internal/config"
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
	if err := dbschema.Apply(ctx, db, inventory.Schema(), inventory.IndexStatements()); err != nil {
		_ = db.Close()
		return fmt.Errorf("apply database schema: %w", err)
	}

	app.db = db
	app.inventory = inventory.NewService(db)
	return nil
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
