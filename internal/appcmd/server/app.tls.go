package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/lwmacct/260614-go-pkg-tlsreload/pkg/tlsreload"
)

func (app *App) bootstrapTLSManager(ctx context.Context) error {
	if !app.cfg.Server.HTTP.TLS.Enabled {
		return nil
	}

	httpTLS := app.cfg.Server.HTTP.TLS
	manager, err := tlsreload.NewFileManager(ctx, tlsreload.FileManagerOptions{
		CertFile:       httpTLS.CertFile,
		KeyFile:        httpTLS.KeyFile,
		AutoReload:     httpTLS.AutoReload,
		ReloadInterval: httpTLS.ReloadInterval,
		RetryInterval:  2 * time.Second,
		Logger:         slog.Default(),
	})
	if err != nil {
		return err
	}

	app.tlsManager = manager
	return nil
}

func (app *App) closeTLSManager() {
	if app.tlsManager == nil {
		return
	}
	app.tlsManager.Close()
	app.tlsManager = nil
}
