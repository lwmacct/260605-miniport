package server

import (
	"github.com/uptrace/bun"

	"github.com/lwmacct/260605-miniport/internal/config"
	"github.com/lwmacct/260605-miniport/internal/domain/inventory"
)

type App struct {
	cfg       *config.Config
	db        *bun.DB
	inventory *inventory.Service
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}
