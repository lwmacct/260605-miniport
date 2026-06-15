package authsession

import (
	"time"

	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"table:user_sessions,alias:us"`

	IDHash        []byte     `bun:"id_hash,pk" json:"-"`
	UserID        int64      `bun:"user_id,notnull" json:"userId"`
	LoginIP       string     `bun:"login_ip,notnull" json:"loginIp"`
	LastIP        string     `bun:"last_ip,notnull" json:"lastIp"`
	UserAgentHash []byte     `bun:"user_agent_hash,notnull" json:"-"`
	ExpiresAt     time.Time  `bun:"expires_at,notnull" json:"expiresAt"`
	RevokedAt     *time.Time `bun:"revoked_at,nullzero" json:"revokedAt,omitempty"`
	CreatedAt     time.Time  `bun:"created_at,notnull" json:"createdAt"`
	LastSeenAt    time.Time  `bun:"last_seen_at,notnull" json:"lastSeenAt"`
}
