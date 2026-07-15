package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/260605-miniport/internal/config"
)

func TestManagerLoadsTLSCertificateFlags(t *testing.T) {
	manager := cfgm.New(config.DefaultConfig(), cfgm.WithoutDefaultPaths())
	var loaded *config.Config
	command := &cli.Command{
		Name: "server",
		Action: manager.Action(func(_ context.Context, _ *cli.Command, cfg *config.Config) error {
			loaded = cfg
			return nil
		}),
	}
	root := &cli.Command{Name: "app", Commands: []*cli.Command{command}}
	manager.MustConfigure(root)

	err := root.Run(t.Context(), []string{
		"app", "server",
		"--http.tls.enabled",
		`--http.tls.certificates={"id":"main","certificate":"/certs/main.pem","private-key":"/certs/main.key"}`,
		`--http.tls.certificates={"id":"api","certificate":"op://vault/api-cert","private-key":"op://vault/api-key"}`,
		"--http.tls.default-certificate=main",
		"--http.tls.retry-interval=4s",
	})
	require.NoError(t, err)
	require.NotNil(t, loaded)
	require.True(t, loaded.Server.HTTP.TLS.Enabled)
	require.Equal(t, "main", loaded.Server.HTTP.TLS.DefaultCertificate)
	require.Len(t, loaded.Server.HTTP.TLS.Certificates, 2)
	require.Equal(t, "api", loaded.Server.HTTP.TLS.Certificates[1].ID)
	require.Equal(t, "op://vault/api-key", loaded.Server.HTTP.TLS.Certificates[1].PrivateKey)
	require.Equal(t, "4s", loaded.Server.HTTP.TLS.RetryInterval.String())
}
