package repository

import (
	"strings"
	"time"
)

type PortsvcHostListFilter struct {
	Query  string
	Status string
}

type PortsvcPortGroupListFilter struct {
	OwnerSubject string
	Admin        bool
	Query        string
	Sort         string
	Status       string
}

type PortsvcHostRecord struct {
	ID        string
	Name      string
	IP        string
	Spec      string
	Status    string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PortsvcPortGroupRecord struct {
	ID           string
	OwnerSubject string
	OwnerName    string
	HostID       string
	Host         *PortsvcHostRecord
	PortStart    int
	PortEnd      int
	ProjectName  string
	ProjectOwner string
	RuntimeMode  string
	RuntimeName  string
	ServiceIP    string
	Status       string
	Tags         string
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PortsvcPortSlotRecord struct {
	ID            string
	PortGroupID   string
	Port          int
	Name          string
	Kind          string
	Protocol      string
	ContainerName string
	Status        string
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PortsvcDependencyRecord struct {
	ID           string
	OwnerSubject string
	PortGroupID  string
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
	PortGroupID  string
	Name         string
	URL          string
	Kind         string
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PortsvcPortGroupChildrenRecord struct {
	PortGroupID  string
	Slots        []PortsvcPortSlotRecord
	Repositories []PortsvcRepositoryRecord
	Dependencies []PortsvcDependencyRecord
}

func utilPortsvcHostRecordFromModel(model *HostsModel) *PortsvcHostRecord {
	if model == nil {
		return nil
	}
	return &PortsvcHostRecord{
		ID:        model.ID,
		Name:      model.Name,
		IP:        model.IP,
		Spec:      model.Spec,
		Status:    model.Status,
		Notes:     model.Notes,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func utilPortsvcPortGroupRecordFromModel(model *PortAllocationsModel) *PortsvcPortGroupRecord {
	if model == nil {
		return nil
	}
	out := &PortsvcPortGroupRecord{
		ID:           model.ID,
		OwnerSubject: model.OwnerSubject,
		HostID:       model.HostID,
		PortStart:    model.PortStart,
		PortEnd:      model.PortEnd,
		ProjectName:  model.ProjectName,
		ProjectOwner: model.ProjectOwner,
		RuntimeMode:  model.RuntimeMode,
		RuntimeName:  model.RuntimeName,
		ServiceIP:    model.ServiceIP,
		Status:       model.Status,
		Tags:         model.Tags,
		Notes:        model.Notes,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
	out.Host = utilPortsvcHostRecordFromModel(model.Host)
	return out
}

func utilPortsvcPortSlotRecordFromModel(model *ServicesModel) *PortsvcPortSlotRecord {
	if model == nil {
		return nil
	}
	return &PortsvcPortSlotRecord{
		ID:            model.ID,
		PortGroupID:   model.PortGroupID,
		Port:          model.Port,
		Name:          model.Name,
		Kind:          model.Kind,
		Protocol:      model.Protocol,
		ContainerName: model.ContainerName,
		Status:        model.Status,
		Notes:         model.Notes,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}

func utilPortsvcDependencyRecordFromModel(model *DependenciesModel) *PortsvcDependencyRecord {
	if model == nil {
		return nil
	}
	return &PortsvcDependencyRecord{
		ID:           model.ID,
		OwnerSubject: model.OwnerSubject,
		PortGroupID:  model.PortGroupID,
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
		PortGroupID:  model.PortGroupID,
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
