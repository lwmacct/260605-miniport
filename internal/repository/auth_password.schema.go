package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type AuthPasswordModel struct {
	bun.BaseModel `bun:"table:auth_passwords,alias:ap"`

	IdentityUserID int64     `bun:"identity_user_id,pk"`
	PasswordHash   string    `bun:"password_hash,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull"`
	UpdatedAt      time.Time `bun:"updated_at,notnull"`
}

func AuthPasswordSchema() []any {
	return []any{(*AuthPasswordModel)(nil)}
}

func (*AuthPasswordModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(identity_user_id) REFERENCES identity_users (id) ON DELETE CASCADE")
	return nil
}
