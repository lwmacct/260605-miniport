package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type InventoryComponentModel struct {
	bun.BaseModel `bun:"table:allocation_dependencies,alias:dependency"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	PortGroupID int64     `bun:"allocation_id,notnull" json:"allocationId"`
	Name        string    `bun:"name,notnull" json:"name"`
	Type        string    `bun:"type,notnull" json:"type"`
	URL         string    `bun:"url" json:"url"`
	Version     string    `bun:"version" json:"version"`
	Notes       string    `bun:"notes" json:"notes"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func InventoryComponentSchema() []any {
	return []any{(*InventoryComponentModel)(nil)}
}

func (*InventoryComponentModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(allocation_id) REFERENCES port_allocations (id) ON DELETE CASCADE")
	return nil
}

func InventoryComponentIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_allocation_dependencies_allocation ON allocation_dependencies(allocation_id)`,
	}
}
