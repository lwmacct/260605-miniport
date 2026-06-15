package database

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"

	"github.com/lwmacct/260605-miniport/internal/config"
)

func Open(ctx context.Context, cfg config.ServerDatabase) (*bun.DB, error) {
	switch strings.ToLower(cfg.Type) {
	case "", "sqlite":
		return openSQLite(ctx, cfg.SQLite)
	case "pgsql":
		return openPGSQL(ctx, cfg.PGSQL)
	default:
		return nil, fmt.Errorf("unsupported database type %q", cfg.Type)
	}
}

func openSQLite(ctx context.Context, sqlitePath string) (*bun.DB, error) {
	if err := os.MkdirAll(filepath.Dir(sqlitePath), 0o750); err != nil {
		return nil, fmt.Errorf("prepare sqlite directory: %w", err)
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, sqlitePath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	if err := sqldb.PingContext(ctx); err != nil {
		_ = sqldb.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	return bun.NewDB(sqldb, sqlitedialect.New()), nil
}

func openPGSQL(ctx context.Context, cfg config.ServerDatabasePGSQL) (*bun.DB, error) {
	host := cfg.Host
	if host == "" {
		host = "localhost"
	}
	port := cfg.Port
	if port == "" {
		port = "5432"
	}

	opts := []pgdriver.Option{
		pgdriver.WithInsecure(true),
		pgdriver.WithAddr(net.JoinHostPort(host, port)),
	}
	if cfg.User != "" {
		opts = append(opts, pgdriver.WithUser(cfg.User))
	}
	if cfg.Database != "" {
		opts = append(opts, pgdriver.WithDatabase(cfg.Database))
	}
	if cfg.Password != "" {
		opts = append(opts, pgdriver.WithPassword(cfg.Password))
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(opts...))
	if err := sqldb.PingContext(ctx); err != nil {
		_ = sqldb.Close()
		return nil, fmt.Errorf("ping pgsql database: %w", err)
	}

	return bun.NewDB(sqldb, pgdialect.New()), nil
}
