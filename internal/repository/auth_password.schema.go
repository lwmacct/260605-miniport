package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type AuthPasswordModel struct {
	bun.BaseModel `bun:"table:auth_passwords,alias:ap"`

	UserID       int64     `bun:"user_id,pk"`
	PasswordHash string    `bun:"password_hash,notnull"`
	CreatedAt    time.Time `bun:"created_at,notnull"`
	UpdatedAt    time.Time `bun:"updated_at,notnull"`
}

func AuthPasswordSchema() []any {
	return []any{(*AuthPasswordModel)(nil)}
}
