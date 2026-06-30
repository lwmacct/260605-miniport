package handler

type RepositoryPayloadDTO struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Kind  string `json:"kind,omitempty"`
	Role  string `json:"role,omitempty"`
	Notes string `json:"notes,omitempty"`
}

type DependencyPayloadDTO struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	URL     string `json:"url,omitempty"`
	Version string `json:"version,omitempty"`
	Role    string `json:"role,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type ServicePayloadDTO struct {
	OwnerSubject     string                 `json:"ownerSubject,omitempty"`
	PortAllocationID string                 `json:"portAllocationId,omitempty"`
	Name             string                 `json:"name,omitempty"`
	ProjectName      string                 `json:"projectName,omitempty"`
	DindIP           string                 `json:"dindIp,omitempty"`
	DindContainer    string                 `json:"dindContainer,omitempty"`
	Status           string                 `json:"status,omitempty"`
	Owner            string                 `json:"owner,omitempty"`
	Tags             string                 `json:"tags,omitempty"`
	Notes            string                 `json:"notes,omitempty"`
	Repositories     []RepositoryPayloadDTO `json:"repositories,omitempty"`
	Dependencies     []DependencyPayloadDTO `json:"dependencies,omitempty"`
}

type PortAllocationPayloadDTO struct {
	OwnerSubject string `json:"ownerSubject,omitempty"`
	PortStart    int    `json:"portStart,omitempty"`
	PortEnd      int    `json:"portEnd,omitempty"`
	Status       string `json:"status,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type PortAllocationDTO struct {
	ID           string `json:"id"`
	OwnerSubject string `json:"ownerSubject"`
	OwnerName    string `json:"ownerName"`
	PortStart    int    `json:"portStart"`
	PortEnd      int    `json:"portEnd"`
	Status       string `json:"status"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type RepositoryDTO struct {
	ID           string `json:"id"`
	OwnerSubject string `json:"ownerSubject"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Kind         string `json:"kind"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type DependencyDTO struct {
	ID           string `json:"id"`
	OwnerSubject string `json:"ownerSubject"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	URL          string `json:"url"`
	Version      string `json:"version"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type ServiceDTO struct {
	ID               string             `json:"id"`
	OwnerSubject     string             `json:"ownerSubject"`
	OwnerName        string             `json:"ownerName"`
	PortAllocationID string             `json:"portAllocationId"`
	PortAllocation   *PortAllocationDTO `json:"portAllocation,omitempty"`
	Name             string             `json:"name"`
	ProjectName      string             `json:"projectName"`
	DindIP           string             `json:"dindIp"`
	DindContainer    string             `json:"dindContainer"`
	Status           string             `json:"status"`
	Owner            string             `json:"owner"`
	Tags             string             `json:"tags"`
	Notes            string             `json:"notes"`
	CreatedAt        string             `json:"createdAt"`
	UpdatedAt        string             `json:"updatedAt"`
	Repositories     []RepositoryDTO    `json:"repositories"`
	Dependencies     []DependencyDTO    `json:"dependencies"`
}

type ServiceListInputDTO struct {
	Session      string `cookie:"web_session"`
	OwnerSubject string `query:"ownerSubject" example:"018f2f9c-1111-7000-8000-000000000001"`
	Query        string `query:"q" example:"postgres"`
	Sort         string `query:"sort" example:"name"`
	Status       string `query:"status" example:"running"`
	ProjectName  string `query:"projectName" example:"miniport"`
}

type ServiceInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id" example:"018f2f9c-1111-7000-8000-000000000002"`
}

type ServiceBodyInputDTO struct {
	Session string `cookie:"web_session"`
	Body    ServicePayloadDTO
}

type ServiceUpdateInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id" example:"018f2f9c-1111-7000-8000-000000000002"`
	Body    ServicePayloadDTO
}

type ServiceBatchDeleteInputDTO struct {
	Session string `cookie:"web_session"`
	Body    struct {
		IDs []string `json:"ids,omitempty"`
	}
}

type PortAllocationListInputDTO struct {
	Session      string `cookie:"web_session"`
	OwnerSubject string `query:"ownerSubject" example:"018f2f9c-1111-7000-8000-000000000001"`
	Sort         string `query:"sort" example:"port"`
	Status       string `query:"status" example:"available"`
}

type PortAllocationInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id" example:"018f2f9c-1111-7000-8000-000000000003"`
}

type PortAllocationBodyInputDTO struct {
	Session string `cookie:"web_session"`
	Body    PortAllocationPayloadDTO
}

type PortAllocationUpdateInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id" example:"018f2f9c-1111-7000-8000-000000000003"`
	Body    PortAllocationPayloadDTO
}

type CSVOutputDTO struct {
	ContentType        string `header:"Content-Type"`
	ContentDisposition string `header:"Content-Disposition"`
	Body               []byte
}
