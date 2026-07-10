package service

import (
	"context"

	githubinfra "github.com/lwmacct/260605-miniport/internal/infra/github"
)

type GithubProvider struct {
	client *githubinfra.Client
}

func NewGithubProvider(cfg GithubConfig) (*GithubProvider, error) {
	client, err := githubinfra.NewClient(githubinfra.Config{
		AppID: cfg.AppID, APIURL: cfg.APIURL, PrivateKeyFile: cfg.PrivateKeyFile,
	}, nil)
	if err != nil {
		return nil, err
	}
	return &GithubProvider{client: client}, nil
}

func (p *GithubProvider) Installation(ctx context.Context, installationID int64) (githubRemoteInstallation, error) {
	item, err := p.client.Installation(ctx, installationID)
	if err != nil {
		return githubRemoteInstallation{}, err
	}
	return githubRemoteInstallation{
		ID: item.ID,
		Account: githubRemoteOwner{
			ID: item.Account.ID, Login: item.Account.Login, Type: item.Account.Type, AvatarURL: item.Account.AvatarURL,
		},
		RepositorySelection: item.RepositorySelection, Permissions: item.Permissions, SuspendedAt: item.SuspendedAt,
	}, nil
}

func (p *GithubProvider) Repositories(ctx context.Context, installationID int64) ([]githubRemoteRepository, error) {
	items, err := p.client.Repositories(ctx, installationID)
	if err != nil {
		return nil, err
	}
	out := make([]githubRemoteRepository, 0, len(items))
	for _, item := range items {
		out = append(out, githubRemoteRepository{
			ID: item.ID, NodeID: item.NodeID,
			Owner: githubRemoteOwner{
				ID: item.Owner.ID, Login: item.Owner.Login, Type: item.Owner.Type, AvatarURL: item.Owner.AvatarURL,
			},
			Name: item.Name, FullName: item.FullName, HTMLURL: item.HTMLURL, Description: item.Description,
			DefaultBranch: item.DefaultBranch, Visibility: item.Visibility, Private: item.Private, Fork: item.Fork,
			Archived: item.Archived, Disabled: item.Disabled, PushedAt: item.PushedAt, UpdatedAt: item.UpdatedAt,
		})
	}
	return out, nil
}
