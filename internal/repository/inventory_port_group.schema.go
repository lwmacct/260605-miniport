package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type InventoryPortGroupModel struct {
	bun.BaseModel `bun:"table:port_allocations,alias:allocation"`

	ID            int64      `bun:"id,pk,autoincrement" json:"id"`
	UserID        int64      `bun:"user_id,notnull" json:"userId"`
	User          *UserModel `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
	PortStart     int        `bun:"port_start,notnull" json:"portStart"`
	PortEnd       int        `bun:"port_end,notnull" json:"portEnd"`
	Name          string     `bun:"name,notnull" json:"name"`
	DindIP        string     `bun:"dind_ip" json:"dindIp"`
	DindContainer string     `bun:"dind_container" json:"dindContainer"`
	Status        string     `bun:"status,notnull" json:"status"`
	Owner         string     `bun:"owner" json:"owner"`
	Tags          string     `bun:"tags" json:"tags"`
	Notes         string     `bun:"notes" json:"notes"`
	CreatedAt     time.Time  `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt     time.Time  `bun:"updated_at,notnull" json:"updatedAt"`
}

func InventoryPortGroupSchema() []any {
	return []any{(*InventoryPortGroupModel)(nil)}
}

func (*InventoryPortGroupModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(user_id) REFERENCES users (id) ON DELETE CASCADE")
	return nil
}

func InventoryPortGroupIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_port_allocations_user_port_start ON port_allocations(user_id, port_start)`,
		`CREATE INDEX IF NOT EXISTS idx_port_allocations_user ON port_allocations(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_port_allocations_status ON port_allocations(status)`,
		`CREATE INDEX IF NOT EXISTS idx_port_allocations_dind_ip ON port_allocations(dind_ip)`,
	}
}
