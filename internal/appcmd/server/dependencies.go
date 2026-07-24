package server

import (
	"context"
	"fmt"

	"github.com/lwmacct/260630-go-hsr-shared/pkg/appmodule"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/database"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/requestctx"

	"github.com/lwmacct/260605-miniport/internal/config"
)

func newDependencies(ctx context.Context, cfg *config.Config) (*dependencies, error) {
	authCfg, err := config.NormalizeAuthMe(cfg.Server.HTTP.AuthMe)
	if err != nil {
		return nil, fmt.Errorf("configure authentication: %w", err)
	}
	cfg.Server.HTTP.AuthMe = authCfg
	authRuntime, err := newAccessAuth(ctx, authCfg)
	if err != nil {
		return nil, fmt.Errorf("configure authentication: %w", err)
	}

	db, err := database.Open(ctx, databaseConfig(cfg.Server.Database))
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	modules, err := appmodule.Build(ctx, db,
		NewCoreSpec(),
		NewGithubSpec(ctx, cfg),
		NewPortsvcSpec(),
	)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return &dependencies{
		db:       db,
		modules:  modules,
		auth:     authRuntime,
		github:   appmodule.MustGet[*GithubModule](modules, "github"),
		portsvc:  appmodule.MustGet[*PortsvcModule](modules, "portsvc"),
		requests: requestctx.NewMiddleware(cfg.Server.HTTP.TrustedProxies),
	}, nil
}

func (d *dependencies) Close() {
	if d == nil {
		return
	}
	if d.tlsStore != nil {
		_ = d.tlsStore.Close()
		d.tlsStore = nil
	}
	if d.modules != nil {
		_ = d.modules.Close()
		d.modules = nil
	}
	if d.db != nil {
		_ = d.db.Close()
		d.db = nil
	}
	d.auth = nil
	d.github = nil
	d.portsvc = nil
}

func databaseConfig(cfg config.ServerDatabase) database.Config {
	return database.Config{
		Type:   cfg.Type,
		SQLite: cfg.SQLite,
		PGSQL: database.PGSQLConfig{
			Host:     cfg.PGSQL.Host,
			Port:     cfg.PGSQL.Port,
			User:     cfg.PGSQL.User,
			Database: cfg.PGSQL.Database,
			Password: cfg.PGSQL.Password,
		},
	}
}
