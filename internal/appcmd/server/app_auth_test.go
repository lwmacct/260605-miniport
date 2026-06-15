package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260605-miniport/internal/config"
	"github.com/lwmacct/260605-miniport/internal/domain/identityuser"
)

func TestAuthPasswordLoginSessionAndAdminUsers(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Server.Database.SQLite = ":memory:"
	cfg.Server.Auth.Admins = []string{"admin"}
	app := NewApp(&cfg)
	require.NoError(t, app.bootstrap(context.Background()))
	defer app.closeDatabase()

	admin, err := app.users.Create(context.Background(), identityuser.CreateInput{Username: "admin"})
	require.NoError(t, err)
	require.NoError(t, app.passwords.Set(context.Background(), admin.Username, admin.ID, "strong-ops-password-123"))

	srv, err := app.newHTTPServer()
	require.NoError(t, err)
	handler := srv.Handler

	body := `{"username":"admin","password":"strong-ops-password-123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/password/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	cookie := rec.Result().Cookies()[0]

	req = httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var session struct {
		Authenticated bool `json:"authenticated"`
		User          struct {
			Username string `json:"username"`
			Admin    bool   `json:"admin"`
		} `json:"user"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &session))
	require.True(t, session.Authenticated)
	require.Equal(t, "admin", session.User.Username)
	require.True(t, session.User.Admin)

	req = httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
}

func TestInventoryWriteRequiresAdmin(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Server.Database.SQLite = ":memory:"
	cfg.Server.Auth.Admins = []string{"root-admin"}
	app := NewApp(&cfg)
	require.NoError(t, app.bootstrap(context.Background()))
	defer app.closeDatabase()

	user, err := app.users.Create(context.Background(), identityuser.CreateInput{Username: "member"})
	require.NoError(t, err)
	require.NoError(t, app.passwords.Set(context.Background(), user.Username, user.ID, "team-secret-123"))

	srv, err := app.newHTTPServer()
	require.NoError(t, err)
	handler := srv.Handler

	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/password/login", strings.NewReader(`{"username":"member","password":"team-secret-123"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.RemoteAddr = "127.0.0.1:12345"
	loginRec := httptest.NewRecorder()
	handler.ServeHTTP(loginRec, loginReq)
	require.Equal(t, http.StatusOK, loginRec.Code, loginRec.Body.String())
	cookie := loginRec.Result().Cookies()[0]

	req := httptest.NewRequest(http.MethodPost, "/api/hosts", strings.NewReader(`{"ip":"172.22.11.12","name":"node-12","network":"","environment":"","notes":""}`))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusForbidden, rec.Code, rec.Body.String())
}
