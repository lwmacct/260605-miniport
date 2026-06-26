package repository

func InventorySchema() []any {
	models := append([]any{}, InventoryHostSchema()...)
	models = append(models, InventoryPortGroupSchema()...)
	models = append(models, InventoryPortSlotSchema()...)
	models = append(models, InventoryComponentSchema()...)
	models = append(models, InventoryRepositoryRefSchema()...)
	return models
}

func InventoryIndexesSchema() []string {
	statements := append([]string{}, InventoryPortGroupIndexesSchema()...)
	statements = append(statements, InventoryPortSlotIndexesSchema()...)
	statements = append(statements, InventoryComponentIndexesSchema()...)
	statements = append(statements, InventoryRepositoryRefIndexesSchema()...)
	return statements
}
