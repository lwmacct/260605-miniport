package server

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/lwmacct/260605-miniport/internal/config"
)

var Command = &cli.Command{
	Name:            "server",
	Usage:           "start HTTP server",
	HideHelpCommand: true,
	Commands:        []*cli.Command{version.Command},
	Action:          config.Manager.Action(action),
}

func action(ctx context.Context, _ *cli.Command, cfg *config.Config) error {
	return NewApp(cfg).Run(ctx)
}
