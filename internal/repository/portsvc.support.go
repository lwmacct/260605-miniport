package repository

import (
	"strings"
	"time"
)

type PortsvcServiceListFilter struct {
	UserID      int64
	Admin       bool
	Query       string
	Sort        string
	Status      string
	ProjectName string
}

type PortsvcPortAllocationListFilter struct {
	UserID int64
	Admin  bool
	Sort   string
	Status string
}

type PortsvcServiceRecord struct {
	ID               int64
	UserID           int64
	Username         string
	PortAllocationID int64
	PortAllocation   *PortsvcPortAllocationRecord
	Name             string
	ProjectName      string
	DindIP           string
	DindContainer    string
	Status           string
	Owner            string
	Tags             string
	Notes            string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type PortsvcPortAllocationRecord struct {
	ID        int64
	UserID    int64
	Username  string
	PortStart int
	PortEnd   int
	Status    string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PortsvcDependencyRecord struct {
	ID        int64
	UserID    int64
	Name      string
	Type      string
	URL       string
	Version   string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PortsvcRepositoryRecord struct {
	ID        int64
	UserID    int64
	Name      string
	URL       string
	Kind      string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PortsvcServiceRepositoryRecord struct {
	ID           int64
	ServiceID    int64
	RepositoryID int64
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PortsvcServiceDependencyRecord struct {
	ID           int64
	ServiceID    int64
	DependencyID int64
	Role         string
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PortsvcServiceChildrenRecord struct {
	ServiceID    int64
	Repositories []PortsvcRepositoryRecord
	Dependencies []PortsvcDependencyRecord
}

func utilPortsvcServiceRecordFromModel(model *ServicesModel) *PortsvcServiceRecord {
	if model == nil {
		return nil
	}
	out := &PortsvcServiceRecord{
		ID:               model.ID,
		UserID:           model.UserID,
		PortAllocationID: model.PortAllocationID,
		Name:             model.Name,
		ProjectName:      model.ProjectName,
		DindIP:           model.DindIP,
		DindContainer:    model.DindContainer,
		Status:           model.Status,
		Owner:            model.Owner,
		Tags:             model.Tags,
		Notes:            model.Notes,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
	if model.User != nil {
		out.Username = model.User.Username
	}
	out.PortAllocation = utilPortsvcPortAllocationRecordFromModel(model.PortAllocation)
	return out
}

func utilPortsvcPortAllocationRecordFromModel(model *PortAllocationsModel) *PortsvcPortAllocationRecord {
	if model == nil {
		return nil
	}
	out := &PortsvcPortAllocationRecord{
		ID:        model.ID,
		UserID:    model.UserID,
		PortStart: model.PortStart,
		PortEnd:   model.PortEnd,
		Status:    model.Status,
		Notes:     model.Notes,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
	if model.User != nil {
		out.Username = model.User.Username
	}
	return out
}

func utilPortsvcDependencyRecordFromModel(model *DependenciesModel) *PortsvcDependencyRecord {
	if model == nil {
		return nil
	}
	return &PortsvcDependencyRecord{
		ID:        model.ID,
		UserID:    model.UserID,
		Name:      model.Name,
		Type:      model.Type,
		URL:       model.URL,
		Version:   model.Version,
		Notes:     model.Notes,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func utilPortsvcRepositoryRecordFromModel(model *RepositoriesModel) *PortsvcRepositoryRecord {
	if model == nil {
		return nil
	}
	return &PortsvcRepositoryRecord{
		ID:        model.ID,
		UserID:    model.UserID,
		Name:      model.Name,
		URL:       model.URL,
		Kind:      model.Kind,
		Notes:     model.Notes,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
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
