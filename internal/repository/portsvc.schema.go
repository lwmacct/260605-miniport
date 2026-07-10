package repository

func PortsvcSchema() []any {
	models := append([]any{}, HostsSchema()...)
	models = append(models, PortAllocationSchema()...)
	models = append(models, ServicesSchema()...)
	models = append(models, ServiceGroupsSchema()...)
	models = append(models, DependencySchema()...)
	models = append(models, PortGroupAssetLinksSchema()...)
	models = append(models, PortGroupRepositoryLinksSchema()...)
	return models
}

func PortsvcIndexesSchema() []string {
	statements := append([]string{}, HostsIndexesSchema()...)
	statements = append(statements, PortAllocationIndexesSchema()...)
	statements = append(statements, ServicesIndexesSchema()...)
	statements = append(statements, ServiceGroupsIndexesSchema()...)
	statements = append(statements, DependencyIndexesSchema()...)
	statements = append(statements, PortGroupAssetLinksIndexesSchema()...)
	statements = append(statements, PortGroupRepositoryLinksIndexesSchema()...)
	return statements
}
