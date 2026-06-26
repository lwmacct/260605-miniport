package repository

import (
	"strings"
	"time"
)

type InventoryHostListFilter struct {
	Environment string
	Query       string
	Sort        string
}

type InventoryPortGroupListFilter struct {
	HostID int64
	Query  string
	Sort   string
	Status string
}

type InventoryHostRecord struct {
	ID          int64
	IP          string
	Name        string
	Network     string
	Environment string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InventoryPortGroupRecord struct {
	ID            int64
	HostID        int64
	Host          *InventoryHostRecord
	PortStart     int
	PortEnd       int
	ServiceName   string
	ContainerName string
	DindHost      string
	Status        string
	Owner         string
	Tags          string
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
	Name        string
	URL         string
	Kind        string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InventoryPortGroupChildrenRecord struct {
	Slots        []InventoryPortSlotRecord
	Components   []InventoryComponentRecord
	Repositories []InventoryRepositoryRefRecord
}

func utilInventoryHostRecordFromModel(model *InventoryHostModel) *InventoryHostRecord {
	if model == nil {
		return nil
	}
	return &InventoryHostRecord{
		ID:          model.ID,
		IP:          model.IP,
		Name:        model.Name,
		Network:     model.Network,
		Environment: model.Environment,
		Notes:       model.Notes,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func utilInventoryPortGroupRecordFromModel(model *InventoryPortGroupModel) *InventoryPortGroupRecord {
	if model == nil {
		return nil
	}
	out := &InventoryPortGroupRecord{
		ID:            model.ID,
		HostID:        model.HostID,
		PortStart:     model.PortStart,
		PortEnd:       model.PortEnd,
		ServiceName:   model.ServiceName,
		ContainerName: model.ContainerName,
		DindHost:      model.DindHost,
		Status:        model.Status,
		Owner:         model.Owner,
		Tags:          model.Tags,
		Notes:         model.Notes,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
	out.Host = utilInventoryHostRecordFromModel(model.InventoryHostModel)
	return out
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
