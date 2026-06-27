package repository

func InventorySchema() []any {
	models := append([]any{}, InventoryPortGroupSchema()...)
	models = append(models, InventoryPortSlotSchema()...)
	models = append(models, InventoryProjectSchema()...)
	models = append(models, InventoryComponentSchema()...)
	models = append(models, InventoryRepositoryRefSchema()...)
	return models
}

func InventoryIndexesSchema() []string {
	statements := append([]string{}, InventoryPortGroupIndexesSchema()...)
	statements = append(statements, InventoryPortSlotIndexesSchema()...)
	statements = append(statements, InventoryProjectIndexesSchema()...)
	statements = append(statements, InventoryComponentIndexesSchema()...)
	statements = append(statements, InventoryRepositoryRefIndexesSchema()...)
	return statements
}
