package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm"
	appserver "github.com/lwmacct/260605-miniport/internal/app/server"
	"github.com/lwmacct/260605-miniport/internal/config"
)

func buildCommands() []*cli.Command {
	return []*cli.Command{
		appserver.Command,
		controlCommand(),
		version.Command,
	}
}

func controlCommand() *cli.Command {
	defaults := config.DefaultConfig()
	return &cli.Command{
		Name:  "control",
		Usage: "向本地控制面发送命令",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "socket",
				Usage: "控制面 Unix socket 路径",
				Value: defaults.Server.Control.Listen,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "reload-cert",
				Usage: "重载 TLS 证书",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					socket := cmd.String("socket")
					if socket == "" {
						return fmt.Errorf("socket is required")
					}
					conn, err := net.Dial("unix", socket)
					if err != nil {
						return err
					}
					defer conn.Close()

					if _, err := io.WriteString(conn, "reload-cert\n"); err != nil {
						return err
					}
					resp, err := io.ReadAll(conn)
					if err != nil {
						return err
					}
					out := strings.TrimSpace(string(resp))
					if strings.HasPrefix(out, "ERR ") {
						return fmt.Errorf("%s", strings.TrimPrefix(out, "ERR "))
					}
					slog.Info("control command completed", "response", out)
					return nil
				},
			},
		},
	}
}

func main() {
	logm.MustInit(logm.PresetAuto())

	cmd := &cli.Command{
		Name:            "webapp",
		Usage:           "Web app skeleton",
		Version:         version.AppVersion,
		Flags:           []cli.Flag{cfgm.ConfigFlag()},
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
