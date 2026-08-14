package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
)

type coreHandler struct{}

func RegisterCore(api huma.API) {
	h := coreHandler{}
	huma.Register(api, huma.Operation{OperationID: "get-health", Method: http.MethodGet, Path: "/health", Summary: "Get service health", Tags: []string{"System"}}, h.health)
	huma.Register(api, huma.Operation{OperationID: "get-meta", Method: http.MethodGet, Path: "/meta", Summary: "Get service metadata", Tags: []string{"System"}}, h.meta)
}

func (coreHandler) health(_ context.Context, _ *struct{}) (*BodyDTO[HealthDTO], error) {
	return &BodyDTO[HealthDTO]{Body: HealthDTO{Status: "ok", Timestamp: time.Now().UTC(), Version: version.AppVersion}}, nil
}

func (coreHandler) meta(_ context.Context, _ *struct{}) (*BodyDTO[MetaDTO], error) {
	return &BodyDTO[MetaDTO]{Body: MetaDTO{Name: "Miniport", Version: version.AppVersion, DocsPath: "/api"}}, nil
}
