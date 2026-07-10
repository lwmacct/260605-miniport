package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type GithubInstallationsModel struct {
	bun.BaseModel `bun:"table:github_installations,alias:github_installation"`

	ID                   string    `bun:"id,pk,type:uuid" json:"id"`
	GithubInstallationID int64     `bun:"github_installation_id,notnull" json:"githubInstallationId"`
	AccountID            int64     `bun:"account_id,notnull" json:"accountId"`
	AccountLogin         string    `bun:"account_login,notnull" json:"accountLogin"`
	AccountType          string    `bun:"account_type,notnull" json:"accountType"`
	AvatarURL            string    `bun:"avatar_url" json:"avatarUrl"`
	RepositorySelection  string    `bun:"repository_selection,notnull" json:"repositorySelection"`
	Permissions          string    `bun:"permissions,notnull" json:"permissions"`
	Status               string    `bun:"status,notnull" json:"status"`
	SuspendedAt          time.Time `bun:"suspended_at,nullzero" json:"suspendedAt"`
	LastSyncedAt         time.Time `bun:"last_synced_at,nullzero" json:"lastSyncedAt"`
	LastSyncError        string    `bun:"last_sync_error" json:"lastSyncError"`
	CreatedAt            time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt            time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func GithubInstallationsSchema() []any { return []any{(*GithubInstallationsModel)(nil)} }

func GithubInstallationsIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_github_installations_external ON github_installations(github_installation_id)`,
		`CREATE INDEX IF NOT EXISTS idx_github_installations_account ON github_installations(account_id)`,
		`CREATE INDEX IF NOT EXISTS idx_github_installations_status ON github_installations(status)`,
	}
}
