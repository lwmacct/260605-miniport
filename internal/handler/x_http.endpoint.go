package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
)

type Endpoint struct {
	config   Config
	services Services
}

func NewEndpoint(config Config, services Services) *Endpoint {
	return &Endpoint{config: config, services: services}
}

func (e *Endpoint) Handler() http.Handler {
	mux := http.NewServeMux()
	api := humago.New(mux, utilHTTPConfig())
	e.Register(api)
	return mux
}

func (e *Endpoint) Register(api huma.API) {
	RegisterAuth(api, e.config, e.services)
	RegisterAdminUser(api, e.config, e.services)
	RegisterInventory(api, e.config, e.services)
	huma.Register(api, huma.Operation{OperationID: "get-health", Method: http.MethodGet, Path: "/health", Summary: "Get service health"}, e.health)
	huma.Register(api, huma.Operation{OperationID: "get-meta", Method: http.MethodGet, Path: "/meta", Summary: "Get service metadata"}, e.meta)
}

func (e *Endpoint) health(_ context.Context, _ *struct{}) (*BodyDTO[HealthDTO], error) {
	return &BodyDTO[HealthDTO]{Body: HealthDTO{Status: "ok", Timestamp: time.Now().UTC(), Version: version.AppVersion}}, nil
}

func (e *Endpoint) meta(_ context.Context, _ *struct{}) (*BodyDTO[MetaDTO], error) {
	return &BodyDTO[MetaDTO]{Body: MetaDTO{Name: "Miniport", Version: version.AppVersion, DocsPath: "/api"}}, nil
}
