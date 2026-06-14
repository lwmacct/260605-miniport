package inventory

func Schema() []any {
	return []any{
		(*Host)(nil),
		(*PortGroup)(nil),
		(*PortSlot)(nil),
		(*Component)(nil),
		(*Repository)(nil),
	}
}

func IndexStatements() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_port_groups_host ON port_groups(host_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_port_slots_group_port ON port_slots(port_group_id, port)`,
		`CREATE INDEX IF NOT EXISTS idx_components_group ON components(port_group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_repositories_group ON repositories(port_group_id)`,
	}
}
