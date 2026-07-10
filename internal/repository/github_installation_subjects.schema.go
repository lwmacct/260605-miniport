package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type GithubInstallationSubjectsModel struct {
	bun.BaseModel `bun:"table:github_installation_subjects,alias:github_installation_subject"`

	ID             string                    `bun:"id,pk,type:uuid" json:"id"`
	InstallationID string                    `bun:"installation_id,type:uuid,notnull" json:"installationId"`
	Installation   *GithubInstallationsModel `bun:"rel:belongs-to,join:installation_id=id" json:"installation,omitempty"`
	OwnerSubject   string                    `bun:"owner_subject,type:uuid,notnull" json:"ownerSubject"`
	CreatedAt      time.Time                 `bun:"created_at,notnull" json:"createdAt"`
}

func GithubInstallationSubjectsSchema() []any { return []any{(*GithubInstallationSubjectsModel)(nil)} }

func (*GithubInstallationSubjectsModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(installation_id) REFERENCES github_installations (id) ON DELETE CASCADE")
	return nil
}

func GithubInstallationSubjectsIndexesSchema() []string {
	return []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_github_installation_subject_unique ON github_installation_subjects(installation_id, owner_subject)`,
		`CREATE INDEX IF NOT EXISTS idx_github_installation_subject_owner ON github_installation_subjects(owner_subject)`,
	}
}
