package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lwmacct/260711-go-pkg-authme/pkg/authme/adapters/statictoken"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260605-miniport/internal/config"
)

func TestPortsvcWriteAllowsAuthenticatedUser(t *testing.T) {
	handler := setupPortsvcTestApp(t)
	cookie := createTestSessionCookie(t, handler)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://localhost/api/port-groups", strings.NewReader(`{"portPrefix":1000,"environmentName":"miniport","runtimeMode":"dind","runtimeName":"miniport-dind-01","serviceIp":"172.22.11.12","slots":[{"port":10000,"name":"redis","kind":"cache","protocol":"redis"}]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost")
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var body struct {
		EnvironmentName string `json:"environmentName"`
		Slots           []struct {
			Name string `json:"name"`
		} `json:"slots"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, "miniport", body.EnvironmentName)
	require.Len(t, body.Slots, 1)
	require.Equal(t, "redis", body.Slots[0].Name)
}

func TestPortsvcReadRejectsMissingSession(t *testing.T) {
	handler := setupPortsvcTestApp(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/port-groups", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code, rec.Body.String())
}

func TestGithubStatusUsesLocalSession(t *testing.T) {
	handler := setupPortsvcTestApp(t)
	cookie := createTestSessionCookie(t, handler)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/github/status", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var status struct {
		Enabled bool `json:"enabled"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &status))
	require.False(t, status.Enabled)

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/github/status", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code, rec.Body.String())
}

func TestGithubSetupDoesNotRequireSessionCookie(t *testing.T) {
	handler := setupPortsvcTestApp(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/github/setup?installation_id=42&state=one-time-state", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusServiceUnavailable, rec.Code, rec.Body.String())
}

func setupPortsvcTestApp(t *testing.T) http.Handler {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.Server.Database.SQLite = "file:" + strings.NewReplacer("/", "-", " ", "-").Replace(t.Name()) + "?mode=memory&cache=shared"
	cfg.Server.HTTP.WebRoot = ""
	cfg.Server.HTTP.AuthMe.Origins = []string{"http://localhost"}
	cfg.Server.HTTP.AuthMe.Session.Keys[0].Secret = base64.RawURLEncoding.EncodeToString([]byte(strings.Repeat("k", 32)))
	cfg.Server.HTTP.AuthMe.StaticToken.Credentials = []statictoken.Credential{{ID: "operator", Name: "Operator", Token: "test-access-token"}}
	app := NewApp(&cfg)
	deps, err := newDependencies(context.Background(), &cfg)
	require.NoError(t, err)
	app.deps = deps
	t.Cleanup(app.deps.Close)
	return app.newHTTPServer().Handler
}

func createTestSessionCookie(t *testing.T, handler http.Handler) *http.Cookie {
	t.Helper()

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://localhost/authme/login/token", strings.NewReader(`{"token":"test-access-token"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost")
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNoContent, rec.Code, rec.Body.String())
	require.NotEmpty(t, rec.Result().Cookies())
	return rec.Result().Cookies()[0]
}
