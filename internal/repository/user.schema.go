package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type UserModel struct {
	bun.BaseModel `bun:"table:users,alias:iu"`

	ID          int64      `bun:"id,pk,autoincrement"`
	Username    string     `bun:"username,notnull,unique"`
	DisplayName string     `bun:"display_name,notnull"`
	Status      string     `bun:"status,notnull"`
	DisabledAt  *time.Time `bun:"disabled_at,nullzero"`
	CreatedAt   time.Time  `bun:"created_at,notnull"`
	UpdatedAt   time.Time  `bun:"updated_at,notnull"`
}

func UserSchema() []any {
	return []any{(*UserModel)(nil)}
}
