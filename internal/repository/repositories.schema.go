package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type RepositoriesModel struct {
	bun.BaseModel `bun:"table:repositories,alias:repository_ref"`

	ID           string    `bun:"id,pk,type:uuid" json:"id"`
	OwnerSubject string    `bun:"owner_subject,type:uuid,notnull" json:"ownerSubject"`
	PortGroupID  string    `bun:"port_group_id,type:uuid,notnull" json:"portGroupId"`
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

func (*RepositoriesModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(port_group_id) REFERENCES port_groups (id) ON DELETE CASCADE")
	return nil
}

func RepositoryRefIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_repositories_owner ON repositories(owner_subject)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_repositories_group_url ON repositories(port_group_id, url)`,
	}
}
