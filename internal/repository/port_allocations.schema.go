package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type PortAllocationsModel struct {
	bun.BaseModel `bun:"table:port_allocations,alias:allocation"`

	ID           string    `bun:"id,pk,type:uuid" json:"id"`
	OwnerSubject string    `bun:"owner_subject,type:uuid,notnull" json:"ownerSubject"`
	PortStart    int       `bun:"port_start,notnull" json:"portStart"`
	PortEnd      int       `bun:"port_end,notnull" json:"portEnd"`
	Status       string    `bun:"status,notnull" json:"status"`
	Notes        string    `bun:"notes" json:"notes"`
	CreatedAt    time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt    time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func PortAllocationSchema() []any {
	return []any{(*PortAllocationsModel)(nil)}
}

func PortAllocationIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_port_allocations_owner_port_start ON port_allocations(owner_subject, port_start)`,
		`CREATE INDEX IF NOT EXISTS idx_port_allocations_owner ON port_allocations(owner_subject)`,
		`CREATE INDEX IF NOT EXISTS idx_port_allocations_status ON port_allocations(status)`,
	}
}
