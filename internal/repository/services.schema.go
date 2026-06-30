package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type ServicesModel struct {
	bun.BaseModel `bun:"table:services,alias:service"`

	ID               string                `bun:"id,pk,type:uuid" json:"id"`
	OwnerSubject     string                `bun:"owner_subject,type:uuid,notnull" json:"ownerSubject"`
	PortAllocationID string                `bun:"port_allocation_id,type:uuid,nullzero" json:"portAllocationId"`
	PortAllocation   *PortAllocationsModel `bun:"rel:belongs-to,join:port_allocation_id=id" json:"portAllocation,omitempty"`
	Name             string                `bun:"name,notnull" json:"name"`
	ProjectName      string                `bun:"project_name" json:"projectName"`
	DindIP           string                `bun:"dind_ip" json:"dindIp"`
	DindContainer    string                `bun:"dind_container" json:"dindContainer"`
	Status           string                `bun:"status,notnull" json:"status"`
	Owner            string                `bun:"owner" json:"owner"`
	Tags             string                `bun:"tags" json:"tags"`
	Notes            string                `bun:"notes" json:"notes"`
	CreatedAt        time.Time             `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt        time.Time             `bun:"updated_at,notnull" json:"updatedAt"`
}

func ServicesSchema() []any {
	return []any{(*ServicesModel)(nil)}
}

func (*ServicesModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(port_allocation_id) REFERENCES port_allocations (id) ON DELETE SET NULL")
	return nil
}

func ServicesIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_services_owner ON services(owner_subject)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_services_port_allocation ON services(port_allocation_id) WHERE port_allocation_id IS NOT NULL`,
		`CREATE INDEX IF NOT EXISTS idx_services_project_name ON services(project_name)`,
		`CREATE INDEX IF NOT EXISTS idx_services_status ON services(status)`,
	}
}
