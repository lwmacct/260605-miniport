package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type PortAllocationsModel struct {
	bun.BaseModel `bun:"table:port_allocations,alias:allocation"`

	ID        int64      `bun:"id,pk,autoincrement" json:"id"`
	UserID    int64      `bun:"user_id,notnull" json:"userId"`
	User      *UserModel `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
	PortStart int        `bun:"port_start,notnull" json:"portStart"`
	PortEnd   int        `bun:"port_end,notnull" json:"portEnd"`
	Status    string     `bun:"status,notnull" json:"status"`
	Notes     string     `bun:"notes" json:"notes"`
	CreatedAt time.Time  `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt time.Time  `bun:"updated_at,notnull" json:"updatedAt"`
}

func PortAllocationSchema() []any {
	return []any{(*PortAllocationsModel)(nil)}
}

func (*PortAllocationsModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(user_id) REFERENCES users (id) ON DELETE CASCADE")
	return nil
}

func PortAllocationIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_port_allocations_user_port_start ON port_allocations(user_id, port_start)`,
		`CREATE INDEX IF NOT EXISTS idx_port_allocations_user ON port_allocations(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_port_allocations_status ON port_allocations(status)`,
	}
}
