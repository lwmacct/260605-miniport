package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type PortGroupRepositoryLinksModel struct {
	bun.BaseModel `bun:"table:port_group_repository_links,alias:repository_link"`

	ID           string                   `bun:"id,pk,type:uuid" json:"id"`
	PortGroupID  string                   `bun:"port_group_id,type:uuid,notnull" json:"portGroupId"`
	PortGroup    *PortAllocationsModel    `bun:"rel:belongs-to,join:port_group_id=id" json:"portGroup,omitempty"`
	PortSlotID   string                   `bun:"port_slot_id,type:uuid,nullzero" json:"portSlotId"`
	PortSlot     *ServicesModel           `bun:"rel:belongs-to,join:port_slot_id=id" json:"portSlot,omitempty"`
	RepositoryID string                   `bun:"repository_id,type:uuid,notnull" json:"repositoryId"`
	Repository   *GithubRepositoriesModel `bun:"rel:belongs-to,join:repository_id=id" json:"repository,omitempty"`
	RelationType string                   `bun:"relation_type,notnull" json:"relationType"`
	Required     bool                     `bun:"required,notnull" json:"required"`
	Notes        string                   `bun:"notes" json:"notes"`
	CreatedAt    time.Time                `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt    time.Time                `bun:"updated_at,notnull" json:"updatedAt"`
}

func PortGroupRepositoryLinksSchema() []any { return []any{(*PortGroupRepositoryLinksModel)(nil)} }

func (*PortGroupRepositoryLinksModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(port_group_id) REFERENCES port_groups (id) ON DELETE CASCADE")
	query.ForeignKey("(port_slot_id) REFERENCES port_slots (id) ON DELETE SET NULL")
	query.ForeignKey("(repository_id) REFERENCES github_repositories (id) ON DELETE RESTRICT")
	return nil
}

func PortGroupRepositoryLinksIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_port_group_repository_links_group ON port_group_repository_links(port_group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_port_group_repository_links_slot ON port_group_repository_links(port_slot_id)`,
		`CREATE INDEX IF NOT EXISTS idx_port_group_repository_links_repository ON port_group_repository_links(repository_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_port_group_repository_links_unique ON port_group_repository_links(port_group_id, port_slot_id, repository_id, relation_type)`,
	}
}
