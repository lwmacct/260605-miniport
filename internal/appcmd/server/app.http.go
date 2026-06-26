package server

import (
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lwmacct/260605-miniport/internal/handler"
	"github.com/lwmacct/260605-miniport/internal/infra/frontend"
)

func (app *App) newHTTPServer() *http.Server {
	router := chi.NewRouter()
	apiHandler := http.StripPrefix("/api", app.newHTTPAPIHandler())
	router.Handle("/api", apiHandler)
	router.Handle("/api/*", apiHandler)

	if !frontend.RegisterRoutes(router, app.cfg.Server.HTTP.WebRoot) {
		router.Get("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message":"miniport backend is running","api":"/api"}`))
		})
	}

	httpHandler := app.container.requests.Wrap(router)
	if maxBodyBytes := app.cfg.Server.HTTP.MaxAPIBodyBytes; maxBodyBytes > 0 {
		httpHandler = http.MaxBytesHandler(httpHandler, maxBodyBytes)
	}

	srv := &http.Server{
		Addr:         app.cfg.Server.HTTP.Listen,
		Handler:      httpHandler,
		ReadTimeout:  app.cfg.Server.HTTP.ReadTimeout,
		WriteTimeout: app.cfg.Server.HTTP.WriteTimeout,
		IdleTimeout:  app.cfg.Server.HTTP.IdleTimeout,
	}
	if app.tlsManager != nil {
		srv.TLSConfig = app.tlsManager.TLSConfig(tls.VersionTLS12)
	}
	return srv
}

func (app *App) newHTTPAPIHandler() http.Handler {
	endpoint := handler.NewEndpoint(app.handlerConfig(), app.handlerServices())
	return endpoint.Handler()
}

func (app *App) handlerConfig() handler.Config {
	return handler.Config{
		LocalLoginEnabled:        app.cfg.Server.Auth.Local.LoginEnabled,
		LocalRegistrationEnabled: app.cfg.Server.Auth.Local.RegistrationEnabled,
		SecureCookies:            app.cfg.Server.HTTP.TLS.Enabled,
		RuntimeAdmins:            app.cfg.Server.Auth.Admins,
		Request:                  handler.RequestFromContext,
	}
}

func (app *App) isRuntimeAdminUsername(username string) bool {
	for _, admin := range app.cfg.Server.Auth.Admins {
		if strings.EqualFold(strings.TrimSpace(admin), username) {
			return true
		}
	}
	return false
}

func (app *App) handlerServices() handler.Services {
	return handler.Services{
		Users:      app.container.users,
		Passwords:  app.container.passwords,
		Sessions:   app.container.sessions,
		Challenges: app.container.challenges,
		AdminUsers: app.container.adminUsers,
		Inventory:  app.container.inventory,
	}
}
