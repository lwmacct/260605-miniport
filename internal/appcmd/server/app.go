package server

import (
	"github.com/lwmacct/260614-go-pkg-tlsreload/pkg/tlsreload"
	"github.com/uptrace/bun"

	"github.com/lwmacct/260605-miniport/internal/config"
	"github.com/lwmacct/260605-miniport/internal/repository"
	"github.com/lwmacct/260605-miniport/internal/service"
)

type App struct {
	cfg         *config.Config
	container   appContainer
	tlsReloader *tlsreload.Reloader
}

type appContainer struct {
	db         *bun.DB
	store      *repository.Store
	inventory  *service.InventoryService
	users      *service.UserService
	passwords  *service.AuthPasswordService
	sessions   *service.AuthSessionService
	challenges *service.AuthChallengeService
	adminUsers *service.AdminUserService
	requests   requestContextMiddleware
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}
