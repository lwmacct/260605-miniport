package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type HostsModel struct {
	bun.BaseModel `bun:"table:hosts,alias:host"`

	ID        string    `bun:"id,pk,type:uuid" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	IP        string    `bun:"ip" json:"ip"`
	Spec      string    `bun:"spec" json:"spec"`
	Status    string    `bun:"status,notnull" json:"status"`
	Notes     string    `bun:"notes" json:"notes"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func HostsSchema() []any {
	return []any{(*HostsModel)(nil)}
}

func HostsIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_hosts_name ON hosts(name)`,
		`CREATE INDEX IF NOT EXISTS idx_hosts_ip ON hosts(ip)`,
		`CREATE INDEX IF NOT EXISTS idx_hosts_status ON hosts(status)`,
	}
}
