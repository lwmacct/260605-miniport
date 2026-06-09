package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestFrontendRootServesIndex(t *testing.T) {
	router := testFrontendRouter(t)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Body.String() != "index" {
		t.Fatalf("body = %q, want %q", recorder.Body.String(), "index")
	}
}

func TestFrontendServesStaticAsset(t *testing.T) {
	router := testFrontendRouter(t)
	request := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Body.String() != "asset" {
		t.Fatalf("body = %q, want %q", recorder.Body.String(), "asset")
	}
}

func TestFrontendRouteFallsBackToIndex(t *testing.T) {
	router := testFrontendRouter(t)
	request := httptest.NewRequest(http.MethodGet, "/hosts/10.0.0.1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Body.String() != "index" {
		t.Fatalf("body = %q, want %q", recorder.Body.String(), "index")
	}
}

func TestFrontendBypassesAPIRoutes(t *testing.T) {
	router := testFrontendRouter(t)
	request := httptest.NewRequest(http.MethodGet, "/api/missing", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestFrontendDisabledWithoutIndex(t *testing.T) {
	router := chi.NewRouter()
	if registerFrontendRoutes(router, t.TempDir()) {
		t.Fatal("registerFrontendRoutes returned true without index.html")
	}
}

func testFrontendRouter(t *testing.T) chi.Router {
	t.Helper()

	router := chi.NewRouter()
	if !registerFrontendRoutes(router, writeTestFrontend(t)) {
		t.Fatal("registerFrontendRoutes returned false")
	}
	return router
}

func writeTestFrontend(t *testing.T) string {
	t.Helper()

	uiDir := t.TempDir()
	assetsDir := filepath.Join(uiDir, "assets")
	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		t.Fatalf("mkdir assets: %v", err)
	}
	if err := os.WriteFile(filepath.Join(uiDir, "index.html"), []byte("index"), 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(assetsDir, "app.js"), []byte("asset"), 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}
	return uiDir
}
