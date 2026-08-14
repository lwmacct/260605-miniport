package repository

import (
	"errors"
	"time"
)

var ErrGithubWebhookDeliveryExists = errors.New("github webhook delivery already processed")

type GithubInstallationRecord struct {
	ID                   string
	GithubInstallationID int64
	AccountID            int64
	AccountLogin         string
	AccountType          string
	AvatarURL            string
	RepositorySelection  string
	Permissions          string
	Status               string
	SuspendedAt          time.Time
	LastSyncedAt         time.Time
	LastSyncError        string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type GithubRepositoryRecord struct {
	ID                 string
	InstallationID     string
	GithubRepositoryID int64
	NodeID             string
	OwnerLogin         string
	Name               string
	FullName           string
	HTMLURL            string
	Description        string
	DefaultBranch      string
	Visibility         string
	Private            bool
	Fork               bool
	Archived           bool
	Disabled           bool
	State              string
	PushedAt           time.Time
	RemoteUpdatedAt    time.Time
	LastSeenAt         time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type GithubConnectionStateRecord struct {
	ID        string
	StateHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type GithubWebhookDeliveryRecord struct {
	ID             string
	DeliveryID     string
	Event          string
	Action         string
	InstallationID int64
	Status         string
	Error          string
	ReceivedAt     time.Time
	ProcessedAt    time.Time
}

func utilGithubInstallationRecord(model *GithubInstallationsModel) *GithubInstallationRecord {
	if model == nil {
		return nil
	}
	return &GithubInstallationRecord{
		ID: model.ID, GithubInstallationID: model.GithubInstallationID, AccountID: model.AccountID,
		AccountLogin: model.AccountLogin, AccountType: model.AccountType, AvatarURL: model.AvatarURL,
		RepositorySelection: model.RepositorySelection, Permissions: model.Permissions, Status: model.Status,
		SuspendedAt: model.SuspendedAt, LastSyncedAt: model.LastSyncedAt, LastSyncError: model.LastSyncError,
		CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt,
	}
}

func utilGithubRepositoryRecord(model *GithubRepositoriesModel) *GithubRepositoryRecord {
	if model == nil {
		return nil
	}
	return &GithubRepositoryRecord{
		ID: model.ID, InstallationID: model.InstallationID, GithubRepositoryID: model.GithubRepositoryID,
		NodeID: model.NodeID, OwnerLogin: model.OwnerLogin, Name: model.Name, FullName: model.FullName,
		HTMLURL: model.HTMLURL, Description: model.Description, DefaultBranch: model.DefaultBranch,
		Visibility: model.Visibility, Private: model.Private, Fork: model.Fork, Archived: model.Archived,
		Disabled: model.Disabled, State: model.State, PushedAt: model.PushedAt,
		RemoteUpdatedAt: model.RemoteUpdatedAt, LastSeenAt: model.LastSeenAt,
		CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt,
	}
}

func utilGithubConnectionStateRecord(model *GithubConnectionStatesModel) *GithubConnectionStateRecord {
	if model == nil {
		return nil
	}
	return &GithubConnectionStateRecord{
		ID: model.ID, StateHash: model.StateHash, ExpiresAt: model.ExpiresAt, CreatedAt: model.CreatedAt,
	}
}
