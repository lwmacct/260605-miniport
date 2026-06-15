package identityuser

import (
	"time"

	"github.com/uptrace/bun"
)

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID          int64      `bun:"id,pk,autoincrement" json:"id"`
	Username    string     `bun:"username,notnull,unique" json:"username"`
	DisplayName string     `bun:"display_name,notnull" json:"displayName"`
	Status      string     `bun:"status,notnull" json:"status"`
	DisabledAt  *time.Time `bun:"disabled_at,nullzero" json:"disabledAt,omitempty"`
	CreatedAt   time.Time  `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time  `bun:"updated_at,notnull" json:"updatedAt"`
}
