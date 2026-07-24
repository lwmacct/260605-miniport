package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type ServiceGroupsModel struct {
	bun.BaseModel `bun:"table:service_groups,alias:service_group"`

	ID          string    `bun:"id,pk,type:uuid" json:"id"`
	Name        string    `bun:"name,notnull" json:"name"`
	Kind        string    `bun:"kind,notnull" json:"kind"`
	Status      string    `bun:"status,notnull" json:"status"`
	Description string    `bun:"description" json:"description"`
	Notes       string    `bun:"notes" json:"notes"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type ServiceGroupsPortGroupsModel struct {
	bun.BaseModel `bun:"table:service_group_port_groups,alias:service_group_port_group"`

	ID             string                `bun:"id,pk,type:uuid" json:"id"`
	ServiceGroupID string                `bun:"service_group_id,type:uuid,notnull" json:"serviceGroupId"`
	ServiceGroup   *ServiceGroupsModel   `bun:"rel:belongs-to,join:service_group_id=id" json:"serviceGroup,omitempty"`
	PortGroupID    string                `bun:"port_group_id,type:uuid,notnull" json:"portGroupId"`
	PortGroup      *PortAllocationsModel `bun:"rel:belongs-to,join:port_group_id=id" json:"portGroup,omitempty"`
	Role           string                `bun:"role" json:"role"`
	Notes          string                `bun:"notes" json:"notes"`
	CreatedAt      time.Time             `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt      time.Time             `bun:"updated_at,notnull" json:"updatedAt"`
}

func ServiceGroupsSchema() []any {
	return []any{(*ServiceGroupsModel)(nil), (*ServiceGroupsPortGroupsModel)(nil)}
}

func (*ServiceGroupsPortGroupsModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(service_group_id) REFERENCES service_groups (id) ON DELETE CASCADE")
	query.ForeignKey("(port_group_id) REFERENCES port_groups (id) ON DELETE CASCADE")
	return nil
}

func ServiceGroupsIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_service_groups_status ON service_groups(status)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_service_groups_name ON service_groups(name)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_service_group_port_groups_unique ON service_group_port_groups(service_group_id, port_group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_service_group_port_groups_service_group ON service_group_port_groups(service_group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_service_group_port_groups_port_group ON service_group_port_groups(port_group_id)`,
	}
}
