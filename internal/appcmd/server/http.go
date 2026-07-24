package server

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/go-chi/chi/v5"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/httpserver"

	"github.com/lwmacct/260605-miniport/internal/infra/frontend"
)

const httpAPIPrefix = "/api"

func (app *App) newHTTPServer() *http.Server {
	httpCfg := app.cfg.Server.HTTP
	srv := &http.Server{
		Addr:         httpCfg.Listen,
		Handler:      app.newHTTPHandler(),
		ReadTimeout:  httpCfg.ReadTimeout,
		WriteTimeout: httpCfg.WriteTimeout,
		IdleTimeout:  httpCfg.IdleTimeout,
	}
	if app.deps.tlsStore != nil {
		srv.TLSConfig = app.tlsConfig()
	}
	return srv
}

func (app *App) newHTTPHandler() http.Handler {
	router := chi.NewRouter()
	apiHandler := http.StripPrefix(httpAPIPrefix, app.newHTTPAPIHandler())
	protectedAPI := app.deps.auth.RequireAccess(apiHandler)
	for _, publicPath := range []string{
		httpAPIPrefix + "/health",
		httpAPIPrefix + "/meta",
		httpAPIPrefix + "/github/setup",
		httpAPIPrefix + "/github/webhooks",
	} {
		router.Handle(publicPath, apiHandler)
	}
	router.Handle(httpAPIPrefix, protectedAPI)
	router.Handle(httpAPIPrefix+"/*", protectedAPI)
	authPrefix := app.deps.auth.PathPrefix()
	router.Handle(authPrefix, app.deps.auth.Handler())
	router.Handle(authPrefix+"/*", app.deps.auth.Handler())

	if !frontend.RegisterRoutes(router, app.cfg.Server.HTTP.WebRoot) {
		router.Get("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message":"miniport backend is running","api":"/api"}`))
		})
	}

	return app.deps.requests.Wrap(router)
}

func (app *App) newHTTPAPIHandler() http.Handler {
	maxBodyBytes := max(app.cfg.Server.HTTP.MaxAPIBodyBytes, 0)
	mux := http.NewServeMux()
	api := humago.New(mux, httpAPIConfig())
	app.deps.modules.Register(api)
	return httpserver.LimitRequestBody(mux, maxBodyBytes)
}

func httpAPIConfig() huma.Config {
	cfg := huma.DefaultConfig("Miniport API", version.AppVersion)
	cfg.Servers = []*huma.Server{{URL: httpAPIPrefix}}
	return cfg
}
