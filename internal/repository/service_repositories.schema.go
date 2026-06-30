package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type ServiceRepositoriesModel struct {
	bun.BaseModel `bun:"table:service_repositories,alias:source_link"`

	ID           string    `bun:"id,pk,type:uuid" json:"id"`
	ServiceID    string    `bun:"service_id,type:uuid,notnull" json:"serviceId"`
	RepositoryID string    `bun:"repository_id,type:uuid,notnull" json:"repositoryId"`
	Role         string    `bun:"role" json:"role"`
	CreatedAt    time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt    time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func ServiceReposSchema() []any {
	return []any{(*ServiceRepositoriesModel)(nil)}
}

func (*ServiceRepositoriesModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(service_id) REFERENCES services (id) ON DELETE CASCADE")
	query.ForeignKey("(repository_id) REFERENCES repositories (id) ON DELETE CASCADE")
	return nil
}

func ServiceReposIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_service_repositories_pair ON service_repositories(service_id, repository_id)`,
		`CREATE INDEX IF NOT EXISTS idx_service_repositories_repository ON service_repositories(repository_id)`,
	}
}
