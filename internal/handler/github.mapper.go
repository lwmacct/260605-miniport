package handler

import "github.com/lwmacct/260605-miniport/internal/service"

func ToGithubInstallationDTO(item service.GithubInstallation) GithubInstallationDTO {
	return GithubInstallationDTO{
		ID: item.ID, GithubInstallationID: item.GithubInstallationID, AccountID: item.AccountID,
		AccountLogin: item.AccountLogin, AccountType: item.AccountType, AvatarURL: item.AvatarURL,
		RepositorySelection: item.RepositorySelection, Status: item.Status, SuspendedAt: item.SuspendedAt,
		LastSyncedAt: item.LastSyncedAt, LastSyncError: item.LastSyncError,
	}
}

func ToGithubInstallationDTOs(items []service.GithubInstallation) []GithubInstallationDTO {
	out := make([]GithubInstallationDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ToGithubInstallationDTO(item))
	}
	return out
}

func ToGithubRepositoryDTO(item service.GithubRepository) GithubRepositoryDTO {
	return GithubRepositoryDTO{
		ID: item.ID, InstallationID: item.InstallationID, GithubRepositoryID: item.GithubRepositoryID,
		OwnerLogin: item.OwnerLogin, Name: item.Name, FullName: item.FullName, HTMLURL: item.HTMLURL,
		Description: item.Description, DefaultBranch: item.DefaultBranch, Visibility: item.Visibility,
		Private: item.Private, Fork: item.Fork, Archived: item.Archived, Disabled: item.Disabled,
		State: item.State, PushedAt: item.PushedAt, RemoteUpdatedAt: item.RemoteUpdatedAt, LastSeenAt: item.LastSeenAt,
	}
}

func ToGithubRepositoryDTOs(items []service.GithubRepository) []GithubRepositoryDTO {
	out := make([]GithubRepositoryDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ToGithubRepositoryDTO(item))
	}
	return out
}
