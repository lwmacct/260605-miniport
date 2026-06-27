package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type InventoryRepositoryRefModel struct {
	bun.BaseModel `bun:"table:allocation_repositories,alias:repository_ref"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	PortGroupID int64     `bun:"allocation_id,notnull" json:"allocationId"`
	ProjectID   int64     `bun:"project_id" json:"projectId"`
	Name        string    `bun:"name,notnull" json:"name"`
	URL         string    `bun:"url,notnull" json:"url"`
	Kind        string    `bun:"kind,notnull" json:"kind"`
	Notes       string    `bun:"notes" json:"notes"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func InventoryRepositoryRefSchema() []any {
	return []any{(*InventoryRepositoryRefModel)(nil)}
}

func (*InventoryRepositoryRefModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(allocation_id) REFERENCES port_allocations (id) ON DELETE CASCADE")
	return nil
}

func InventoryRepositoryRefIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_allocation_repositories_allocation ON allocation_repositories(allocation_id)`,
		`CREATE INDEX IF NOT EXISTS idx_allocation_repositories_project ON allocation_repositories(project_id)`,
	}
}
