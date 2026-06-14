package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *App) Run(ctx context.Context) error {
	if err := app.bootstrap(ctx); err != nil {
		return err
	}
	defer app.closeDatabase()

	srv, err := app.newHTTPServer()
	if err != nil {
		return err
	}

	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", srv.Addr)
	if err != nil {
		return err
	}
	defer func() { _ = ln.Close() }()

	errCh := make(chan error, 1)
	go func() {
		cfg := app.cfg.Server.HTTP
		slog.Info("web service starting", "listen", srv.Addr, "https", cfg.TLS.Enabled(), "web_root", cfg.WebRoot)
		var serveErr error
		if cfg.TLS.Enabled() {
			serveErr = srv.ServeTLS(ln, cfg.TLS.CertFile, cfg.TLS.KeyFile)
		} else {
			serveErr = srv.Serve(ln)
		}
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errCh <- serveErr
		}
		close(errCh)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(sigCh)

	select {
	case <-ctx.Done():
		return app.shutdown(ctx, srv)
	case sig := <-sigCh:
		slog.Info("received shutdown signal", "signal", sig.String())
		return app.shutdown(ctx, srv)
	case err := <-errCh:
		return err
	}
}

func (app *App) shutdown(ctx context.Context, srv *http.Server) error {
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	slog.Info("web service stopped")
	return nil
}
