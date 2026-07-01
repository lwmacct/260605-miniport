package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type PortGroupAssetLinksModel struct {
	bun.BaseModel `bun:"table:port_group_asset_links,alias:asset_link"`

	ID           string                `bun:"id,pk,type:uuid" json:"id"`
	PortGroupID  string                `bun:"port_group_id,type:uuid,notnull" json:"portGroupId"`
	PortGroup    *PortAllocationsModel `bun:"rel:belongs-to,join:port_group_id=id" json:"portGroup,omitempty"`
	PortSlotID   string                `bun:"port_slot_id,type:uuid,nullzero" json:"portSlotId"`
	PortSlot     *ServicesModel        `bun:"rel:belongs-to,join:port_slot_id=id" json:"portSlot,omitempty"`
	AssetID      string                `bun:"asset_id,type:uuid,notnull" json:"assetId"`
	Asset        *DependenciesModel    `bun:"rel:belongs-to,join:asset_id=id" json:"asset,omitempty"`
	RelationType string                `bun:"relation_type,notnull" json:"relationType"`
	Required     bool                  `bun:"required,notnull" json:"required"`
	Notes        string                `bun:"notes" json:"notes"`
	CreatedAt    time.Time             `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt    time.Time             `bun:"updated_at,notnull" json:"updatedAt"`
}

func PortGroupAssetLinksSchema() []any {
	return []any{(*PortGroupAssetLinksModel)(nil)}
}

func (*PortGroupAssetLinksModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(port_group_id) REFERENCES port_groups (id) ON DELETE CASCADE")
	query.ForeignKey("(port_slot_id) REFERENCES port_slots (id) ON DELETE SET NULL")
	query.ForeignKey("(asset_id) REFERENCES dependency_assets (id) ON DELETE CASCADE")
	return nil
}

func PortGroupAssetLinksIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_port_group_asset_links_group ON port_group_asset_links(port_group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_port_group_asset_links_slot ON port_group_asset_links(port_slot_id)`,
		`CREATE INDEX IF NOT EXISTS idx_port_group_asset_links_asset ON port_group_asset_links(asset_id)`,
		`CREATE INDEX IF NOT EXISTS idx_port_group_asset_links_relation ON port_group_asset_links(relation_type)`,
	}
}
