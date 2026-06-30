package server

import (
	"context"
	"crypto/tls"
	"log/slog"
	"time"

	"github.com/lwmacct/260614-go-pkg-tlsreload/pkg/tlsreload"
)

func (app *App) bootstrapTLSManager(ctx context.Context) error {
	if !app.cfg.Server.HTTP.TLS.Enabled {
		return nil
	}

	httpTLS := app.cfg.Server.HTTP.TLS
	reloader, err := tlsreload.New(ctx, tlsreload.Config{
		CertFile:       httpTLS.CertFile,
		KeyFile:        httpTLS.KeyFile,
		ReloadInterval: httpTLS.ReloadInterval,
		MinVersion:     tls.VersionTLS12,
		RetryInterval:  2 * time.Second,
		Logger:         slog.Default(),
	})
	if err != nil {
		return err
	}

	app.deps.tlsReloader = reloader
	return nil
}
