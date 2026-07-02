package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type PortAllocationsModel struct {
	bun.BaseModel `bun:"table:port_groups,alias:port_group"`

	ID               string      `bun:"id,pk,type:uuid" json:"id"`
	OwnerSubject     string      `bun:"owner_subject,type:uuid,notnull" json:"ownerSubject"`
	HostID           string      `bun:"host_id,type:uuid,nullzero" json:"hostId"`
	Host             *HostsModel `bun:"rel:belongs-to,join:host_id=id" json:"host,omitempty"`
	PortPrefix       int         `bun:"port_prefix,notnull" json:"portPrefix"`
	EnvironmentName  string      `bun:"environment_name" json:"environmentName"`
	EnvironmentOwner string      `bun:"environment_owner" json:"environmentOwner"`
	RuntimeMode      string      `bun:"runtime_mode,notnull" json:"runtimeMode"`
	RuntimeName      string      `bun:"runtime_name" json:"runtimeName"`
	ServiceIP        string      `bun:"service_ip" json:"serviceIp"`
	Status           string      `bun:"status,notnull" json:"status"`
	Tags             string      `bun:"tags" json:"tags"`
	Notes            string      `bun:"notes" json:"notes"`
	CreatedAt        time.Time   `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt        time.Time   `bun:"updated_at,notnull" json:"updatedAt"`
}

func PortAllocationSchema() []any {
	return []any{(*PortAllocationsModel)(nil)}
}

func (*PortAllocationsModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(host_id) REFERENCES hosts (id) ON DELETE SET NULL")
	return nil
}

func PortAllocationIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_port_groups_owner_port_prefix ON port_groups(owner_subject, port_prefix)`,
		`CREATE INDEX IF NOT EXISTS idx_port_groups_owner ON port_groups(owner_subject)`,
		`CREATE INDEX IF NOT EXISTS idx_port_groups_host ON port_groups(host_id)`,
		`CREATE INDEX IF NOT EXISTS idx_port_groups_status ON port_groups(status)`,
		`CREATE INDEX IF NOT EXISTS idx_port_groups_environment_name ON port_groups(environment_name)`,
	}
}
