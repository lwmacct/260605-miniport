package server

import (
	"context"
	"crypto/tls"
	"log/slog"

	"github.com/lwmacct/260614-go-pkg-tlsreload/pkg/tlsreload"
)

func (app *App) bootstrapTLSStore(ctx context.Context) error {
	if !app.cfg.Server.HTTP.TLS.Enabled {
		return nil
	}

	httpTLS := app.cfg.Server.HTTP.TLS
	store, err := tlsreload.New(ctx, httpTLS.ReloadConfig(), tlsreload.WithLogger(slog.Default()))
	if err != nil {
		return err
	}

	app.deps.tlsStore = store
	return nil
}

func (app *App) tlsConfig() *tls.Config {
	return &tls.Config{
		MinVersion:     tls.VersionTLS12,
		GetCertificate: app.deps.tlsStore.GetCertificate,
	}
}
