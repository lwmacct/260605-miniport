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

type PortsvcDependencyAssetListFilter struct {
	OwnerSubject string
	Admin        bool
	Query        string
	AssetKind    string
	AssetType    string
	Provider     string
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

type PortsvcDependencyAssetRecord struct {
	ID              string
	OwnerSubject    string
	OwnerName       string
	Name            string
	AssetKind       string
	AssetType       string
	Provider        string
	URL             string
	FullName        string
	ExternalID      string
	Visibility      string
	Controllability string
	Status          string
	Description     string
	Metadata        string
	LastSyncedAt    time.Time
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type PortsvcPortGroupAssetLinkRecord struct {
	ID           string
	PortGroupID  string
	PortSlotID   string
	AssetID      string
	Asset        *PortsvcDependencyAssetRecord
	RelationType string
	Required     bool
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PortsvcPortGroupChildrenRecord struct {
	PortGroupID string
	Slots       []PortsvcPortSlotRecord
	AssetLinks  []PortsvcPortGroupAssetLinkRecord
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

func utilPortsvcDependencyAssetRecordFromModel(model *DependenciesModel) *PortsvcDependencyAssetRecord {
	if model == nil {
		return nil
	}
	return &PortsvcDependencyAssetRecord{
		ID:              model.ID,
		OwnerSubject:    model.OwnerSubject,
		Name:            model.Name,
		AssetKind:       model.AssetKind,
		AssetType:       model.AssetType,
		Provider:        model.Provider,
		URL:             model.URL,
		FullName:        model.FullName,
		ExternalID:      model.ExternalID,
		Visibility:      model.Visibility,
		Controllability: model.Controllability,
		Status:          model.Status,
		Description:     model.Description,
		Metadata:        model.Metadata,
		LastSyncedAt:    model.LastSyncedAt,
		Notes:           model.Notes,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func utilPortsvcPortGroupAssetLinkRecordFromModel(model *PortGroupAssetLinksModel) *PortsvcPortGroupAssetLinkRecord {
	if model == nil {
		return nil
	}
	out := &PortsvcPortGroupAssetLinkRecord{
		ID:           model.ID,
		PortGroupID:  model.PortGroupID,
		PortSlotID:   model.PortSlotID,
		AssetID:      model.AssetID,
		RelationType: model.RelationType,
		Required:     model.Required,
		Notes:        model.Notes,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
	out.Asset = utilPortsvcDependencyAssetRecordFromModel(model.Asset)
	return out
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
