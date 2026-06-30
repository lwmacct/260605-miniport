package server

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/appmodule"
	"github.com/uptrace/bun"

	"github.com/lwmacct/260605-miniport/internal/handler"
)

type CoreModule struct{}

var _ appmodule.Module = (*CoreModule)(nil)

func NewCoreSpec() appmodule.Spec {
	module := &CoreModule{}
	return appmodule.Spec{
		Name:   module.Name(),
		Schema: func(context.Context, *bun.DB) error { return nil },
		Build: func(*appmodule.Context) (appmodule.Module, error) {
			return module, nil
		},
	}
}

func (m *CoreModule) Name() string {
	return "core"
}

func (m *CoreModule) Register(api huma.API) {
	huma.Register(api, huma.Operation{OperationID: "get-health", Method: http.MethodGet, Path: "/health", Summary: "Get service health"}, m.health)
	huma.Register(api, huma.Operation{OperationID: "get-meta", Method: http.MethodGet, Path: "/meta", Summary: "Get service metadata"}, m.meta)
}

func (m *CoreModule) health(_ context.Context, _ *struct{}) (*handler.BodyDTO[handler.HealthDTO], error) {
	return &handler.BodyDTO[handler.HealthDTO]{Body: handler.HealthDTO{Status: "ok", Timestamp: time.Now().UTC(), Version: version.AppVersion}}, nil
}

func (m *CoreModule) meta(_ context.Context, _ *struct{}) (*handler.BodyDTO[handler.MetaDTO], error) {
	return &handler.BodyDTO[handler.MetaDTO]{Body: handler.MetaDTO{Name: "Miniport", Version: version.AppVersion, DocsPath: "/api"}}, nil
}
