package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type InventoryHostModel struct {
	bun.BaseModel `bun:"table:inventory_hosts,alias:host"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	IP          string    `bun:"ip,notnull,unique" json:"ip"`
	Name        string    `bun:"name" json:"name"`
	Network     string    `bun:"network" json:"network"`
	Environment string    `bun:"environment" json:"environment"`
	Notes       string    `bun:"notes" json:"notes"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func InventoryHostSchema() []any {
	return []any{(*InventoryHostModel)(nil)}
}
