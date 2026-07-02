package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lwmacct/260630-go-hsr-auth/pkg/auth"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260605-miniport/internal/config"
)

func TestPortsvcWriteAllowsAuthenticatedUser(t *testing.T) {
	app, handler := setupPortsvcTestApp(t)
	cookie := createTestSessionCookie(t, app, "member")

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/port-groups", strings.NewReader(`{"portPrefix":1000,"environmentName":"miniport","runtimeMode":"dind","runtimeName":"miniport-dind-01","serviceIp":"172.22.11.12","slots":[{"port":10000,"name":"redis","kind":"cache","protocol":"redis"}]}`))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var body struct {
		OwnerSubject    string `json:"ownerSubject"`
		OwnerName       string `json:"ownerName"`
		EnvironmentName string `json:"environmentName"`
		Slots           []struct {
			Name string `json:"name"`
		} `json:"slots"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.NotEmpty(t, body.OwnerSubject)
	require.Equal(t, "member", body.OwnerName)
	require.Equal(t, "miniport", body.EnvironmentName)
	require.Len(t, body.Slots, 1)
	require.Equal(t, "redis", body.Slots[0].Name)
}

func TestPortsvcReadRejectsMissingSession(t *testing.T) {
	_, handler := setupPortsvcTestApp(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/port-groups", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code, rec.Body.String())
}

func setupPortsvcTestApp(t *testing.T) (*App, http.Handler) {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.Server.Database.SQLite = "file:" + strings.NewReplacer("/", "-", " ", "-").Replace(t.Name()) + "?mode=memory&cache=shared"
	cfg.Server.HTTP.WebRoot = ""
	app := NewApp(&cfg)
	deps, err := newDependencies(context.Background(), &cfg)
	require.NoError(t, err)
	app.deps = deps
	t.Cleanup(app.deps.Close)
	return app, app.newHTTPServer().Handler
}

func createTestSessionCookie(t *testing.T, app *App, ownerName string) *http.Cookie {
	t.Helper()

	user, err := app.deps.auth.CreateExternalUser(context.Background(), auth.ExternalUserInput{
		Username:    ownerName,
		DisplayName: ownerName,
	})
	require.NoError(t, err)
	session, err := app.deps.auth.CreateSession(auth.ContextWithRequest(context.Background(), testSessionRequest()), user.ID, testSessionRequest())
	require.NoError(t, err)

	cookie, err := http.ParseSetCookie(session.SetCookie)
	require.NoError(t, err)
	return cookie
}

func testSessionRequest() auth.SessionRequest {
	return auth.SessionRequest{
		IP:         "127.0.0.1",
		Scheme:     "http",
		Host:       "example.test",
		UserAgent:  "test",
		Method:     "GET",
		Path:       "/api/port-groups",
		RemoteAddr: "127.0.0.1:12345",
	}
}
