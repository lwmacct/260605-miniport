package handler

type AdminUserListInputDTO struct {
	Session string `cookie:"web_session"`
}

type AdminUserDTO struct {
	ID          int64   `json:"id"`
	Username    string  `json:"username"`
	DisplayName string  `json:"displayName"`
	Status      string  `json:"status"`
	Admin       bool    `json:"admin"`
	DisabledAt  *string `json:"disabledAt,omitempty"`
}
