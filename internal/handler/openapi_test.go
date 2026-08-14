package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/stretchr/testify/require"
)

func TestOpenAPIContract(t *testing.T) {
	cfg := huma.DefaultConfig("Miniport API", "1.0.0")
	cfg.Servers = []*huma.Server{{URL: "/api"}}
	api := humago.New(http.NewServeMux(), cfg)
	RegisterCore(api)
	RegisterGithub(api, nil)
	RegisterPortsvc(api, Services{})

	document, err := json.Marshal(api.OpenAPI())
	require.NoError(t, err)
	var spec struct {
		Paths map[string]map[string]struct {
			OperationID string `json:"operationId"`
			Responses   map[string]struct {
				Content map[string]json.RawMessage `json:"content"`
			} `json:"responses"`
		} `json:"paths"`
		Components struct {
			Schemas map[string]struct {
				Properties map[string]struct {
					MinItems *int `json:"minItems"`
				} `json:"properties"`
			} `json:"schemas"`
		} `json:"components"`
	}
	require.NoError(t, json.Unmarshal(document, &spec))

	expectedPaths := []string{
		"/console/dependency-assets",
		"/console/github/connections",
		"/console/github/installations",
		"/console/github/installations/sync",
		"/console/github/repositories",
		"/console/github/status",
		"/console/hosts",
		"/console/port-groups",
		"/console/port-groups/export.csv",
		"/console/service-groups",
		"/health",
		"/integrations/github/setup",
		"/integrations/github/webhooks",
		"/meta",
	}
	require.ElementsMatch(t, expectedPaths, mapKeys(spec.Paths))
	require.NotContains(t, spec.Paths, "/hosts")
	require.NotContains(t, spec.Paths, "/port-slots/{id}")

	operationIDs := map[string]struct{}{}
	for _, path := range spec.Paths {
		for _, operation := range path {
			if operation.OperationID == "" {
				continue
			}
			_, exists := operationIDs[operation.OperationID]
			require.False(t, exists, "duplicate operation ID %s", operation.OperationID)
			operationIDs[operation.OperationID] = struct{}{}
		}
	}

	_, hasRedirect := spec.Paths["/integrations/github/setup"]["get"].Responses["303"]
	require.True(t, hasRedirect)
	csvResponse := spec.Paths["/console/port-groups/export.csv"]["get"].Responses["200"]
	require.Contains(t, csvResponse.Content, "text/csv")
	items := spec.Components.Schemas["HostCreateBatchDTO"].Properties["items"]
	require.NotNil(t, items.MinItems)
	require.Equal(t, 1, *items.MinItems)
}

func mapKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}
