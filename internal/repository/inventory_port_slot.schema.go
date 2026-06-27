package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type InventoryPortSlotModel struct {
	bun.BaseModel `bun:"table:allocation_ports,alias:port_slot"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	PortGroupID int64     `bun:"allocation_id,notnull" json:"allocationId"`
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
	query.ForeignKey("(allocation_id) REFERENCES port_allocations (id) ON DELETE CASCADE")
	return nil
}

func InventoryPortSlotIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_allocation_ports_allocation_port ON allocation_ports(allocation_id, port)`,
	}
}
