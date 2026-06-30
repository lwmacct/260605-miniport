package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type DependenciesModel struct {
	bun.BaseModel `bun:"table:dependencies,alias:dependency"`

	ID           string    `bun:"id,pk,type:uuid" json:"id"`
	OwnerSubject string    `bun:"owner_subject,type:uuid,notnull" json:"ownerSubject"`
	Name         string    `bun:"name,notnull" json:"name"`
	Type         string    `bun:"type,notnull" json:"type"`
	URL          string    `bun:"url" json:"url"`
	Version      string    `bun:"version" json:"version"`
	Notes        string    `bun:"notes" json:"notes"`
	CreatedAt    time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt    time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func DependencySchema() []any {
	return []any{(*DependenciesModel)(nil)}
}

func DependencyIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_dependencies_owner ON dependencies(owner_subject)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_dependencies_owner_name_type_version ON dependencies(owner_subject, name, type, version)`,
	}
}
