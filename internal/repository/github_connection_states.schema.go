package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type GithubConnectionStatesModel struct {
	bun.BaseModel `bun:"table:github_connection_states,alias:github_connection_state"`

	StateHash    string    `bun:"state_hash,pk" json:"stateHash"`
	OwnerSubject string    `bun:"owner_subject,type:uuid,notnull" json:"ownerSubject"`
	ExpiresAt    time.Time `bun:"expires_at,notnull" json:"expiresAt"`
	CreatedAt    time.Time `bun:"created_at,notnull" json:"createdAt"`
}

func GithubConnectionStatesSchema() []any { return []any{(*GithubConnectionStatesModel)(nil)} }

func GithubConnectionStatesIndexesSchema() []string {
	return []string{`CREATE INDEX IF NOT EXISTS idx_github_connection_states_expiry ON github_connection_states(expires_at)`}
}
