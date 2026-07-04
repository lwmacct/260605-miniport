package server

import (
	"github.com/lwmacct/260614-go-pkg-tlsreload/pkg/tlsreload"

	"github.com/lwmacct/260605-miniport/internal/config"
)

func (app *App) validateHTTPTLS() error {
	cfg := app.cfg.Server.HTTP.TLS
	if !cfg.Enabled {
		return nil
	}
	return validateHTTPTLSRefs(cfg)
}

func validateHTTPTLSRefs(cfg config.ServerHTTPTLS) error {
	return tlsreload.Config{
		Enabled:      cfg.Enabled,
		CertFile:     cfg.CertFile,
		KeyFile:      cfg.KeyFile,
		PollInterval: cfg.PollInterval,
	}.Validate()
}
