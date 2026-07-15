package server

import "github.com/lwmacct/260605-miniport/internal/config"

func (app *App) validateHTTPTLS() error {
	cfg := app.cfg.Server.HTTP.TLS
	if !cfg.Enabled {
		return nil
	}
	return validateHTTPTLSRefs(cfg)
}

func validateHTTPTLSRefs(cfg config.ServerHTTPTLS) error {
	return cfg.ReloadConfig().Validate()
}
