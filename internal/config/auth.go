package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lwmacct/260711-go-pkg-authme/pkg/authme"
	"github.com/lwmacct/260711-go-pkg-authme/pkg/authme/adapters/dexgithub"
)

func NormalizeAuthMe(cfg ServerAuthMe) (ServerAuthMe, error) {
	core, err := (authme.Config{Prefix: cfg.PathPrefix, Origins: cfg.Origins, Session: cfg.Session}).Normalize()
	if err != nil {
		return cfg, fmt.Errorf("normalize authme: %w", err)
	}
	cfg.PathPrefix = core.Prefix
	cfg.Origins = core.Origins
	cfg.Session = core.Session

	if cfg.StaticToken.Enabled {
		if len(cfg.StaticToken.Credentials) != 1 {
			return cfg, errors.New("normalize authme: exactly one static token credential is required")
		}
		normalized, err := cfg.StaticToken.Normalize()
		if err != nil {
			return cfg, fmt.Errorf("normalize authme static token: %w", err)
		}
		cfg.StaticToken = normalized
	}

	if cfg.DexGitHub.Enabled {
		normalized, err := cfg.DexGitHub.Normalize()
		if err != nil {
			return cfg, fmt.Errorf("normalize authme dexgithub: %w", err)
		}
		cfg.DexGitHub = normalized
		cfg.AllowedGitHubUser = strings.ToLower(strings.TrimSpace(cfg.AllowedGitHubUser))
		if _, err := dexgithub.NewUsernameAuthorizer([]string{cfg.AllowedGitHubUser}); err != nil {
			return cfg, fmt.Errorf("normalize authme GitHub user: %w", err)
		}
	}

	if !cfg.StaticToken.Enabled && !cfg.DexGitHub.Enabled {
		return cfg, errors.New("normalize authme: at least one authentication method must be enabled")
	}
	return cfg, nil
}
