package frontend

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router, configuredDir string) bool {
	uiDir, ok := detectDir(configuredDir)
	if !ok {
		return false
	}

	indexPath := filepath.Join(uiDir, "index.html")
	serveIndex := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexPath)
	}

	router.MethodFunc(http.MethodGet, "/", serveIndex)
	router.MethodFunc(http.MethodHead, "/", serveIndex)
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		if shouldBypass(requestPath, r.Method) {
			http.NotFound(w, r)
			return
		}

		if assetPath, ok := assetPath(uiDir, requestPath); ok && fileExists(assetPath) {
			http.ServeFile(w, r, assetPath)
			return
		}

		http.ServeFile(w, r, indexPath)
	})

	return true
}

func detectDir(configuredDir string) (string, bool) {
	candidate := strings.TrimSpace(configuredDir)
	if candidate == "" || !fileExists(filepath.Join(candidate, "index.html")) {
		return "", false
	}

	absPath, err := filepath.Abs(candidate)
	if err != nil {
		return candidate, true
	}
	return absPath, true
}

func shouldBypass(requestPath, method string) bool {
	if method != http.MethodGet && method != http.MethodHead {
		return true
	}
	return requestPath == "/api" || strings.HasPrefix(requestPath, "/api/")
}

func assetPath(uiDir, requestPath string) (string, bool) {
	cleanPath := filepath.Clean("/" + strings.TrimPrefix(requestPath, "/"))
	if cleanPath == "/" {
		return "", false
	}

	candidate := filepath.Join(uiDir, strings.TrimPrefix(cleanPath, "/"))
	relPath, err := filepath.Rel(uiDir, candidate)
	if err != nil || relPath == ".." || strings.HasPrefix(relPath, "../") {
		return "", false
	}
	return candidate, true
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
