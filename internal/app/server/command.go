package server

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/lwmacct/260605-miniport/internal/config"
)

var (
	defaults = config.DefaultConfig()
	usage    = cfgm.Schema(defaults).Command("server")
)

// Command 是 Web API 服务命令。
var Command = &cli.Command{
	Name:            "server",
	Usage:           "启动 Web API 服务",
	Action:          action,
	Commands:        []*cli.Command{version.Command},
	HideHelpCommand: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "http.listen",
			Usage: usage.MustUsage("http.listen"),
			Value: defaults.Server.HTTP.Listen,
		},
		&cli.StringFlag{
			Name:  "http.ssl-cert-file",
			Usage: usage.MustUsage("http.ssl-cert-file"),
			Value: defaults.Server.HTTP.SSLCertFile,
		},
		&cli.StringFlag{
			Name:  "http.ssl-key-file",
			Usage: usage.MustUsage("http.ssl-key-file"),
			Value: defaults.Server.HTTP.SSLKeyFile,
		},
		&cli.StringFlag{
			Name:  "db.type",
			Usage: usage.MustUsage("db.type"),
			Value: defaults.Server.DB.Type,
		},
		&cli.StringFlag{
			Name:  "db.sqlite",
			Usage: usage.MustUsage("db.sqlite"),
			Value: defaults.Server.DB.SQLite,
		},
		&cli.StringFlag{
			Name:  "db.pgsql.host",
			Usage: usage.MustUsage("db.pgsql.host"),
			Value: defaults.Server.DB.PGSQL.Host,
		},
		&cli.StringFlag{
			Name:  "db.pgsql.port",
			Usage: usage.MustUsage("db.pgsql.port"),
			Value: defaults.Server.DB.PGSQL.Port,
		},
		&cli.StringFlag{
			Name:  "db.pgsql.user",
			Usage: usage.MustUsage("db.pgsql.user"),
			Value: defaults.Server.DB.PGSQL.User,
		},
		&cli.StringFlag{
			Name:  "db.pgsql.database",
			Usage: usage.MustUsage("db.pgsql.database"),
			Value: defaults.Server.DB.PGSQL.Database,
		},
		&cli.StringFlag{
			Name:  "db.pgsql.password",
			Usage: usage.MustUsage("db.pgsql.password"),
			Value: defaults.Server.DB.PGSQL.Password,
		},
		&cli.StringFlag{
			Name:  "control.listen",
			Usage: usage.MustUsage("control.listen"),
			Value: defaults.Server.Control.Listen,
		},
	},
}

func action(ctx context.Context, cmd *cli.Command) error {
	cfg := cfgm.MustLoadCmd(cmd, config.DefaultConfig(), "")
	return Run(ctx, cfg)
}
