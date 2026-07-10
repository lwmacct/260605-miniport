package server

import (
	"github.com/lwmacct/260614-go-pkg-tlsreload/pkg/tlsreload"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/appmodule"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/requestctx"
	"github.com/uptrace/bun"

	"github.com/lwmacct/260605-miniport/internal/config"
)

type App struct {
	cfg  *config.Config
	deps *dependencies
}

type dependencies struct {
	db          *bun.DB
	modules     *appmodule.Runtime
	auth        *AuthModule
	github      *GithubModule
	portsvc     *PortsvcModule
	requests    requestctx.Middleware
	tlsReloader *tlsreload.Manager
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}
