package repository

import (
	"strings"
	"time"
)

type InventoryPortGroupListFilter struct {
	UserID      int64
	Admin       bool
	Query       string
	Sort        string
	Status      string
	ProjectName string
	DindIP      string
}

type InventoryPortGroupRecord struct {
	ID            int64
	UserID        int64
	Username      string
	PortStart     int
	PortEnd       int
	Name          string
	DindIP        string
	DindContainer string
	Status        string
	Owner         string
	Tags          string
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type InventoryProjectRecord struct {
	ID          int64
	PortGroupID int64
	Name        string
	Description string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InventoryPortSlotRecord struct {
	ID          int64
	PortGroupID int64
	Port        int
	Name        string
	Protocol    string
	Purpose     string
	Status      string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InventoryComponentRecord struct {
	ID          int64
	PortGroupID int64
	Name        string
	Type        string
	URL         string
	Version     string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InventoryRepositoryRefRecord struct {
	ID          int64
	PortGroupID int64
	ProjectID   int64
	Name        string
	URL         string
	Kind        string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InventoryPortGroupChildrenRecord struct {
	PortGroupID  int64
	Slots        []InventoryPortSlotRecord
	Projects     []InventoryProjectRecord
	Components   []InventoryComponentRecord
	Repositories []InventoryRepositoryRefRecord
}

func utilInventoryPortGroupRecordFromModel(model *InventoryPortGroupModel) *InventoryPortGroupRecord {
	if model == nil {
		return nil
	}
	out := &InventoryPortGroupRecord{
		ID:            model.ID,
		UserID:        model.UserID,
		PortStart:     model.PortStart,
		PortEnd:       model.PortEnd,
		Name:          model.Name,
		DindIP:        model.DindIP,
		DindContainer: model.DindContainer,
		Status:        model.Status,
		Owner:         model.Owner,
		Tags:          model.Tags,
		Notes:         model.Notes,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
	if model.User != nil {
		out.Username = model.User.Username
	}
	return out
}

func utilInventoryProjectRecordFromModel(model *InventoryProjectModel) *InventoryProjectRecord {
	if model == nil {
		return nil
	}
	return &InventoryProjectRecord{
		ID:          model.ID,
		PortGroupID: model.PortGroupID,
		Name:        model.Name,
		Description: model.Description,
		Notes:       model.Notes,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func utilInventoryPortSlotRecordFromModel(model *InventoryPortSlotModel) *InventoryPortSlotRecord {
	if model == nil {
		return nil
	}
	return &InventoryPortSlotRecord{
		ID:          model.ID,
		PortGroupID: model.PortGroupID,
		Port:        model.Port,
		Name:        model.Name,
		Protocol:    model.Protocol,
		Purpose:     model.Purpose,
		Status:      model.Status,
		Notes:       model.Notes,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func utilInventoryComponentRecordFromModel(model *InventoryComponentModel) *InventoryComponentRecord {
	if model == nil {
		return nil
	}
	return &InventoryComponentRecord{
		ID:          model.ID,
		PortGroupID: model.PortGroupID,
		Name:        model.Name,
		Type:        model.Type,
		URL:         model.URL,
		Version:     model.Version,
		Notes:       model.Notes,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func utilInventoryRepositoryRefRecordFromModel(model *InventoryRepositoryRefModel) *InventoryRepositoryRefRecord {
	if model == nil {
		return nil
	}
	return &InventoryRepositoryRefRecord{
		ID:          model.ID,
		PortGroupID: model.PortGroupID,
		ProjectID:   model.ProjectID,
		Name:        model.Name,
		URL:         model.URL,
		Kind:        model.Kind,
		Notes:       model.Notes,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func utilCompactString(value string) string {
	return strings.TrimSpace(value)
}

func utilSearchPattern(value string) string {
	return "%" + strings.ToLower(strings.TrimSpace(value)) + "%"
}

func utilJoinSearchClauses(columns []string) string {
	clauses := make([]string, 0, len(columns))
	for _, column := range columns {
		clauses = append(clauses, "LOWER("+column+") LIKE ?")
	}
	return "(" + strings.Join(clauses, " OR ") + ")"
}

func utilJoinSearchArgs(pattern string, count int) []any {
	args := make([]any, 0, count)
	for range count {
		args = append(args, pattern)
	}
	return args
}
