package server

import (
	"context"

	"github.com/lwmacct/260711-go-pkg-authme/pkg/authme"
	"github.com/lwmacct/260711-go-pkg-authme/pkg/authme/adapters/dexgithub"
	"github.com/lwmacct/260711-go-pkg-authme/pkg/authme/adapters/statictoken"

	"github.com/lwmacct/260605-miniport/internal/config"
)

func newAccessAuth(ctx context.Context, cfg config.ServerAuthMe) (*authme.Auth, error) {
	methods := make([]authme.Method, 0, 2)
	var authorizers []authme.Authorizer
	if cfg.StaticToken.Enabled {
		method, err := statictoken.New(cfg.StaticToken)
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}
	if cfg.DexGitHub.Enabled {
		method, err := dexgithub.New(ctx, cfg.DexGitHub)
		if err != nil {
			return nil, err
		}
		authorizer, err := dexgithub.NewUsernameAuthorizer([]string{cfg.AllowedGitHubUser})
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
		authorizers = append(authorizers, authorizer)
	}
	return authme.New(
		authme.Config{Prefix: cfg.PathPrefix, Origins: cfg.Origins, Session: cfg.Session},
		authme.WithMethods(methods...),
		authme.WithAuthorizer(authme.Chain(authorizers...)),
	)
}
