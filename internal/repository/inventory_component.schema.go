package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type InventoryComponentModel struct {
	bun.BaseModel `bun:"table:inventory_components,alias:component"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	PortGroupID int64     `bun:"port_group_id,notnull" json:"portGroupId"`
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
	query.ForeignKey("(port_group_id) REFERENCES inventory_port_groups (id) ON DELETE CASCADE")
	return nil
}

func InventoryComponentIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_inventory_components_group ON inventory_components(port_group_id)`,
	}
}
