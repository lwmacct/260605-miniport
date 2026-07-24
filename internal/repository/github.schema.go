package repository

func GithubSchema() []any {
	models := append([]any{}, GithubInstallationsSchema()...)
	models = append(models, GithubRepositoriesSchema()...)
	models = append(models, GithubConnectionStatesSchema()...)
	models = append(models, GithubWebhookDeliveriesSchema()...)
	return models
}

func GithubIndexesSchema() []string {
	statements := append([]string{}, GithubInstallationsIndexesSchema()...)
	statements = append(statements, GithubRepositoriesIndexesSchema()...)
	statements = append(statements, GithubConnectionStatesIndexesSchema()...)
	statements = append(statements, GithubWebhookDeliveriesIndexesSchema()...)
	return statements
}
