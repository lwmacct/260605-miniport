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

type DependencyAssetPayloadDTO struct {
	Name            string `json:"name,omitempty"`
	AssetKind       string `json:"assetKind,omitempty"`
	AssetType       string `json:"assetType,omitempty"`
	Provider        string `json:"provider,omitempty"`
	URL             string `json:"url,omitempty"`
	FullName        string `json:"fullName,omitempty"`
	ExternalID      string `json:"externalId,omitempty"`
	Visibility      string `json:"visibility,omitempty"`
	Controllability string `json:"controllability,omitempty"`
	Status          string `json:"status,omitempty"`
	Description     string `json:"description,omitempty"`
	Metadata        string `json:"metadata,omitempty"`
	Notes           string `json:"notes,omitempty"`
}

type PortGroupAssetLinkPayloadDTO struct {
	ID           string `json:"id,omitempty"`
	PortSlotID   string `json:"portSlotId,omitempty"`
	AssetID      string `json:"assetId,omitempty"`
	RelationType string `json:"relationType,omitempty"`
	Required     bool   `json:"required,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type PortGroupRepositoryLinkPayloadDTO struct {
	ID           string `json:"id,omitempty"`
	PortSlotID   string `json:"portSlotId,omitempty"`
	RepositoryID string `json:"repositoryId,omitempty"`
	RelationType string `json:"relationType,omitempty"`
	Required     bool   `json:"required,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type PortGroupPayloadDTO struct {
	HostID           string                              `json:"hostId,omitempty"`
	PortPrefix       int                                 `json:"portPrefix,omitempty"`
	EnvironmentName  string                              `json:"environmentName,omitempty"`
	EnvironmentOwner string                              `json:"environmentOwner,omitempty"`
	RuntimeMode      string                              `json:"runtimeMode,omitempty"`
	RuntimeName      string                              `json:"runtimeName,omitempty"`
	ServiceIP        string                              `json:"serviceIp,omitempty"`
	Status           string                              `json:"status,omitempty"`
	Tags             string                              `json:"tags,omitempty"`
	Notes            string                              `json:"notes,omitempty"`
	Slots            []PortSlotPayloadDTO                `json:"slots,omitempty"`
	AssetLinks       []PortGroupAssetLinkPayloadDTO      `json:"assetLinks,omitempty"`
	RepositoryLinks  []PortGroupRepositoryLinkPayloadDTO `json:"repositoryLinks,omitempty"`
}

type ServiceGroupPortGroupPayloadDTO struct {
	ID          string `json:"id,omitempty"`
	PortGroupID string `json:"portGroupId,omitempty"`
	Role        string `json:"role,omitempty"`
	Notes       string `json:"notes,omitempty"`
}

type ServiceGroupPayloadDTO struct {
	Name        string                            `json:"name,omitempty"`
	Kind        string                            `json:"kind,omitempty"`
	Status      string                            `json:"status,omitempty"`
	Description string                            `json:"description,omitempty"`
	Notes       string                            `json:"notes,omitempty"`
	PortGroups  []ServiceGroupPortGroupPayloadDTO `json:"portGroups,omitempty"`
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
	ID               string                       `json:"id"`
	HostID           string                       `json:"hostId"`
	Host             *HostDTO                     `json:"host,omitempty"`
	PortPrefix       int                          `json:"portPrefix"`
	EnvironmentName  string                       `json:"environmentName"`
	EnvironmentOwner string                       `json:"environmentOwner"`
	RuntimeMode      string                       `json:"runtimeMode"`
	RuntimeName      string                       `json:"runtimeName"`
	ServiceIP        string                       `json:"serviceIp"`
	Status           string                       `json:"status"`
	Tags             string                       `json:"tags"`
	Notes            string                       `json:"notes"`
	CreatedAt        string                       `json:"createdAt"`
	UpdatedAt        string                       `json:"updatedAt"`
	Slots            []PortSlotDTO                `json:"slots"`
	AssetLinks       []PortGroupAssetLinkDTO      `json:"assetLinks"`
	RepositoryLinks  []PortGroupRepositoryLinkDTO `json:"repositoryLinks"`
}

type ServiceGroupPortGroupDTO struct {
	ID             string        `json:"id"`
	ServiceGroupID string        `json:"serviceGroupId"`
	PortGroupID    string        `json:"portGroupId"`
	PortGroup      *PortGroupDTO `json:"portGroup,omitempty"`
	Role           string        `json:"role"`
	Notes          string        `json:"notes"`
	CreatedAt      string        `json:"createdAt"`
	UpdatedAt      string        `json:"updatedAt"`
}

type ServiceGroupDTO struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Kind        string                     `json:"kind"`
	Status      string                     `json:"status"`
	Description string                     `json:"description"`
	Notes       string                     `json:"notes"`
	CreatedAt   string                     `json:"createdAt"`
	UpdatedAt   string                     `json:"updatedAt"`
	PortGroups  []ServiceGroupPortGroupDTO `json:"portGroups"`
}

type DependencyAssetDTO struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	AssetKind       string `json:"assetKind"`
	AssetType       string `json:"assetType"`
	Provider        string `json:"provider"`
	URL             string `json:"url"`
	FullName        string `json:"fullName"`
	ExternalID      string `json:"externalId"`
	Visibility      string `json:"visibility"`
	Controllability string `json:"controllability"`
	Status          string `json:"status"`
	Description     string `json:"description"`
	Metadata        string `json:"metadata"`
	LastSyncedAt    string `json:"lastSyncedAt"`
	Notes           string `json:"notes"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

type PortGroupAssetLinkDTO struct {
	ID           string              `json:"id"`
	PortGroupID  string              `json:"portGroupId"`
	PortSlotID   string              `json:"portSlotId"`
	AssetID      string              `json:"assetId"`
	Asset        *DependencyAssetDTO `json:"asset,omitempty"`
	RelationType string              `json:"relationType"`
	Required     bool                `json:"required"`
	Notes        string              `json:"notes"`
	CreatedAt    string              `json:"createdAt"`
	UpdatedAt    string              `json:"updatedAt"`
}

type PortGroupRepositoryLinkDTO struct {
	ID           string               `json:"id"`
	PortGroupID  string               `json:"portGroupId"`
	PortSlotID   string               `json:"portSlotId"`
	RepositoryID string               `json:"repositoryId"`
	Repository   *GithubRepositoryDTO `json:"repository,omitempty"`
	RelationType string               `json:"relationType"`
	Required     bool                 `json:"required"`
	Notes        string               `json:"notes"`
	CreatedAt    string               `json:"createdAt"`
	UpdatedAt    string               `json:"updatedAt"`
}

type HostListInputDTO struct {
	Query  string `query:"q" example:"4h4g"`
	Status string `query:"status" example:"active"`
}

type HostInputDTO struct {
	ID string `path:"id" example:"018f2f9c-1111-7000-8000-000000000001"`
}

type HostBodyInputDTO struct {
	Body HostPayloadDTO
}

type HostUpdateInputDTO struct {
	ID   string `path:"id" example:"018f2f9c-1111-7000-8000-000000000001"`
	Body HostPayloadDTO
}

type PortGroupListInputDTO struct {
	Query  string `query:"q" example:"miniport"`
	Sort   string `query:"sort" example:"port"`
	Status string `query:"status" example:"running"`
}

type ServiceGroupListInputDTO struct {
	Query  string `query:"q" example:"etcd"`
	Status string `query:"status" example:"active"`
}

type DependencyAssetListInputDTO struct {
	Query     string `query:"q" example:"github"`
	AssetKind string `query:"assetKind" example:"repository"`
	AssetType string `query:"assetType" example:"owned"`
	Provider  string `query:"provider" example:"github"`
	Status    string `query:"status" example:"active"`
}

type DependencyAssetInputDTO struct {
	ID string `path:"id" example:"018f2f9c-1111-7000-8000-000000000004"`
}

type DependencyAssetBodyInputDTO struct {
	Body DependencyAssetPayloadDTO
}

type DependencyAssetUpdateInputDTO struct {
	ID   string `path:"id" example:"018f2f9c-1111-7000-8000-000000000004"`
	Body DependencyAssetPayloadDTO
}

type PortGroupInputDTO struct {
	ID string `path:"id" example:"018f2f9c-1111-7000-8000-000000000002"`
}

type PortGroupBodyInputDTO struct {
	Body PortGroupPayloadDTO
}

type PortGroupUpdateInputDTO struct {
	ID   string `path:"id" example:"018f2f9c-1111-7000-8000-000000000002"`
	Body PortGroupPayloadDTO
}

type ServiceGroupInputDTO struct {
	ID string `path:"id" example:"018f2f9c-1111-7000-8000-000000000005"`
}

type ServiceGroupBodyInputDTO struct {
	Body ServiceGroupPayloadDTO
}

type ServiceGroupUpdateInputDTO struct {
	ID   string `path:"id" example:"018f2f9c-1111-7000-8000-000000000005"`
	Body ServiceGroupPayloadDTO
}

type PortSlotBodyInputDTO struct {
	ID   string `path:"id" example:"018f2f9c-1111-7000-8000-000000000002"`
	Body PortSlotPayloadDTO
}

type PortSlotUpdateInputDTO struct {
	ID   string `path:"id" example:"018f2f9c-1111-7000-8000-000000000003"`
	Body PortSlotPayloadDTO
}

type CSVOutputDTO struct {
	ContentType        string `header:"Content-Type"`
	ContentDisposition string `header:"Content-Disposition"`
	Body               []byte
}
