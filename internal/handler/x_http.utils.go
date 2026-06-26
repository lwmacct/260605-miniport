package handler

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
)

func utilHTTPConfig() huma.Config {
	cfg := huma.DefaultConfig("Miniport API", version.AppVersion)
	cfg.Servers = []*huma.Server{{URL: "/api"}}
	return cfg
}
