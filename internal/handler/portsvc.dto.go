package handler

type HostDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IP        string `json:"ip"`
	Spec      string `json:"spec"`
	Status    string `json:"status" enum:"active,stopped"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"createdAt" format:"date-time"`
	UpdatedAt string `json:"updatedAt" format:"date-time"`
}

type HostListDTO struct {
	Items []HostDTO `json:"items" nullable:"false"`
}

type HostBatchDTO struct {
	Items []HostDTO `json:"items" nullable:"false"`
}

type HostCreateDTO struct {
	Name   string `json:"name" minLength:"1"`
	IP     string `json:"ip,omitempty"`
	Spec   string `json:"spec,omitempty"`
	Status string `json:"status,omitempty" enum:"active,stopped"`
	Notes  string `json:"notes,omitempty"`
}

type HostUpdateDTO struct {
	ID     string `json:"id" minLength:"1"`
	Name   string `json:"name" minLength:"1"`
	IP     string `json:"ip"`
	Spec   string `json:"spec"`
	Status string `json:"status" enum:"active,stopped"`
	Notes  string `json:"notes"`
}

type HostCreateBatchDTO struct {
	Items []HostCreateDTO `json:"items" nullable:"false" minItems:"1"`
}

type HostUpdateBatchDTO struct {
	Items []HostUpdateDTO `json:"items" nullable:"false" minItems:"1"`
}

type HostCreateInputDTO struct{ Body HostCreateBatchDTO }
type HostUpdateInputDTO struct{ Body HostUpdateBatchDTO }
type HostDeleteInputDTO struct{ Body BatchIDsDTO }

type HostListInputDTO struct {
	Query  string `query:"q"`
	Status string `query:"status"`
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
	Status          string `json:"status" enum:"active,disabled"`
	Description     string `json:"description"`
	Metadata        string `json:"metadata"`
	LastSyncedAt    string `json:"lastSyncedAt,omitempty" format:"date-time"`
	Notes           string `json:"notes"`
	CreatedAt       string `json:"createdAt" format:"date-time"`
	UpdatedAt       string `json:"updatedAt" format:"date-time"`
}

type DependencyAssetListDTO struct {
	Items []DependencyAssetDTO `json:"items" nullable:"false"`
}

type DependencyAssetBatchDTO struct {
	Items []DependencyAssetDTO `json:"items" nullable:"false"`
}

type DependencyAssetCreateDTO struct {
	Name            string `json:"name" minLength:"1"`
	AssetKind       string `json:"assetKind,omitempty"`
	AssetType       string `json:"assetType,omitempty"`
	Provider        string `json:"provider,omitempty"`
	URL             string `json:"url,omitempty"`
	FullName        string `json:"fullName,omitempty"`
	ExternalID      string `json:"externalId,omitempty"`
	Visibility      string `json:"visibility,omitempty"`
	Controllability string `json:"controllability,omitempty"`
	Status          string `json:"status,omitempty" enum:"active,disabled"`
	Description     string `json:"description,omitempty"`
	Metadata        string `json:"metadata,omitempty"`
	Notes           string `json:"notes,omitempty"`
}

type DependencyAssetUpdateDTO struct {
	ID              string `json:"id" minLength:"1"`
	Name            string `json:"name" minLength:"1"`
	AssetKind       string `json:"assetKind"`
	AssetType       string `json:"assetType"`
	Provider        string `json:"provider"`
	URL             string `json:"url"`
	FullName        string `json:"fullName"`
	ExternalID      string `json:"externalId"`
	Visibility      string `json:"visibility"`
	Controllability string `json:"controllability"`
	Status          string `json:"status" enum:"active,disabled"`
	Description     string `json:"description"`
	Metadata        string `json:"metadata"`
	Notes           string `json:"notes"`
}

type DependencyAssetCreateBatchDTO struct {
	Items []DependencyAssetCreateDTO `json:"items" nullable:"false" minItems:"1"`
}

type DependencyAssetUpdateBatchDTO struct {
	Items []DependencyAssetUpdateDTO `json:"items" nullable:"false" minItems:"1"`
}

type DependencyAssetCreateInputDTO struct{ Body DependencyAssetCreateBatchDTO }
type DependencyAssetUpdateInputDTO struct{ Body DependencyAssetUpdateBatchDTO }
type DependencyAssetDeleteInputDTO struct{ Body BatchIDsDTO }

type DependencyAssetListInputDTO struct {
	Query     string `query:"q"`
	AssetKind string `query:"assetKind"`
	AssetType string `query:"assetType"`
	Provider  string `query:"provider"`
	Status    string `query:"status"`
}

type PortSlotDTO struct {
	ID            string `json:"id"`
	PortGroupID   string `json:"portGroupId"`
	Port          int    `json:"port" minimum:"1"`
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	Protocol      string `json:"protocol"`
	ContainerName string `json:"containerName"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`
	CreatedAt     string `json:"createdAt" format:"date-time"`
	UpdatedAt     string `json:"updatedAt" format:"date-time"`
}

type PortSlotInputDTO struct {
	ID            string `json:"id,omitempty"`
	Port          int    `json:"port,omitempty"`
	Name          string `json:"name,omitempty"`
	Kind          string `json:"kind,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	ContainerName string `json:"containerName,omitempty"`
	Status        string `json:"status,omitempty"`
	Notes         string `json:"notes,omitempty"`
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
	CreatedAt    string              `json:"createdAt" format:"date-time"`
	UpdatedAt    string              `json:"updatedAt" format:"date-time"`
}

type PortGroupAssetLinkInputDTO struct {
	ID           string `json:"id,omitempty"`
	PortSlotID   string `json:"portSlotId,omitempty"`
	AssetID      string `json:"assetId,omitempty"`
	RelationType string `json:"relationType,omitempty"`
	Required     bool   `json:"required,omitempty"`
	Notes        string `json:"notes,omitempty"`
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
	CreatedAt    string               `json:"createdAt" format:"date-time"`
	UpdatedAt    string               `json:"updatedAt" format:"date-time"`
}

type PortGroupRepositoryLinkInputDTO struct {
	ID           string `json:"id,omitempty"`
	PortSlotID   string `json:"portSlotId,omitempty"`
	RepositoryID string `json:"repositoryId,omitempty"`
	RelationType string `json:"relationType,omitempty"`
	Required     bool   `json:"required,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type PortGroupDTO struct {
	ID               string                       `json:"id"`
	HostID           string                       `json:"hostId"`
	Host             *HostDTO                     `json:"host,omitempty"`
	PortPrefix       int                          `json:"portPrefix"`
	EnvironmentName  string                       `json:"environmentName"`
	EnvironmentOwner string                       `json:"environmentOwner"`
	RuntimeMode      string                       `json:"runtimeMode" enum:"dind,host"`
	RuntimeName      string                       `json:"runtimeName"`
	ServiceIP        string                       `json:"serviceIp"`
	Status           string                       `json:"status"`
	Tags             string                       `json:"tags"`
	Notes            string                       `json:"notes"`
	CreatedAt        string                       `json:"createdAt" format:"date-time"`
	UpdatedAt        string                       `json:"updatedAt" format:"date-time"`
	Slots            []PortSlotDTO                `json:"slots" nullable:"false"`
	AssetLinks       []PortGroupAssetLinkDTO      `json:"assetLinks" nullable:"false"`
	RepositoryLinks  []PortGroupRepositoryLinkDTO `json:"repositoryLinks" nullable:"false"`
}

type PortGroupListDTO struct {
	Items []PortGroupDTO `json:"items" nullable:"false"`
}

type PortGroupBatchDTO struct {
	Items []PortGroupDTO `json:"items" nullable:"false"`
}

type PortGroupCreateDTO struct {
	HostID           string                            `json:"hostId,omitempty"`
	PortPrefix       int                               `json:"portPrefix,omitempty" minimum:"1000" maximum:"5999"`
	EnvironmentName  string                            `json:"environmentName,omitempty"`
	EnvironmentOwner string                            `json:"environmentOwner,omitempty"`
	RuntimeMode      string                            `json:"runtimeMode,omitempty" enum:"dind,host"`
	RuntimeName      string                            `json:"runtimeName,omitempty"`
	ServiceIP        string                            `json:"serviceIp,omitempty"`
	Status           string                            `json:"status,omitempty"`
	Tags             string                            `json:"tags,omitempty"`
	Notes            string                            `json:"notes,omitempty"`
	Slots            []PortSlotInputDTO                `json:"slots,omitempty" nullable:"false"`
	AssetLinks       []PortGroupAssetLinkInputDTO      `json:"assetLinks,omitempty" nullable:"false"`
	RepositoryLinks  []PortGroupRepositoryLinkInputDTO `json:"repositoryLinks,omitempty" nullable:"false"`
}

type PortGroupUpdateDTO struct {
	ID               string                            `json:"id" minLength:"1"`
	HostID           string                            `json:"hostId"`
	PortPrefix       int                               `json:"portPrefix" minimum:"1000" maximum:"5999"`
	EnvironmentName  string                            `json:"environmentName"`
	EnvironmentOwner string                            `json:"environmentOwner"`
	RuntimeMode      string                            `json:"runtimeMode" enum:"dind,host"`
	RuntimeName      string                            `json:"runtimeName"`
	ServiceIP        string                            `json:"serviceIp"`
	Status           string                            `json:"status"`
	Tags             string                            `json:"tags"`
	Notes            string                            `json:"notes"`
	Slots            []PortSlotInputDTO                `json:"slots" nullable:"false"`
	AssetLinks       []PortGroupAssetLinkInputDTO      `json:"assetLinks" nullable:"false"`
	RepositoryLinks  []PortGroupRepositoryLinkInputDTO `json:"repositoryLinks" nullable:"false"`
}

type PortGroupCreateBatchDTO struct {
	Items []PortGroupCreateDTO `json:"items" nullable:"false" minItems:"1"`
}

type PortGroupUpdateBatchDTO struct {
	Items []PortGroupUpdateDTO `json:"items" nullable:"false" minItems:"1"`
}

type PortGroupListInputDTO struct {
	Query  string `query:"q"`
	Sort   string `query:"sort" enum:"port,environment,status,updated_desc"`
	Status string `query:"status"`
}

type PortGroupCreateInputDTO struct{ Body PortGroupCreateBatchDTO }
type PortGroupUpdateInputDTO struct{ Body PortGroupUpdateBatchDTO }
type PortGroupDeleteInputDTO struct{ Body BatchIDsDTO }

type ServiceGroupPortGroupDTO struct {
	ID             string        `json:"id"`
	ServiceGroupID string        `json:"serviceGroupId"`
	PortGroupID    string        `json:"portGroupId"`
	PortGroup      *PortGroupDTO `json:"portGroup,omitempty"`
	Role           string        `json:"role"`
	Notes          string        `json:"notes"`
	CreatedAt      string        `json:"createdAt" format:"date-time"`
	UpdatedAt      string        `json:"updatedAt" format:"date-time"`
}

type ServiceGroupPortGroupInputDTO struct {
	ID          string `json:"id,omitempty"`
	PortGroupID string `json:"portGroupId,omitempty"`
	Role        string `json:"role,omitempty"`
	Notes       string `json:"notes,omitempty"`
}

type ServiceGroupDTO struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Kind        string                     `json:"kind"`
	Status      string                     `json:"status"`
	Description string                     `json:"description"`
	Notes       string                     `json:"notes"`
	CreatedAt   string                     `json:"createdAt" format:"date-time"`
	UpdatedAt   string                     `json:"updatedAt" format:"date-time"`
	PortGroups  []ServiceGroupPortGroupDTO `json:"portGroups" nullable:"false"`
}

type ServiceGroupListDTO struct {
	Items []ServiceGroupDTO `json:"items" nullable:"false"`
}

type ServiceGroupBatchDTO struct {
	Items []ServiceGroupDTO `json:"items" nullable:"false"`
}

type ServiceGroupCreateDTO struct {
	Name        string                          `json:"name" minLength:"1"`
	Kind        string                          `json:"kind,omitempty"`
	Status      string                          `json:"status,omitempty"`
	Description string                          `json:"description,omitempty"`
	Notes       string                          `json:"notes,omitempty"`
	PortGroups  []ServiceGroupPortGroupInputDTO `json:"portGroups,omitempty" nullable:"false"`
}

type ServiceGroupUpdateDTO struct {
	ID          string                          `json:"id" minLength:"1"`
	Name        string                          `json:"name" minLength:"1"`
	Kind        string                          `json:"kind"`
	Status      string                          `json:"status"`
	Description string                          `json:"description"`
	Notes       string                          `json:"notes"`
	PortGroups  []ServiceGroupPortGroupInputDTO `json:"portGroups" nullable:"false"`
}

type ServiceGroupCreateBatchDTO struct {
	Items []ServiceGroupCreateDTO `json:"items" nullable:"false" minItems:"1"`
}

type ServiceGroupUpdateBatchDTO struct {
	Items []ServiceGroupUpdateDTO `json:"items" nullable:"false" minItems:"1"`
}

type ServiceGroupListInputDTO struct {
	Query  string `query:"q"`
	Status string `query:"status"`
}

type ServiceGroupCreateInputDTO struct{ Body ServiceGroupCreateBatchDTO }
type ServiceGroupUpdateInputDTO struct{ Body ServiceGroupUpdateBatchDTO }
type ServiceGroupDeleteInputDTO struct{ Body BatchIDsDTO }

type CSVOutputDTO struct {
	ContentType        string `header:"Content-Type"`
	ContentDisposition string `header:"Content-Disposition"`
	Body               []byte `contentType:"text/csv"`
}
