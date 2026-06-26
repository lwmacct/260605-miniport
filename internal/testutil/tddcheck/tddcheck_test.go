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
	return tddcheck.DefaultConfig()
}
