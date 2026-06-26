package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260605-miniport/internal/config"
	"github.com/lwmacct/260605-miniport/internal/service"
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
	return service.AuthChallengeProviderImage
}

func (p *testChallengeProvider) PublicConfig() service.AuthChallengePublicConfig {
	return service.AuthChallengePublicConfig{Provider: service.AuthChallengeProviderImage}
}

func (p *testChallengeProvider) Create(context.Context, service.AuthChallengeInput) (*service.AuthChallenge, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.nextID++
	id := "cap_test_" + strconv.Itoa(p.nextID)
	p.answers[id] = "PASS"
	return &service.AuthChallenge{
		Provider:    service.AuthChallengeProviderImage,
		ChallengeID: id,
		Image:       "data:image/png;base64,test",
		ExpiresAt:   time.Now().Add(time.Minute),
	}, nil
}

func (p *testChallengeProvider) Verify(_ context.Context, response service.AuthChallengeAnswer, _ service.AuthChallengeInput) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	answer, ok := p.answers[response.ChallengeID]
	delete(p.answers, response.ChallengeID)
	if !ok || response.Provider != service.AuthChallengeProviderImage || response.Answer != answer {
		return service.ErrAuthChallengeInvalid
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
	app.container.challenges = service.NewAuthChallengeService(newTestChallengeProvider())
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
	t.Helper()
	return `{"username":"` + username + `","password":"` + password + `",` + newChallengeBody(t, handler) + `}`
}

func TestAuthPasswordLoginSessionAndAdminUsers(t *testing.T) {
	app, handler := setupAuthTestApp(t, []string{"admin"})

	admin, err := app.container.users.Create(context.Background(), service.CreateIdentityUserInput{Username: "admin"})
	require.NoError(t, err)
	require.NoError(t, app.container.passwords.Set(context.Background(), admin.Username, admin.ID, "strong-ops-password-123"))

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

func TestAuthStateReturnsPublicConfigForGuest(t *testing.T) {
	_, handler := setupAuthTestApp(t, []string{"admin"})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/auth/state", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var state struct {
		Config struct {
			Challenge struct {
				Provider string `json:"provider"`
			} `json:"challenge"`
		} `json:"config"`
		Session struct {
			Authenticated bool `json:"authenticated"`
		} `json:"session"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &state))
	require.Equal(t, "image", state.Config.Challenge.Provider)
	require.False(t, state.Session.Authenticated)
}

func TestAuthStateReturnsCurrentUser(t *testing.T) {
	app, handler := setupAuthTestApp(t, []string{"admin"})
	user, err := app.container.users.Create(context.Background(), service.CreateIdentityUserInput{Username: "admin"})
	require.NoError(t, err)
	require.NoError(t, app.container.passwords.Set(context.Background(), user.Username, user.ID, "strong-ops-password-123"))

	loginReq := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/auth/password/login", strings.NewReader(loginBody(t, handler, "admin", "strong-ops-password-123")))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.RemoteAddr = "127.0.0.1:12345"
	loginRec := httptest.NewRecorder()
	handler.ServeHTTP(loginRec, loginReq)
	require.Equal(t, http.StatusOK, loginRec.Code, loginRec.Body.String())

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/auth/state", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(loginRec.Result().Cookies()[0])
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var state struct {
		Session struct {
			Authenticated bool `json:"authenticated"`
			User          struct {
				Username string `json:"username"`
				Admin    bool   `json:"admin"`
			} `json:"user"`
		} `json:"session"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &state))
	require.True(t, state.Session.Authenticated)
	require.Equal(t, "admin", state.Session.User.Username)
	require.True(t, state.Session.User.Admin)
}

func TestInventoryWriteRequiresAdmin(t *testing.T) {
	app, handler := setupAuthTestApp(t, []string{"root-admin"})

	user, err := app.container.users.Create(context.Background(), service.CreateIdentityUserInput{Username: "member"})
	require.NoError(t, err)
	require.NoError(t, app.container.passwords.Set(context.Background(), user.Username, user.ID, "team-secret-123"))

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

func TestProtectedReadRejectsDuplicateSessionCookie(t *testing.T) {
	app, handler := setupAuthTestApp(t, []string{"admin"})

	user, err := app.container.users.Create(context.Background(), service.CreateIdentityUserInput{Username: "admin"})
	require.NoError(t, err)
	require.NoError(t, app.container.passwords.Set(context.Background(), user.Username, user.ID, "strong-ops-password-123"))

	loginReq := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/auth/password/login", strings.NewReader(loginBody(t, handler, "admin", "strong-ops-password-123")))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.RemoteAddr = "127.0.0.1:12345"
	loginRec := httptest.NewRecorder()
	handler.ServeHTTP(loginRec, loginReq)
	require.Equal(t, http.StatusOK, loginRec.Code, loginRec.Body.String())
	validCookie := loginRec.Result().Cookies()[0]

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/hosts", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("Cookie", validCookie.Name+"=sess_stale; "+validCookie.Name+"="+validCookie.Value)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code, rec.Body.String())

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/hosts", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.AddCookie(validCookie)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
}

func TestAuthPasswordLoginRequiresChallenge(t *testing.T) {
	app, handler := setupAuthTestApp(t, []string{"admin"})
	user, err := app.container.users.Create(context.Background(), service.CreateIdentityUserInput{Username: "admin"})
	require.NoError(t, err)
	require.NoError(t, app.container.passwords.Set(context.Background(), user.Username, user.ID, "strong-ops-password-123"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/auth/password/login", strings.NewReader(`{"username":"admin","password":"strong-ops-password-123","challenge":{"provider":"image","challengeId":"missing","answer":"PASS"}}`))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code, rec.Body.String())
}

func TestChallengeProviderConsumesChallengeOnce(t *testing.T) {
	provider := newTestChallengeProvider()
	challengeService := service.NewAuthChallengeService(provider)
	challenge, err := challengeService.Create(context.Background(), service.AuthChallengeInput{})
	require.NoError(t, err)

	response := service.AuthChallengeAnswer{Provider: challenge.Provider, ChallengeID: challenge.ChallengeID, Answer: "PASS"}
	require.NoError(t, challengeService.Verify(context.Background(), response, service.AuthChallengeInput{}))
	err = challengeService.Verify(context.Background(), response, service.AuthChallengeInput{})
	require.ErrorIs(t, err, service.ErrAuthChallengeInvalid)
}
