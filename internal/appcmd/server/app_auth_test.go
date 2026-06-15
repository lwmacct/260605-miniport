package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260605-miniport/internal/config"
	"github.com/lwmacct/260605-miniport/internal/domain/authchallenge"
	"github.com/lwmacct/260605-miniport/internal/domain/identityuser"
)

type testChallengeProvider struct {
	mu      sync.Mutex
	nextID  int
	answers map[string]string
}

func newTestChallengeProvider() *testChallengeProvider {
	return &testChallengeProvider{answers: map[string]string{}}
}

func (p *testChallengeProvider) Name() string {
	return authchallenge.ProviderImage
}

func (p *testChallengeProvider) PublicConfig() authchallenge.PublicConfigDTO {
	return authchallenge.PublicConfigDTO{Provider: authchallenge.ProviderImage}
}

func (p *testChallengeProvider) Create(context.Context, authchallenge.RequestDTO) (*authchallenge.ChallengeDTO, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.nextID++
	id := "cap_test_" + string(rune('0'+p.nextID))
	p.answers[id] = "PASS"
	return &authchallenge.ChallengeDTO{
		Provider:    authchallenge.ProviderImage,
		ChallengeID: id,
		Image:       "data:image/png;base64,test",
		ExpiresAt:   time.Now().Add(time.Minute),
	}, nil
}

func (p *testChallengeProvider) Verify(_ context.Context, response authchallenge.ResponseDTO, _ authchallenge.RequestDTO) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	answer, ok := p.answers[response.ChallengeID]
	delete(p.answers, response.ChallengeID)
	if !ok || response.Provider != authchallenge.ProviderImage || response.Answer != answer {
		return authchallenge.ErrInvalidChallenge
	}
	return nil
}

func setupAuthTestApp(t *testing.T, admins []string) (*App, http.Handler) {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.Server.Database.SQLite = ":memory:"
	cfg.Server.Auth.Admins = admins
	app := NewApp(&cfg)
	require.NoError(t, app.bootstrap(context.Background()))
	app.challenges = authchallenge.NewService(newTestChallengeProvider())
	t.Cleanup(app.closeDatabase)
	return app, app.newHTTPServer().Handler
}

func newChallengeBody(t *testing.T, handler http.Handler) string {
	t.Helper()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/auth/challenges", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var challenge struct {
		Provider    string `json:"provider"`
		ChallengeID string `json:"challengeId"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &challenge))
	return `"challenge":{"provider":"` + challenge.Provider + `","challengeId":"` + challenge.ChallengeID + `","answer":"PASS"}`
}

func loginBody(t *testing.T, handler http.Handler, username string, password string) string {
	return `{"username":"` + username + `","password":"` + password + `",` + newChallengeBody(t, handler) + `}`
}

func TestAuthPasswordLoginSessionAndAdminUsers(t *testing.T) {
	app, handler := setupAuthTestApp(t, []string{"admin"})

	admin, err := app.users.Create(context.Background(), identityuser.CreateInput{Username: "admin"})
	require.NoError(t, err)
	require.NoError(t, app.passwords.Set(context.Background(), admin.Username, admin.ID, "strong-ops-password-123"))

	body := loginBody(t, handler, "admin", "strong-ops-password-123")
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/auth/password/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	cookie := rec.Result().Cookies()[0]

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/auth/me", nil)
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

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/admin/users", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
}

func TestInventoryWriteRequiresAdmin(t *testing.T) {
	app, handler := setupAuthTestApp(t, []string{"root-admin"})

	user, err := app.users.Create(context.Background(), identityuser.CreateInput{Username: "member"})
	require.NoError(t, err)
	require.NoError(t, app.passwords.Set(context.Background(), user.Username, user.ID, "team-secret-123"))

	loginReq := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/auth/password/login", strings.NewReader(loginBody(t, handler, "member", "team-secret-123")))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.RemoteAddr = "127.0.0.1:12345"
	loginRec := httptest.NewRecorder()
	handler.ServeHTTP(loginRec, loginReq)
	require.Equal(t, http.StatusOK, loginRec.Code, loginRec.Body.String())
	cookie := loginRec.Result().Cookies()[0]

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/hosts", strings.NewReader(`{"ip":"172.22.11.12","name":"node-12","network":"","environment":"","notes":""}`))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusForbidden, rec.Code, rec.Body.String())
}

func TestAuthPasswordLoginRequiresChallenge(t *testing.T) {
	app, handler := setupAuthTestApp(t, []string{"admin"})
	user, err := app.users.Create(context.Background(), identityuser.CreateInput{Username: "admin"})
	require.NoError(t, err)
	require.NoError(t, app.passwords.Set(context.Background(), user.Username, user.ID, "strong-ops-password-123"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/auth/password/login", strings.NewReader(`{"username":"admin","password":"strong-ops-password-123","challenge":{"provider":"image","challengeId":"missing","answer":"PASS"}}`))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code, rec.Body.String())
}

func TestChallengeProviderConsumesChallengeOnce(t *testing.T) {
	provider := newTestChallengeProvider()
	service := authchallenge.NewService(provider)
	challenge, err := service.Create(context.Background(), authchallenge.RequestDTO{})
	require.NoError(t, err)

	response := authchallenge.ResponseDTO{Provider: challenge.Provider, ChallengeID: challenge.ChallengeID, Answer: "PASS"}
	require.NoError(t, service.Verify(context.Background(), response, authchallenge.RequestDTO{}))
	err = service.Verify(context.Background(), response, authchallenge.RequestDTO{})
	require.True(t, errors.Is(err, authchallenge.ErrInvalidChallenge))
}
