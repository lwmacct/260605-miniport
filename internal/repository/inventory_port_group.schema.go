package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type InventoryPortGroupModel struct {
	bun.BaseModel `bun:"table:inventory_port_groups,alias:port_group"`

	ID            int64               `bun:"id,pk,autoincrement" json:"id"`
	HostID        int64               `bun:"host_id,notnull" json:"hostId"`
	Host          *InventoryHostModel `bun:"rel:belongs-to,join:host_id=id" json:"host,omitempty"`
	PortStart     int                 `bun:"port_start,notnull" json:"portStart"`
	PortEnd       int                 `bun:"port_end,notnull" json:"portEnd"`
	ServiceName   string              `bun:"service_name,notnull" json:"serviceName"`
	ContainerName string              `bun:"container_name" json:"containerName"`
	DindHost      string              `bun:"dind_host" json:"dindHost"`
	Status        string              `bun:"status,notnull" json:"status"`
	Owner         string              `bun:"owner" json:"owner"`
	Tags          string              `bun:"tags" json:"tags"`
	Notes         string              `bun:"notes" json:"notes"`
	CreatedAt     time.Time           `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt     time.Time           `bun:"updated_at,notnull" json:"updatedAt"`
}

func InventoryPortGroupSchema() []any {
	return []any{(*InventoryPortGroupModel)(nil)}
}

func (*InventoryPortGroupModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(host_id) REFERENCES inventory_hosts (id)")
	return nil
}

func InventoryPortGroupIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_inventory_port_groups_host ON inventory_port_groups(host_id)`,
	}
}
