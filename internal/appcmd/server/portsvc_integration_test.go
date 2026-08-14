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

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://localhost/api/console/port-groups", strings.NewReader(`{"items":[{"portPrefix":1000,"environmentName":"miniport","runtimeMode":"dind","runtimeName":"miniport-dind-01","serviceIp":"172.22.11.12","slots":[{"port":10000,"name":"redis","kind":"cache","protocol":"redis"}]}]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost")
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())

	var body struct {
		Items []struct {
			EnvironmentName string `json:"environmentName"`
			Slots           []struct {
				Name string `json:"name"`
			} `json:"slots"`
		} `json:"items"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Len(t, body.Items, 1)
	require.Equal(t, "miniport", body.Items[0].EnvironmentName)
	require.Len(t, body.Items[0].Slots, 1)
	require.Equal(t, "redis", body.Items[0].Slots[0].Name)
}

func TestPortsvcReadRejectsMissingSession(t *testing.T) {
	handler := setupPortsvcTestApp(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/console/port-groups", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code, rec.Body.String())
}

func TestGithubStatusUsesLocalSession(t *testing.T) {
	handler := setupPortsvcTestApp(t)
	cookie := createTestSessionCookie(t, handler)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/console/github/status", nil)
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

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/console/github/status", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code, rec.Body.String())
}

func TestGithubSetupDoesNotRequireSessionCookie(t *testing.T) {
	handler := setupPortsvcTestApp(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/integrations/github/setup?installation_id=42&state=one-time-state", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusServiceUnavailable, rec.Code, rec.Body.String())
}

func TestPortsvcBatchWriteRollsBackOnValidationFailure(t *testing.T) {
	handler := setupPortsvcTestApp(t)
	cookie := createTestSessionCookie(t, handler)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://localhost/api/console/hosts", strings.NewReader(`{"items":[{"name":"valid-host"},{"name":" "}]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost")
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/console/hosts", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var body struct {
		Items []json.RawMessage `json:"items"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Empty(t, body.Items)
}

func TestPortsvcBatchConflictRollsBack(t *testing.T) {
	handler := setupPortsvcTestApp(t)
	cookie := createTestSessionCookie(t, handler)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://localhost/api/console/hosts", strings.NewReader(`{"items":[{"name":"duplicate-host"},{"name":"duplicate-host"}]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost")
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusConflict, rec.Code, rec.Body.String())
	var problem struct {
		Type string `json:"type"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &problem))
	require.Equal(t, "urn:problem:portsvc-resource-conflict", problem.Type)

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/console/hosts", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var body struct {
		Items []json.RawMessage `json:"items"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Empty(t, body.Items)
}

func TestPortsvcPortGroupBatchRollsBackAggregateWrites(t *testing.T) {
	handler := setupPortsvcTestApp(t)
	cookie := createTestSessionCookie(t, handler)

	payload := `{"items":[` +
		`{"portPrefix":1000,"environmentName":"first","runtimeMode":"dind","slots":[{"port":10000,"name":"api"}]},` +
		`{"portPrefix":1000,"environmentName":"conflict","runtimeMode":"dind"}` +
		`]}`
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://localhost/api/console/port-groups", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost")
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusConflict, rec.Code, rec.Body.String())

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/console/port-groups", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var body struct {
		Items []json.RawMessage `json:"items"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Empty(t, body.Items)
}

func TestLegacyAPIPathIsNotRegistered(t *testing.T) {
	handler := setupPortsvcTestApp(t)
	cookie := createTestSessionCookie(t, handler)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/api/port-groups", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNotFound, rec.Code, rec.Body.String())
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
