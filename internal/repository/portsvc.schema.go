package repository

func PortsvcSchema() []any {
	models := append([]any{}, HostsSchema()...)
	models = append(models, PortAllocationSchema()...)
	models = append(models, ServicesSchema()...)
	models = append(models, DependencySchema()...)
	models = append(models, PortGroupAssetLinksSchema()...)
	return models
}

func PortsvcIndexesSchema() []string {
	statements := append([]string{}, HostsIndexesSchema()...)
	statements = append(statements, PortAllocationIndexesSchema()...)
	statements = append(statements, ServicesIndexesSchema()...)
	statements = append(statements, DependencyIndexesSchema()...)
	statements = append(statements, PortGroupAssetLinksIndexesSchema()...)
	return statements
}
