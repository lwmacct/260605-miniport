package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type DependenciesModel struct {
	bun.BaseModel `bun:"table:dependencies,alias:dependency"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID    int64     `bun:"user_id,notnull" json:"userId"`
	Name      string    `bun:"name,notnull" json:"name"`
	Type      string    `bun:"type,notnull" json:"type"`
	URL       string    `bun:"url" json:"url"`
	Version   string    `bun:"version" json:"version"`
	Notes     string    `bun:"notes" json:"notes"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func DependencySchema() []any {
	return []any{(*DependenciesModel)(nil)}
}

func (*DependenciesModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(user_id) REFERENCES users (id) ON DELETE CASCADE")
	return nil
}

func DependencyIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_dependencies_user ON dependencies(user_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_dependencies_user_name_type_version ON dependencies(user_id, name, type, version)`,
	}
}
