package server

import (
	"errors"

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
	if (cfg.CertFile == "") != (cfg.KeyFile == "") {
		return errors.New("http tls.cert-file and tls.key-file must be configured together")
	}
	if cfg.ReloadInterval < 0 {
		return errors.New("http tls.reload-interval must not be negative")
	}
	return nil
}
