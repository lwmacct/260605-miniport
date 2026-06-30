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

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/services", strings.NewReader(`{"name":"member-service","projectName":"miniport","dindIp":"172.22.11.12"}`))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var body struct {
		UserID   int64  `json:"userId"`
		Username string `json:"username"`
		Name     string `json:"name"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.NotZero(t, body.UserID)
	require.Equal(t, "member", body.Username)
	require.Equal(t, "member-service", body.Name)
}

func TestPortsvcReadRejectsMissingSession(t *testing.T) {
	_, handler := setupPortsvcTestApp(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/services", nil)
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

func createTestSessionCookie(t *testing.T, app *App, username string) *http.Cookie {
	t.Helper()

	user, err := app.deps.auth.CreateExternalUser(context.Background(), auth.ExternalUserInput{
		Username:    username,
		DisplayName: username,
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
		Path:       "/api/services",
		RemoteAddr: "127.0.0.1:12345",
	}
}
