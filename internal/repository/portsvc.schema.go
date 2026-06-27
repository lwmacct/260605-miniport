package repository

func PortsvcSchema() []any {
	models := append([]any{}, PortAllocationSchema()...)
	models = append(models, ServicesSchema()...)
	models = append(models, DependencySchema()...)
	models = append(models, RepositoryRefSchema()...)
	models = append(models, ServiceReposSchema()...)
	models = append(models, ServiceDepsSchema()...)
	return models
}

func PortsvcIndexesSchema() []string {
	statements := append([]string{}, PortAllocationIndexesSchema()...)
	statements = append(statements, ServicesIndexesSchema()...)
	statements = append(statements, DependencyIndexesSchema()...)
	statements = append(statements, RepositoryRefIndexesSchema()...)
	statements = append(statements, ServiceReposIndexesSchema()...)
	statements = append(statements, ServiceDepsIndexesSchema()...)
	return statements
}
