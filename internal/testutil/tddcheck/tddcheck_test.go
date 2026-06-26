package tddcheck_test

import (
	"testing"

	"github.com/lwmacct/260622-go-pkg-tddcheck/pkg/tddcheck"
)

func TestRules(t *testing.T) {
	tddcheck.ProjectRules{
		Root:   "internal",
		Config: tddcheck.DefaultConfig(),
	}.Assert(t)
}
