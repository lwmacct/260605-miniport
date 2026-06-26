package handler

import "github.com/lwmacct/260605-miniport/internal/service"

type HostDTO struct {
	ID          int64  `json:"id"`
	IP          string `json:"ip"`
	Name        string `json:"name"`
	Network     string `json:"network"`
	Environment string `json:"environment"`
	Notes       string `json:"notes"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type PortSlotDTO struct {
	ID          int64  `json:"id"`
	PortGroupID int64  `json:"portGroupId"`
	Port        int    `json:"port"`
	Name        string `json:"name"`
	Protocol    string `json:"protocol"`
	Purpose     string `json:"purpose"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type ComponentDTO struct {
	ID          int64  `json:"id"`
	PortGroupID int64  `json:"portGroupId"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Version     string `json:"version"`
	Notes       string `json:"notes"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type RepositoryDTO struct {
	ID          int64  `json:"id"`
	PortGroupID int64  `json:"portGroupId"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Kind        string `json:"kind"`
	Notes       string `json:"notes"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type PortGroupDTO struct {
	ID            int64           `json:"id"`
	HostID        int64           `json:"hostId"`
	Host          *HostDTO        `json:"host,omitempty"`
	PortStart     int             `json:"portStart"`
	PortEnd       int             `json:"portEnd"`
	ServiceName   string          `json:"serviceName"`
	ContainerName string          `json:"containerName"`
	DindHost      string          `json:"dindHost"`
	Status        string          `json:"status"`
	Owner         string          `json:"owner"`
	Tags          string          `json:"tags"`
	Notes         string          `json:"notes"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
	Slots         []PortSlotDTO   `json:"slots"`
	Components    []ComponentDTO  `json:"components"`
	Repositories  []RepositoryDTO `json:"repositories"`
}

type HostListInputDTO struct {
	Session     string `cookie:"web_session"`
	Environment string `query:"environment" example:"dev"`
	Query       string `query:"q" example:"172.22.11"`
	Sort        string `query:"sort" example:"ip"`
}

type HostInputDTO struct {
	Session string `cookie:"web_session"`
	ID      int64  `path:"id" example:"1"`
}

type HostBodyInputDTO struct {
	Session string `cookie:"web_session"`
	Body    service.HostPayload
}

type HostUpdateInputDTO struct {
	Session string `cookie:"web_session"`
	ID      int64  `path:"id" example:"1"`
	Body    service.HostPayload
}

type PortGroupListInputDTO struct {
	Session string `cookie:"web_session"`
	HostID  int64  `query:"hostId" example:"1"`
	Query   string `query:"q" example:"order-service"`
	Sort    string `query:"sort" example:"host_port"`
	Status  string `query:"status" example:"running"`
}

type PortGroupInputDTO struct {
	Session string `cookie:"web_session"`
	ID      int64  `path:"id" example:"1"`
}

type PortGroupBodyInputDTO struct {
	Session string `cookie:"web_session"`
	Body    service.PortGroupPayload
}

type PortGroupUpdateInputDTO struct {
	Session string `cookie:"web_session"`
	ID      int64  `path:"id" example:"1"`
	Body    service.PortGroupPayload
}

type PortGroupBatchUpdateInputDTO struct {
	Session string `cookie:"web_session"`
	Body    service.PortGroupBatchUpdateInput
}

type PortGroupBatchDeleteInputDTO struct {
	Session string `cookie:"web_session"`
	Body    service.PortGroupBatchDeleteInput
}

type CSVOutputDTO struct {
	ContentType        string `header:"Content-Type"`
	ContentDisposition string `header:"Content-Disposition"`
	Body               []byte
}
