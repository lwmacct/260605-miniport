package server

import (
	"github.com/lwmacct/260605-miniport/internal/adapter/httpauth"
	"github.com/lwmacct/260605-miniport/internal/domain/authchallenge"
	"github.com/uptrace/bun"

	"github.com/lwmacct/260605-miniport/internal/config"
	"github.com/lwmacct/260605-miniport/internal/domain/authpassword"
	"github.com/lwmacct/260605-miniport/internal/domain/authsession"
	"github.com/lwmacct/260605-miniport/internal/domain/identityuser"
	"github.com/lwmacct/260605-miniport/internal/domain/inventory"
)

type App struct {
	cfg        *config.Config
	db         *bun.DB
	inventory  *inventory.Service
	users      *identityuser.Service
	passwords  *authpassword.Service
	sessions   *authsession.Service
	challenges *authchallenge.Service
	httpAuth   *httpauth.Service
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}
