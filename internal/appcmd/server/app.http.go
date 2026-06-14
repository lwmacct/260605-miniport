package server

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"

	"github.com/lwmacct/260605-miniport/internal/adapter/inventoryhttp"
	"github.com/lwmacct/260605-miniport/internal/infra/frontend"
)

type healthOutput struct {
	Body struct {
		Status    string    `json:"status" example:"ok"`
		Timestamp time.Time `json:"timestamp" example:"2026-06-15T12:00:00Z"`
		Version   string    `json:"version" example:"0.1.0"`
		Database  string    `json:"database" example:".local/data/sqlite.db"`
	}
}

type metaOutput struct {
	Body struct {
		Name     string `json:"name" example:"Miniport"`
		Version  string `json:"version" example:"0.1.0"`
		Listen   string `json:"listen" example:":40238"`
		Database string `json:"database" example:".local/data/sqlite.db"`
		DocsPath string `json:"docsPath" example:"/api"`
	}
}

func (app *App) newHTTPServer() (*http.Server, error) {
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	router.Mount("/api", apiRouter)

	api := humachi.New(apiRouter, httpAPIConfig())
	app.registerHTTPAPI(api)

	if !frontend.RegisterRoutes(router, app.cfg.Server.HTTP.WebRoot) {
		router.Get("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message":"miniport backend is running","api":"/api"}`))
		})
	}

	handler := http.Handler(router)
	if maxBodyBytes := app.cfg.Server.HTTP.MaxAPIBodyBytes; maxBodyBytes > 0 {
		handler = http.MaxBytesHandler(handler, maxBodyBytes)
	}

	return &http.Server{
		Addr:         app.cfg.Server.HTTP.Listen,
		Handler:      handler,
		ReadTimeout:  app.cfg.Server.HTTP.ReadTimeout,
		WriteTimeout: app.cfg.Server.HTTP.WriteTimeout,
		IdleTimeout:  app.cfg.Server.HTTP.IdleTimeout,
	}, nil
}

func httpAPIConfig() huma.Config {
	cfg := huma.DefaultConfig("Miniport API", version.AppVersion)
	cfg.Servers = []*huma.Server{{URL: "/api"}}
	return cfg
}

func (app *App) registerHTTPAPI(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-health",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Get service health",
	}, func(ctx context.Context, input *struct{}) (*healthOutput, error) {
		_ = input
		out := &healthOutput{}
		out.Body.Status = "ok"
		out.Body.Timestamp = time.Now().UTC()
		out.Body.Version = version.AppVersion
		out.Body.Database = databaseDisplay(app.cfg.Server.Database)
		if err := app.db.DB.PingContext(ctx); err != nil {
			out.Body.Status = "degraded"
		}
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-meta",
		Method:      http.MethodGet,
		Path:        "/meta",
		Summary:     "Get service metadata",
	}, func(_ context.Context, input *struct{}) (*metaOutput, error) {
		_ = input
		out := &metaOutput{}
		out.Body.Name = "Miniport"
		out.Body.Version = version.AppVersion
		out.Body.Listen = app.cfg.Server.HTTP.Listen
		out.Body.Database = databaseDisplay(app.cfg.Server.Database)
		out.Body.DocsPath = "/api"
		return out, nil
	})

	inventoryhttp.Register(api, app.inventory)
}
