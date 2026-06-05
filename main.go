package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"

	appserver "github.com/lwmacct/webapp/internal/app/server"
	"github.com/lwmacct/webapp/internal/config"
	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm"
)

func buildCommands() []*cli.Command {
	return []*cli.Command{
		serveCommand(),
		version.Command,
	}
}

func serveCommand() *cli.Command {
	defaults := config.DefaultConfig()
	return &cli.Command{
		Name:            "serve",
		Usage:           "启动 Web API 服务",
		Action:          serveAction,
		Commands:        []*cli.Command{version.Command},
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "http.listen",
				Usage: "HTTP 服务监听地址",
				Value: defaults.Server.HTTP.Listen,
			},
			&cli.StringFlag{
				Name:  "database",
				Usage: "SQLite 数据库文件路径",
				Value: defaults.Server.Database,
			},
		},
	}
}

func serveAction(ctx context.Context, cmd *cli.Command) error {
	cfg := cfgm.MustLoadCmd(cmd, config.DefaultConfig(), version.AppVersion)
	return appserver.Run(ctx, cfg)
}

func main() {
	logm.MustInit(logm.PresetAuto())

	cmd := &cli.Command{
		Name:            "webapp",
		Usage:           "Web app skeleton",
		Version:         version.AppVersion,
		Commands:        buildCommands(),
		HideHelpCommand: true,
		Action: func(ctx context.Context, c *cli.Command) error {
			return cli.ShowSubcommandHelp(c)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("command failed", "error", err)
		os.Exit(1)
	}
}

