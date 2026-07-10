package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"

	"github.com/lwmacct/260605-miniport/internal/infra/dbschema"
	"github.com/lwmacct/260605-miniport/internal/repository"
)

type fakeGithubAPI struct {
	installation githubRemoteInstallation
	repositories []githubRemoteRepository
	err          error
}

func (f fakeGithubAPI) Installation(context.Context, int64) (githubRemoteInstallation, error) {
	return f.installation, f.err
}

func (f fakeGithubAPI) Repositories(context.Context, int64) ([]githubRemoteRepository, error) {
	return f.repositories, f.err
}

func TestGithubCompleteConnectionSyncsPrivateRepositories(t *testing.T) {
	ctx := t.Context()
	store := newGithubTestStore(t, ctx)
	now := time.Date(2026, 7, 10, 8, 0, 0, 0, time.UTC)
	ownerSubject := "018f2f9c-1111-7000-8000-000000000001"
	state := "one-time-state"
	require.NoError(t, store.AddGithubConnectionState(ctx, utilHashGithubState(state), ownerSubject, now.Add(time.Minute), now))

	githubService := &GithubService{
		store: store,
		cfg:   GithubConfig{Enabled: true},
		now:   func() time.Time { return now },
		api: fakeGithubAPI{
			installation: githubRemoteInstallation{
				ID: 42, Account: githubRemoteOwner{ID: 7, Login: "acme", Type: "Organization"},
				RepositorySelection: "all", Permissions: map[string]string{"metadata": "read"},
			},
			repositories: []githubRemoteRepository{{
				ID: 99, Owner: githubRemoteOwner{Login: "acme"}, Name: "private-api", FullName: "acme/private-api",
				HTMLURL: "https://github.com/acme/private-api", Private: true, Visibility: "private", DefaultBranch: "main",
			}},
		},
	}

	require.NoError(t, githubService.CompleteConnection(ctx, ownerSubject, state, 42))
	installations, err := githubService.ListInstallations(ctx, ownerSubject)
	require.NoError(t, err)
	require.Len(t, installations, 1)
	require.Equal(t, "all", installations[0].RepositorySelection)

	repositories, err := githubService.ListRepositories(ctx, ownerSubject, "private", "")
	require.NoError(t, err)
	require.Len(t, repositories, 1)
	require.Equal(t, "acme/private-api", repositories[0].FullName)
	require.True(t, repositories[0].Private)

	err = githubService.CompleteConnection(ctx, ownerSubject, state, 42)
	require.ErrorIs(t, err, ErrGithubInvalidState)
}

func TestVerifyGithubWebhookSignature(t *testing.T) {
	body := []byte(`{"action":"created"}`)
	mac := hmac.New(sha256.New, []byte("secret"))
	_, _ = mac.Write(body)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	require.True(t, validateGithubWebhookSignature("secret", signature, body))
	require.False(t, validateGithubWebhookSignature("wrong", signature, body))
	require.False(t, validateGithubWebhookSignature("secret", "sha1=bad", body))
}

func newGithubTestStore(t *testing.T, ctx context.Context) *repository.Store {
	t.Helper()
	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	require.NoError(t, err)
	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, dbschema.Apply(ctx, db, repository.GithubSchema(), repository.GithubIndexesSchema()))
	return repository.NewStore(db)
}
