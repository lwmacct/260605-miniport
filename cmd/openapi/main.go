package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"github.com/lwmacct/260605-miniport/internal/handler"
)

const openAPIOutput = "openapi/openapi.json"

func main() {
	check := flag.Bool("check", false, "check that the OpenAPI document is current")
	flag.Parse()
	if err := generateOpenAPI(*check); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "generate OpenAPI: %v\n", err)
		os.Exit(1)
	}
}

func generateOpenAPI(check bool) error {
	if err := os.MkdirAll(filepath.Dir(openAPIOutput), 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	mux := http.NewServeMux()
	cfg := huma.DefaultConfig("Miniport API", "1.0.0")
	cfg.Servers = []*huma.Server{{URL: "/api"}}
	api := humago.New(mux, cfg)
	handler.RegisterCore(api)
	handler.RegisterGithub(api, nil)
	handler.RegisterPortsvc(api, handler.Services{})

	document, err := json.MarshalIndent(api.OpenAPI(), "", "  ")
	if err != nil {
		return fmt.Errorf("encode document: %w", err)
	}
	document = append(document, '\n')
	if check {
		current, readErr := os.ReadFile(openAPIOutput)
		if readErr != nil {
			return fmt.Errorf("read current document: %w", readErr)
		}
		if !bytes.Equal(current, document) {
			return fmt.Errorf("%s is stale; run pnpm generate:api", openAPIOutput)
		}
		return nil
	}
	if err := os.WriteFile(openAPIOutput, document, 0o600); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	return nil
}
