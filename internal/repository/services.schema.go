package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type ServicesModel struct {
	bun.BaseModel `bun:"table:port_slots,alias:slot"`

	ID            string                `bun:"id,pk,type:uuid" json:"id"`
	PortGroupID   string                `bun:"port_group_id,type:uuid,notnull" json:"portGroupId"`
	PortGroup     *PortAllocationsModel `bun:"rel:belongs-to,join:port_group_id=id" json:"portGroup,omitempty"`
	Port          int                   `bun:"port,notnull" json:"port"`
	Name          string                `bun:"name,notnull" json:"name"`
	Kind          string                `bun:"kind,notnull" json:"kind"`
	Protocol      string                `bun:"protocol,notnull" json:"protocol"`
	ContainerName string                `bun:"container_name" json:"containerName"`
	Status        string                `bun:"status,notnull" json:"status"`
	Notes         string                `bun:"notes" json:"notes"`
	CreatedAt     time.Time             `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt     time.Time             `bun:"updated_at,notnull" json:"updatedAt"`
}

func ServicesSchema() []any {
	return []any{(*ServicesModel)(nil)}
}

func (*ServicesModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(port_group_id) REFERENCES port_groups (id) ON DELETE CASCADE")
	return nil
}

func ServicesIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_port_slots_group_port ON port_slots(port_group_id, port)`,
		`CREATE INDEX IF NOT EXISTS idx_port_slots_group ON port_slots(port_group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_port_slots_name ON port_slots(name)`,
		`CREATE INDEX IF NOT EXISTS idx_port_slots_status ON port_slots(status)`,
	}
}
