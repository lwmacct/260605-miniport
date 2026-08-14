package server

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
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
	handler.RegisterCore(api)
}
