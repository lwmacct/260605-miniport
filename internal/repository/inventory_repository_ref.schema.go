package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type InventoryRepositoryRefModel struct {
	bun.BaseModel `bun:"table:inventory_repository_refs,alias:repository_ref"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	PortGroupID int64     `bun:"port_group_id,notnull" json:"portGroupId"`
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
	query.ForeignKey("(port_group_id) REFERENCES inventory_port_groups (id) ON DELETE CASCADE")
	return nil
}

func InventoryRepositoryRefIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_inventory_repository_refs_group ON inventory_repository_refs(port_group_id)`,
	}
}
