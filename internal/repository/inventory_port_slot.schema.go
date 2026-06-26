package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type InventoryPortSlotModel struct {
	bun.BaseModel `bun:"table:inventory_port_slots,alias:port_slot"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	PortGroupID int64     `bun:"port_group_id,notnull" json:"portGroupId"`
	Port        int       `bun:"port,notnull" json:"port"`
	Name        string    `bun:"name" json:"name"`
	Protocol    string    `bun:"protocol,notnull" json:"protocol"`
	Purpose     string    `bun:"purpose" json:"purpose"`
	Status      string    `bun:"status,notnull" json:"status"`
	Notes       string    `bun:"notes" json:"notes"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func InventoryPortSlotSchema() []any {
	return []any{(*InventoryPortSlotModel)(nil)}
}

func (*InventoryPortSlotModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(port_group_id) REFERENCES inventory_port_groups (id) ON DELETE CASCADE")
	return nil
}

func InventoryPortSlotIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_inventory_port_slots_group_port ON inventory_port_slots(port_group_id, port)`,
	}
}
