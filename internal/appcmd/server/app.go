package server

import (
	"github.com/lwmacct/260614-go-pkg-tlsreload/pkg/tlsreload"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/appmodule"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/requestctx"
	"github.com/lwmacct/260711-go-pkg-authme/pkg/authme"
	"github.com/uptrace/bun"

	"github.com/lwmacct/260605-miniport/internal/config"
)

type App struct {
	cfg  *config.Config
	deps *dependencies
}

type dependencies struct {
	db       *bun.DB
	modules  *appmodule.Runtime
	auth     *authme.Auth
	github   *GithubModule
	portsvc  *PortsvcModule
	requests requestctx.Middleware
	tlsStore *tlsreload.Store
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}
