package handler

type PortSlotPayloadDTO struct {
	Port     int    `json:"port,omitempty"`
	Name     string `json:"name,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Purpose  string `json:"purpose,omitempty"`
	Status   string `json:"status,omitempty"`
	Notes    string `json:"notes,omitempty"`
}

type ProjectPayloadDTO struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Notes       string `json:"notes,omitempty"`
}

type ComponentPayloadDTO struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	URL     string `json:"url,omitempty"`
	Version string `json:"version,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type RepositoryPayloadDTO struct {
	ProjectID int64  `json:"projectId,omitempty"`
	Name      string `json:"name,omitempty"`
	URL       string `json:"url,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

type AllocationPayloadDTO struct {
	UserID        int64                  `json:"userId,omitempty"`
	PortStart     int                    `json:"portStart,omitempty"`
	PortEnd       int                    `json:"portEnd,omitempty"`
	Name          string                 `json:"name,omitempty"`
	DindIP        string                 `json:"dindIp,omitempty"`
	DindContainer string                 `json:"dindContainer,omitempty"`
	Status        string                 `json:"status,omitempty"`
	Owner         string                 `json:"owner,omitempty"`
	Tags          string                 `json:"tags,omitempty"`
	Notes         string                 `json:"notes,omitempty"`
	Slots         []PortSlotPayloadDTO   `json:"slots,omitempty"`
	Projects      []ProjectPayloadDTO    `json:"projects,omitempty"`
	Components    []ComponentPayloadDTO  `json:"components,omitempty"`
	Repositories  []RepositoryPayloadDTO `json:"repositories,omitempty"`
}

type AllocationBatchUpdateDTO struct {
	IDs    []int64 `json:"ids,omitempty"`
	Owner  *string `json:"owner,omitempty"`
	Status *string `json:"status,omitempty"`
	Tags   *string `json:"tags,omitempty"`
}

type AllocationBatchDeleteDTO struct {
	IDs []int64 `json:"ids,omitempty"`
}

type PortSlotDTO struct {
	ID           int64  `json:"id"`
	AllocationID int64  `json:"allocationId"`
	Port         int    `json:"port"`
	Name         string `json:"name"`
	Protocol     string `json:"protocol"`
	Purpose      string `json:"purpose"`
	Status       string `json:"status"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type ProjectDTO struct {
	ID           int64  `json:"id"`
	AllocationID int64  `json:"allocationId"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type ComponentDTO struct {
	ID           int64  `json:"id"`
	AllocationID int64  `json:"allocationId"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	URL          string `json:"url"`
	Version      string `json:"version"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type RepositoryDTO struct {
	ID           int64  `json:"id"`
	AllocationID int64  `json:"allocationId"`
	ProjectID    int64  `json:"projectId"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Kind         string `json:"kind"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type AllocationDTO struct {
	ID            int64           `json:"id"`
	UserID        int64           `json:"userId"`
	Username      string          `json:"username"`
	PortStart     int             `json:"portStart"`
	PortEnd       int             `json:"portEnd"`
	Name          string          `json:"name"`
	DindIP        string          `json:"dindIp"`
	DindContainer string          `json:"dindContainer"`
	Status        string          `json:"status"`
	Owner         string          `json:"owner"`
	Tags          string          `json:"tags"`
	Notes         string          `json:"notes"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
	Slots         []PortSlotDTO   `json:"slots"`
	Projects      []ProjectDTO    `json:"projects"`
	Components    []ComponentDTO  `json:"components"`
	Repositories  []RepositoryDTO `json:"repositories"`
}

type AllocationListInputDTO struct {
	Session     string `cookie:"web_session"`
	UserID      int64  `query:"userId" example:"1"`
	Query       string `query:"q" example:"postgres"`
	Sort        string `query:"sort" example:"port"`
	Status      string `query:"status" example:"running"`
	ProjectName string `query:"projectName" example:"miniport"`
	DindIP      string `query:"dindIp" example:"172.20.0.12"`
}

type AllocationInputDTO struct {
	Session string `cookie:"web_session"`
	ID      int64  `path:"id" example:"1"`
}

type AllocationBodyInputDTO struct {
	Session string `cookie:"web_session"`
	Body    AllocationPayloadDTO
}

type AllocationUpdateInputDTO struct {
	Session string `cookie:"web_session"`
	ID      int64  `path:"id" example:"1"`
	Body    AllocationPayloadDTO
}

type AllocationBatchUpdateInputDTO struct {
	Session string `cookie:"web_session"`
	Body    AllocationBatchUpdateDTO
}

type AllocationBatchDeleteInputDTO struct {
	Session string `cookie:"web_session"`
	Body    AllocationBatchDeleteDTO
}

type CSVOutputDTO struct {
	ContentType        string `header:"Content-Type"`
	ContentDisposition string `header:"Content-Disposition"`
	Body               []byte
}
