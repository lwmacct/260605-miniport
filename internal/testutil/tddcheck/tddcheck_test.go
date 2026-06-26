package tddcheck_test

import (
	"testing"

	"github.com/lwmacct/260622-go-pkg-tddcheck/pkg/tddcheck"
)

func TestRules(t *testing.T) {
	tddcheck.ProjectRules{
		Root:   "internal",
		Config: projectRulesConfig(),
	}.Assert(t)
}

func projectRulesConfig() tddcheck.Config {
	config := tddcheck.DefaultConfig()
	config.LayerDirs = []string{"handler", "service", "repository"}
	config.DependencyLayerDirs = []string{
		"adapter",
		"appcmd",
		"config",
		"handler",
		"infra",
		"repository",
		"service",
		"testutil",
	}
	config.LayerFileKinds["service"] = []string{
		"commands",
		"constants",
		"dto",
		"entity",
		"errors",
		"mapper",
		"provider",
		"service",
		"support",
		"utils",
		"validation",
	}
	for _, sourceLayer := range []string{"handler", "service", "repository", "infra", "appcmd"} {
		config.LayerRules = append(config.LayerRules, tddcheck.LayerDependencyRule{
			SourceLayer: sourceLayer,
			TargetLayer: "adapter",
			Message:     sourceLayer + " must not import adapter",
		})
	}
	return config
}
