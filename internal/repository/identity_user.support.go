package repository

import "time"

const (
	IdentityUserModelStatusActive   = "active"
	IdentityUserModelStatusDisabled = "disabled"
)

type IdentityUserRecord struct {
	ID          int64
	Username    string
	DisplayName string
	Status      string
	DisabledAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func utilIdentityUserRecordFromModel(model *IdentityUserModel) *IdentityUserRecord {
	if model == nil {
		return nil
	}
	return &IdentityUserRecord{
		ID:          model.ID,
		Username:    model.Username,
		DisplayName: model.DisplayName,
		Status:      model.Status,
		DisabledAt:  model.DisabledAt,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
