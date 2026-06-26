package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type AuthSessionModel struct {
	bun.BaseModel `bun:"table:auth_sessions,alias:as"`

	IDHash        string     `bun:"id_hash,pk"`
	UserID        int64      `bun:"user_id,notnull"`
	LoginIP       string     `bun:"login_ip"`
	LastIP        string     `bun:"last_ip"`
	UserAgentHash string     `bun:"user_agent_hash"`
	ExpiresAt     time.Time  `bun:"expires_at,notnull"`
	RevokedAt     *time.Time `bun:"revoked_at,nullzero"`
	CreatedAt     time.Time  `bun:"created_at,notnull"`
	LastSeenAt    time.Time  `bun:"last_seen_at,notnull"`
}

func AuthSessionSchema() []any {
	return []any{(*AuthSessionModel)(nil)}
}
