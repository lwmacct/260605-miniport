package inventory

type HostPayload struct {
	IP          string `json:"ip" example:"172.22.11.12"`
	Name        string `json:"name,omitempty" example:"node-12"`
	Network     string `json:"network,omitempty" example:"172.22.11.0/24"`
	Environment string `json:"environment,omitempty" example:"dev"`
	Notes       string `json:"notes,omitempty"`
}

type PortSlotPayload struct {
	Port     int    `json:"port" example:"11120"`
	Name     string `json:"name,omitempty" example:"nginx"`
	Protocol string `json:"protocol,omitempty" example:"tcp"`
	Purpose  string `json:"purpose,omitempty" example:"reverse proxy"`
	Status   string `json:"status,omitempty" example:"used"`
	Notes    string `json:"notes,omitempty"`
}

type ComponentPayload struct {
	Name    string `json:"name" example:"redis"`
	Type    string `json:"type,omitempty" example:"opensource"`
	URL     string `json:"url,omitempty" example:"https://redis.io"`
	Version string `json:"version,omitempty" example:"7"`
	Notes   string `json:"notes,omitempty"`
}

type RepositoryPayload struct {
	Name  string `json:"name" example:"service-api"`
	URL   string `json:"url" example:"https://git.example.com/team/service-api"`
	Kind  string `json:"kind,omitempty" example:"backend"`
	Notes string `json:"notes,omitempty"`
}

type PortGroupPayload struct {
	HostID        int64               `json:"hostId" example:"1"`
	PortStart     int                 `json:"portStart" example:"11120"`
	PortEnd       int                 `json:"portEnd" example:"11129"`
	ServiceName   string              `json:"serviceName" example:"order-service"`
	ContainerName string              `json:"containerName,omitempty" example:"order-service-dind"`
	DindHost      string              `json:"dindHost,omitempty" example:"dind-01"`
	Status        string              `json:"status,omitempty" example:"running"`
	Owner         string              `json:"owner,omitempty" example:"platform"`
	Tags          string              `json:"tags,omitempty" example:"api,redis"`
	Notes         string              `json:"notes,omitempty"`
	Slots         []PortSlotPayload   `json:"slots,omitempty"`
	Components    []ComponentPayload  `json:"components,omitempty"`
	Repositories  []RepositoryPayload `json:"repositories,omitempty"`
}

type HostListParams struct {
	Environment string
	Query       string
	Sort        string
}

type PortGroupView struct {
	PortGroup

	Slots        []PortSlot   `json:"slots"`
	Components   []Component  `json:"components"`
	Repositories []Repository `json:"repositories"`
}

type PortGroupListParams struct {
	HostID int64
	Query  string
	Sort   string
	Status string
}

type PortGroupBatchUpdateInput struct {
	IDs    []int64 `json:"ids"`
	Owner  *string `json:"owner,omitempty"`
	Status *string `json:"status,omitempty"`
	Tags   *string `json:"tags,omitempty"`
}

type PortGroupBatchDeleteInput struct {
	IDs []int64 `json:"ids"`
}
