package authpassword

import (
	"time"

	"github.com/uptrace/bun"
)

type Password struct {
	bun.BaseModel `bun:"table:user_passwords,alias:up"`

	UserID       int64     `bun:"user_id,pk" json:"userId"`
	PasswordHash string    `bun:"password_hash,notnull" json:"-"`
	CreatedAt    time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt    time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}
