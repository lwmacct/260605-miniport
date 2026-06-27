package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type InventoryProjectModel struct {
	bun.BaseModel `bun:"table:allocation_projects,alias:project"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	PortGroupID int64     `bun:"allocation_id,notnull" json:"allocationId"`
	Name        string    `bun:"name,notnull" json:"name"`
	Description string    `bun:"description" json:"description"`
	Notes       string    `bun:"notes" json:"notes"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func InventoryProjectSchema() []any {
	return []any{(*InventoryProjectModel)(nil)}
}

func (*InventoryProjectModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(allocation_id) REFERENCES port_allocations (id) ON DELETE CASCADE")
	return nil
}

func InventoryProjectIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_allocation_projects_allocation ON allocation_projects(allocation_id)`,
	}
}
