package inventory

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
)

func ApplySchema(ctx context.Context, db *bun.DB) error {
	if db.Dialect().Name() == dialect.SQLite {
		if _, err := db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
			return fmt.Errorf("enable sqlite foreign keys: %w", err)
		}
	}

	models := []any{
		(*Host)(nil),
		(*PortGroup)(nil),
		(*PortSlot)(nil),
		(*Component)(nil),
		(*Repository)(nil),
	}
	for _, model := range models {
		if _, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			return fmt.Errorf("create inventory table: %w", err)
		}
	}

	statements := []string{
		`CREATE INDEX IF NOT EXISTS idx_port_groups_host ON port_groups(host_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_port_slots_group_port ON port_slots(port_group_id, port)`,
		`CREATE INDEX IF NOT EXISTS idx_components_group ON components(port_group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_repositories_group ON repositories(port_group_id)`,
	}
	for _, stmt := range statements {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("create inventory index: %w", err)
		}
	}
	return nil
}
