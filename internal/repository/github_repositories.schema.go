package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type GithubRepositoriesModel struct {
	bun.BaseModel `bun:"table:github_repositories,alias:github_repository"`

	ID                 string                    `bun:"id,pk,type:uuid" json:"id"`
	InstallationID     string                    `bun:"installation_id,type:uuid,notnull" json:"installationId"`
	Installation       *GithubInstallationsModel `bun:"rel:belongs-to,join:installation_id=id" json:"installation,omitempty"`
	GithubRepositoryID int64                     `bun:"github_repository_id,notnull" json:"githubRepositoryId"`
	NodeID             string                    `bun:"node_id" json:"nodeId"`
	OwnerLogin         string                    `bun:"owner_login,notnull" json:"ownerLogin"`
	Name               string                    `bun:"name,notnull" json:"name"`
	FullName           string                    `bun:"full_name,notnull" json:"fullName"`
	HTMLURL            string                    `bun:"html_url,notnull" json:"htmlUrl"`
	Description        string                    `bun:"description" json:"description"`
	DefaultBranch      string                    `bun:"default_branch" json:"defaultBranch"`
	Visibility         string                    `bun:"visibility,notnull" json:"visibility"`
	Private            bool                      `bun:"private,notnull" json:"private"`
	Fork               bool                      `bun:"fork,notnull" json:"fork"`
	Archived           bool                      `bun:"archived,notnull" json:"archived"`
	Disabled           bool                      `bun:"disabled,notnull" json:"disabled"`
	State              string                    `bun:"state,notnull" json:"state"`
	PushedAt           time.Time                 `bun:"pushed_at,nullzero" json:"pushedAt"`
	RemoteUpdatedAt    time.Time                 `bun:"remote_updated_at,nullzero" json:"remoteUpdatedAt"`
	LastSeenAt         time.Time                 `bun:"last_seen_at,notnull" json:"lastSeenAt"`
	CreatedAt          time.Time                 `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt          time.Time                 `bun:"updated_at,notnull" json:"updatedAt"`
}

func GithubRepositoriesSchema() []any { return []any{(*GithubRepositoriesModel)(nil)} }

func (*GithubRepositoriesModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(installation_id) REFERENCES github_installations (id) ON DELETE CASCADE")
	return nil
}

func GithubRepositoriesIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_github_repositories_external ON github_repositories(github_repository_id)`,
		`CREATE INDEX IF NOT EXISTS idx_github_repositories_installation ON github_repositories(installation_id)`,
		`CREATE INDEX IF NOT EXISTS idx_github_repositories_full_name ON github_repositories(full_name)`,
		`CREATE INDEX IF NOT EXISTS idx_github_repositories_state ON github_repositories(state)`,
	}
}
