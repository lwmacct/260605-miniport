package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type GithubConnectionStatesModel struct {
	bun.BaseModel `bun:"table:github_connection_states,alias:github_connection_state"`

	ID        string    `bun:"id,pk,type:uuid" json:"id"`
	StateHash string    `bun:"state_hash,notnull" json:"stateHash"`
	ExpiresAt time.Time `bun:"expires_at,notnull" json:"expiresAt"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"createdAt"`
}

func GithubConnectionStatesSchema() []any { return []any{(*GithubConnectionStatesModel)(nil)} }

func GithubConnectionStatesIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_github_connection_states_hash ON github_connection_states(state_hash)`,
		`CREATE INDEX IF NOT EXISTS idx_github_connection_states_expiry ON github_connection_states(expires_at)`,
	}
}
