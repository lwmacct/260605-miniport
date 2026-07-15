package server

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/lwmacct/260605-miniport/internal/config"
)

var binding = config.Definition.Bind(cfgm.Scope("server"))

var Command = &cli.Command{
	Name:            "server",
	Usage:           "start HTTP server",
	HideHelpCommand: true,
	Commands:        []*cli.Command{version.Command},
	Action:          action,
	Flags:           binding.Flags(),
}

func action(ctx context.Context, cmd *cli.Command) error {
	cfg := binding.MustLoad(ctx, cmd)
	return NewApp(cfg).Run(ctx)
}
