package config

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/lwmacct/260711-go-pkg-authme/pkg/authme/adapters/statictoken"
	"github.com/stretchr/testify/require"
)

func TestNormalizeAuthMeDefaults(t *testing.T) {
	cfg := validTestAuthMeConfig()

	normalized, err := NormalizeAuthMe(cfg)
	require.NoError(t, err)
	require.Equal(t, "/authme", normalized.PathPrefix)
	require.Equal(t, []string{
		"http://localhost:40238",
		"http://localhost:40239",
	}, normalized.Origins)
	require.Len(t, normalized.StaticToken.Credentials, 1)
	require.Equal(t, "lwmacct", normalized.AllowedGitHubUser)
}

func TestNormalizeAuthMeRequiresOneStaticCredential(t *testing.T) {
	cfg := validTestAuthMeConfig()
	cfg.StaticToken.Credentials = append(cfg.StaticToken.Credentials, statictoken.Credential{
		ID: "second", Name: "Second operator", Token: "second-access-token",
	})

	_, err := NormalizeAuthMe(cfg)
	require.EqualError(t, err, "normalize authme: exactly one static token credential is required")
}

func validTestAuthMeConfig() ServerAuthMe {
	cfg := DefaultConfig().Server.HTTP.AuthMe
	cfg.Session.Keys[0].Secret = base64.RawURLEncoding.EncodeToString([]byte(strings.Repeat("k", 32)))
	cfg.StaticToken.Credentials[0].Token = "test-access-token"
	cfg.DexGitHub.ClientSecret = "test-client-secret"
	return cfg
}
