package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type DependenciesModel struct {
	bun.BaseModel `bun:"table:dependency_assets,alias:asset"`

	ID              string    `bun:"id,pk,type:uuid" json:"id"`
	OwnerSubject    string    `bun:"owner_subject,type:uuid,notnull" json:"ownerSubject"`
	Name            string    `bun:"name,notnull" json:"name"`
	AssetKind       string    `bun:"asset_kind,notnull" json:"assetKind"`
	AssetType       string    `bun:"asset_type,notnull" json:"assetType"`
	Provider        string    `bun:"provider,notnull" json:"provider"`
	URL             string    `bun:"url" json:"url"`
	FullName        string    `bun:"full_name" json:"fullName"`
	ExternalID      string    `bun:"external_id" json:"externalId"`
	Visibility      string    `bun:"visibility,notnull" json:"visibility"`
	Controllability string    `bun:"controllability,notnull" json:"controllability"`
	Status          string    `bun:"status,notnull" json:"status"`
	Description     string    `bun:"description" json:"description"`
	Metadata        string    `bun:"metadata" json:"metadata"`
	LastSyncedAt    time.Time `bun:"last_synced_at,nullzero" json:"lastSyncedAt"`
	Notes           string    `bun:"notes" json:"notes"`
	CreatedAt       time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt       time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

func DependencySchema() []any {
	return []any{(*DependenciesModel)(nil)}
}

func DependencyIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_dependency_assets_owner ON dependency_assets(owner_subject)`,
		`CREATE INDEX IF NOT EXISTS idx_dependency_assets_kind ON dependency_assets(asset_kind)`,
		`CREATE INDEX IF NOT EXISTS idx_dependency_assets_type ON dependency_assets(asset_type)`,
		`CREATE INDEX IF NOT EXISTS idx_dependency_assets_provider ON dependency_assets(provider)`,
		`CREATE INDEX IF NOT EXISTS idx_dependency_assets_status ON dependency_assets(status)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_dependency_assets_owner_name_kind ON dependency_assets(owner_subject, name, asset_kind)`,
	}
}
