package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type ServiceDependenciesModel struct {
	bun.BaseModel `bun:"table:service_dependencies,alias:dependency_link"`

	ID           int64     `bun:"id,pk,autoincrement" json:"id"`
	ServiceID    int64     `bun:"service_id,notnull" json:"serviceId"`
	DependencyID int64     `bun:"dependency_id,notnull" json:"dependencyId"`
	Role         string    `bun:"role" json:"role"`
	Notes        string    `bun:"notes" json:"notes"`
	CreatedAt    time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt    time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func ServiceDepsSchema() []any {
	return []any{(*ServiceDependenciesModel)(nil)}
}

func (*ServiceDependenciesModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(service_id) REFERENCES services (id) ON DELETE CASCADE")
	query.ForeignKey("(dependency_id) REFERENCES dependencies (id) ON DELETE CASCADE")
	return nil
}

func ServiceDepsIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_service_dependencies_pair ON service_dependencies(service_id, dependency_id)`,
		`CREATE INDEX IF NOT EXISTS idx_service_dependencies_dependency ON service_dependencies(dependency_id)`,
	}
}
