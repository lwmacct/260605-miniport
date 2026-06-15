package tddcheck_test

import (
	"testing"

	"github.com/lwmacct/260610-go-pkg-tddcheck/pkg/tddcheck/rules/files/context"
	"github.com/lwmacct/260610-go-pkg-tddcheck/pkg/tddcheck/rules/files/cqrs"
	"github.com/lwmacct/260610-go-pkg-tddcheck/pkg/tddcheck/rules/files/entity"
	"github.com/lwmacct/260610-go-pkg-tddcheck/pkg/tddcheck/rules/files/errors"
	"github.com/lwmacct/260610-go-pkg-tddcheck/pkg/tddcheck/rules/files/handler"
	"github.com/lwmacct/260610-go-pkg-tddcheck/pkg/tddcheck/rules/files/mapper"
	"github.com/lwmacct/260610-go-pkg-tddcheck/pkg/tddcheck/rules/other/databasetest"
	"github.com/lwmacct/260610-go-pkg-tddcheck/pkg/tddcheck/rules/other/errorprefix"
	"github.com/lwmacct/260610-go-pkg-tddcheck/pkg/tddcheck/rules/other/layerdeps"
	"github.com/lwmacct/260610-go-pkg-tddcheck/pkg/tddcheck/rules/other/packagename"
	"github.com/lwmacct/260610-go-pkg-tddcheck/pkg/tddcheck/rules/other/publicapi"
)

type assertRule interface {
	Assert(t *testing.T)
}

func TestRules(t *testing.T) {
	tests := []struct {
		name string
		rule assertRule
	}{
		{"dependency-layerdeps", layerdeps.New("internal")},
		{"file-cqrs", cqrs.New("internal")},
		{"file-entity", entity.New("internal")},
		{"name-error-prefix", errorprefix.New("internal")},
		{"file-errors", errors.New("internal")},
		{"file-mapper", mapper.New("internal")},
		{"file-handler", handler.New("internal")},
		{"file-context", context.New("internal")},
		{"name-public-api", publicapi.New("internal")},
		{"name-package", packagename.New("internal")},
		{"test-database", databasetest.New(".")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.rule.Assert(t)
		})
	}
}
