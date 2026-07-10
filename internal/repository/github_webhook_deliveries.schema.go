package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type GithubWebhookDeliveriesModel struct {
	bun.BaseModel `bun:"table:github_webhook_deliveries,alias:github_webhook_delivery"`

	DeliveryID     string    `bun:"delivery_id,pk" json:"deliveryId"`
	Event          string    `bun:"event,notnull" json:"event"`
	Action         string    `bun:"action" json:"action"`
	InstallationID int64     `bun:"installation_id" json:"installationId"`
	Status         string    `bun:"status,notnull" json:"status"`
	Error          string    `bun:"error" json:"error"`
	ReceivedAt     time.Time `bun:"received_at,notnull" json:"receivedAt"`
	ProcessedAt    time.Time `bun:"processed_at,nullzero" json:"processedAt"`
}

func GithubWebhookDeliveriesSchema() []any { return []any{(*GithubWebhookDeliveriesModel)(nil)} }

func GithubWebhookDeliveriesIndexesSchema() []string {
	return []string{
		`CREATE INDEX IF NOT EXISTS idx_github_webhook_deliveries_received ON github_webhook_deliveries(received_at)`,
		`CREATE INDEX IF NOT EXISTS idx_github_webhook_deliveries_status ON github_webhook_deliveries(status)`,
	}
}
