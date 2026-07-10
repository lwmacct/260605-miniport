package server

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/appmodule"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/identity"
	"github.com/uptrace/bun"

	"github.com/lwmacct/260605-miniport/internal/handler"
	"github.com/lwmacct/260605-miniport/internal/infra/dbschema"
	"github.com/lwmacct/260605-miniport/internal/repository"
	"github.com/lwmacct/260605-miniport/internal/service"
)

type PortsvcModule struct {
	identity identity.SessionResolver
	value    *service.PortsvcService
}

var _ appmodule.Module = (*PortsvcModule)(nil)

func NewPortsvcSpec() appmodule.Spec {
	module := &PortsvcModule{}
	return appmodule.Spec{
		Name:     module.Name(),
		Requires: []string{"auth", "github"},
		Schema:   applyPortsvcSchema,
		Build: func(ctx *appmodule.Context) (appmodule.Module, error) {
			store := repository.NewStore(ctx.DB())
			authModule := appmodule.MustContextGet[*AuthModule](ctx, "auth")
			return &PortsvcModule{
				identity: authModule,
				value:    service.NewPortsvcService(store, authModule),
			}, nil
		},
	}
}

func (m *PortsvcModule) Name() string {
	return "portsvc"
}

func (m *PortsvcModule) Register(api huma.API) {
	handler.RegisterPortsvc(api, handler.Config{
		Identity: m.identity,
	}, handler.Services{
		Portsvc: m.value,
	})
}

func applyPortsvcSchema(ctx context.Context, db *bun.DB) error {
	if err := dbschema.Apply(ctx, db, repository.PortsvcSchema(), repository.PortsvcIndexesSchema()); err != nil {
		return fmt.Errorf("apply portsvc schema: %w", err)
	}
	return nil
}
