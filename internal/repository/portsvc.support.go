package repository

import (
	"strings"
	"time"
)

type PortsvcServiceListFilter struct {
	OwnerSubject string
	Admin        bool
	Query        string
	Sort         string
	Status       string
	ProjectName  string
}

type PortsvcPortAllocationListFilter struct {
	OwnerSubject string
	Admin        bool
	Sort         string
	Status       string
}

type PortsvcServiceRecord struct {
	ID               string
	OwnerSubject     string
	OwnerName        string
	PortAllocationID string
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
	ID           string
	OwnerSubject string
	OwnerName    string
	PortStart    int
	PortEnd      int
	Status       string
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PortsvcDependencyRecord struct {
	ID           string
	OwnerSubject string
	Name         string
	Type         string
	URL          string
	Version      string
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PortsvcRepositoryRecord struct {
	ID           string
	OwnerSubject string
	Name         string
	URL          string
	Kind         string
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PortsvcServiceRepositoryRecord struct {
	ID           string
	ServiceID    string
	RepositoryID string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PortsvcServiceDependencyRecord struct {
	ID           string
	ServiceID    string
	DependencyID string
	Role         string
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PortsvcServiceChildrenRecord struct {
	ServiceID    string
	Repositories []PortsvcRepositoryRecord
	Dependencies []PortsvcDependencyRecord
}

func utilPortsvcServiceRecordFromModel(model *ServicesModel) *PortsvcServiceRecord {
	if model == nil {
		return nil
	}
	out := &PortsvcServiceRecord{
		ID:               model.ID,
		OwnerSubject:     model.OwnerSubject,
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
	out.PortAllocation = utilPortsvcPortAllocationRecordFromModel(model.PortAllocation)
	return out
}

func utilPortsvcPortAllocationRecordFromModel(model *PortAllocationsModel) *PortsvcPortAllocationRecord {
	if model == nil {
		return nil
	}
	out := &PortsvcPortAllocationRecord{
		ID:           model.ID,
		OwnerSubject: model.OwnerSubject,
		PortStart:    model.PortStart,
		PortEnd:      model.PortEnd,
		Status:       model.Status,
		Notes:        model.Notes,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
	return out
}

func utilPortsvcDependencyRecordFromModel(model *DependenciesModel) *PortsvcDependencyRecord {
	if model == nil {
		return nil
	}
	return &PortsvcDependencyRecord{
		ID:           model.ID,
		OwnerSubject: model.OwnerSubject,
		Name:         model.Name,
		Type:         model.Type,
		URL:          model.URL,
		Version:      model.Version,
		Notes:        model.Notes,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func utilPortsvcRepositoryRecordFromModel(model *RepositoriesModel) *PortsvcRepositoryRecord {
	if model == nil {
		return nil
	}
	return &PortsvcRepositoryRecord{
		ID:           model.ID,
		OwnerSubject: model.OwnerSubject,
		Name:         model.Name,
		URL:          model.URL,
		Kind:         model.Kind,
		Notes:        model.Notes,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
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
