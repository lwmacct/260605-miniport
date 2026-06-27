package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type RepositoriesModel struct {
	bun.BaseModel `bun:"table:repositories,alias:repository_ref"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID    int64     `bun:"user_id,notnull" json:"userId"`
	Name      string    `bun:"name,notnull" json:"name"`
	URL       string    `bun:"url,notnull" json:"url"`
	Kind      string    `bun:"kind,notnull" json:"kind"`
	Notes     string    `bun:"notes" json:"notes"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func RepositoryRefSchema() []any {
	return []any{(*RepositoriesModel)(nil)}
}

func (*RepositoriesModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(user_id) REFERENCES users (id) ON DELETE CASCADE")
	return nil
}

func RepositoryRefIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_repositories_user ON repositories(user_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_repositories_user_url ON repositories(user_id, url)`,
	}
}
