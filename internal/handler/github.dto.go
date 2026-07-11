package handler

import "time"

type GithubStatusDTO struct {
	Enabled    bool   `json:"enabled"`
	AppSlug    string `json:"appSlug"`
	InstallURL string `json:"installUrl"`
}

type GithubConnectionStartDTO struct {
	URL string `json:"url"`
}

type GithubInstallationDTO struct {
	ID                   string    `json:"id"`
	GithubInstallationID int64     `json:"githubInstallationId"`
	AccountID            int64     `json:"accountId"`
	AccountLogin         string    `json:"accountLogin"`
	AccountType          string    `json:"accountType"`
	AvatarURL            string    `json:"avatarUrl"`
	RepositorySelection  string    `json:"repositorySelection"`
	Status               string    `json:"status"`
	SuspendedAt          time.Time `json:"suspendedAt,omitzero"`
	LastSyncedAt         time.Time `json:"lastSyncedAt,omitzero"`
	LastSyncError        string    `json:"lastSyncError,omitempty"`
}

type GithubRepositoryDTO struct {
	ID                 string    `json:"id"`
	InstallationID     string    `json:"installationId"`
	GithubRepositoryID int64     `json:"githubRepositoryId"`
	OwnerLogin         string    `json:"ownerLogin"`
	Name               string    `json:"name"`
	FullName           string    `json:"fullName"`
	HTMLURL            string    `json:"htmlUrl"`
	Description        string    `json:"description"`
	DefaultBranch      string    `json:"defaultBranch"`
	Visibility         string    `json:"visibility"`
	Private            bool      `json:"private"`
	Fork               bool      `json:"fork"`
	Archived           bool      `json:"archived"`
	Disabled           bool      `json:"disabled"`
	State              string    `json:"state"`
	PushedAt           time.Time `json:"pushedAt,omitzero"`
	RemoteUpdatedAt    time.Time `json:"remoteUpdatedAt,omitzero"`
	LastSeenAt         time.Time `json:"lastSeenAt"`
}

type GithubSessionInputDTO struct {
	Session string `cookie:"web_session"`
}

type GithubSetupInputDTO struct {
	InstallationID int64  `query:"installation_id" required:"true"`
	SetupAction    string `query:"setup_action"`
	State          string `query:"state" required:"true"`
}

type GithubInstallationInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id"`
}

type GithubRepositoryListInputDTO struct {
	Session string `cookie:"web_session"`
	Query   string `query:"q"`
	State   string `query:"state"`
}

type GithubWebhookInputDTO struct {
	Delivery  string `header:"X-GitHub-Delivery" required:"true"`
	Event     string `header:"X-GitHub-Event" required:"true"`
	Signature string `header:"X-Hub-Signature-256" required:"true"`
	RawBody   []byte `contentType:"application/json"`
}

type RedirectDTO struct {
	Status   int
	Location string `header:"Location"`
}
