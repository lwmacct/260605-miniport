package handler

type HostPayloadDTO struct {
	Name   string `json:"name,omitempty"`
	IP     string `json:"ip,omitempty"`
	Spec   string `json:"spec,omitempty"`
	Status string `json:"status,omitempty"`
	Notes  string `json:"notes,omitempty"`
}

type PortSlotPayloadDTO struct {
	ID            string `json:"id,omitempty"`
	Port          int    `json:"port,omitempty"`
	Name          string `json:"name,omitempty"`
	Kind          string `json:"kind,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	ContainerName string `json:"containerName,omitempty"`
	Status        string `json:"status,omitempty"`
	Notes         string `json:"notes,omitempty"`
}

type RepositoryPayloadDTO struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Kind  string `json:"kind,omitempty"`
	Notes string `json:"notes,omitempty"`
}

type DependencyPayloadDTO struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	URL     string `json:"url,omitempty"`
	Version string `json:"version,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type PortGroupPayloadDTO struct {
	OwnerSubject string                 `json:"ownerSubject,omitempty"`
	HostID       string                 `json:"hostId,omitempty"`
	PortStart    int                    `json:"portStart,omitempty"`
	PortEnd      int                    `json:"portEnd,omitempty"`
	ProjectName  string                 `json:"projectName,omitempty"`
	ProjectOwner string                 `json:"projectOwner,omitempty"`
	RuntimeMode  string                 `json:"runtimeMode,omitempty"`
	RuntimeName  string                 `json:"runtimeName,omitempty"`
	ServiceIP    string                 `json:"serviceIp,omitempty"`
	Status       string                 `json:"status,omitempty"`
	Tags         string                 `json:"tags,omitempty"`
	Notes        string                 `json:"notes,omitempty"`
	Slots        []PortSlotPayloadDTO   `json:"slots,omitempty"`
	Repositories []RepositoryPayloadDTO `json:"repositories,omitempty"`
	Dependencies []DependencyPayloadDTO `json:"dependencies,omitempty"`
}

type HostDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IP        string `json:"ip"`
	Spec      string `json:"spec"`
	Status    string `json:"status"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type PortSlotDTO struct {
	ID            string `json:"id"`
	PortGroupID   string `json:"portGroupId"`
	Port          int    `json:"port"`
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	Protocol      string `json:"protocol"`
	ContainerName string `json:"containerName"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type PortGroupDTO struct {
	ID           string          `json:"id"`
	OwnerSubject string          `json:"ownerSubject"`
	OwnerName    string          `json:"ownerName"`
	HostID       string          `json:"hostId"`
	Host         *HostDTO        `json:"host,omitempty"`
	PortStart    int             `json:"portStart"`
	PortEnd      int             `json:"portEnd"`
	ProjectName  string          `json:"projectName"`
	ProjectOwner string          `json:"projectOwner"`
	RuntimeMode  string          `json:"runtimeMode"`
	RuntimeName  string          `json:"runtimeName"`
	ServiceIP    string          `json:"serviceIp"`
	Status       string          `json:"status"`
	Tags         string          `json:"tags"`
	Notes        string          `json:"notes"`
	CreatedAt    string          `json:"createdAt"`
	UpdatedAt    string          `json:"updatedAt"`
	Slots        []PortSlotDTO   `json:"slots"`
	Repositories []RepositoryDTO `json:"repositories"`
	Dependencies []DependencyDTO `json:"dependencies"`
}

type RepositoryDTO struct {
	ID           string `json:"id"`
	OwnerSubject string `json:"ownerSubject"`
	PortGroupID  string `json:"portGroupId"`
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
	PortGroupID  string `json:"portGroupId"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	URL          string `json:"url"`
	Version      string `json:"version"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type HostListInputDTO struct {
	Session string `cookie:"web_session"`
	Query   string `query:"q" example:"4h4g"`
	Status  string `query:"status" example:"active"`
}

type HostInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id" example:"018f2f9c-1111-7000-8000-000000000001"`
}

type HostBodyInputDTO struct {
	Session string `cookie:"web_session"`
	Body    HostPayloadDTO
}

type HostUpdateInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id" example:"018f2f9c-1111-7000-8000-000000000001"`
	Body    HostPayloadDTO
}

type PortGroupListInputDTO struct {
	Session      string `cookie:"web_session"`
	OwnerSubject string `query:"ownerSubject" example:"018f2f9c-1111-7000-8000-000000000001"`
	Query        string `query:"q" example:"miniport"`
	Sort         string `query:"sort" example:"port"`
	Status       string `query:"status" example:"running"`
}

type PortGroupInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id" example:"018f2f9c-1111-7000-8000-000000000002"`
}

type PortGroupBodyInputDTO struct {
	Session string `cookie:"web_session"`
	Body    PortGroupPayloadDTO
}

type PortGroupUpdateInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id" example:"018f2f9c-1111-7000-8000-000000000002"`
	Body    PortGroupPayloadDTO
}

type PortSlotBodyInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id" example:"018f2f9c-1111-7000-8000-000000000002"`
	Body    PortSlotPayloadDTO
}

type PortSlotUpdateInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id" example:"018f2f9c-1111-7000-8000-000000000003"`
	Body    PortSlotPayloadDTO
}

type CSVOutputDTO struct {
	ContentType        string `header:"Content-Type"`
	ContentDisposition string `header:"Content-Disposition"`
	Body               []byte
}
