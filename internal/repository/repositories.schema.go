package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type RepositoriesModel struct {
	bun.BaseModel `bun:"table:repositories,alias:repository_ref"`

	ID           string    `bun:"id,pk,type:uuid" json:"id"`
	OwnerSubject string    `bun:"owner_subject,type:uuid,notnull" json:"ownerSubject"`
	Name         string    `bun:"name,notnull" json:"name"`
	URL          string    `bun:"url,notnull" json:"url"`
	Kind         string    `bun:"kind,notnull" json:"kind"`
	Notes        string    `bun:"notes" json:"notes"`
	CreatedAt    time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt    time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func RepositoryRefSchema() []any {
	return []any{(*RepositoriesModel)(nil)}
}

func RepositoryRefIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_repositories_owner ON repositories(owner_subject)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_repositories_owner_url ON repositories(owner_subject, url)`,
	}
}
