package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"

	"github.com/lwmacct/260605-miniport/internal/infra/dbschema"
)

func TestListPortsvcPortGroupsSortsByPort(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	now := time.Now().UTC()

	_, err := store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		PortPrefix: 1002, RuntimeMode: "dind", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		PortPrefix: 1000, RuntimeMode: "dind", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)

	groups, err := store.ListPortsvcPortGroups(ctx, PortsvcPortGroupListFilter{})
	require.NoError(t, err)
	require.Len(t, groups, 2)
	require.Equal(t, 1000, groups[0].PortPrefix)
	require.Equal(t, 1002, groups[1].PortPrefix)
}

func TestPortGroupsAreGloballyUnique(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	now := time.Now().UTC()

	_, err := store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		PortPrefix: 1000, RuntimeMode: "dind", Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		PortPrefix: 1000, RuntimeMode: "dind", Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.Error(t, err)
}

func TestPortSlotsAreUniqueInsideGroup(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	now := time.Now().UTC()
	group, err := store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		PortPrefix: 1000, RuntimeMode: "dind", Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcPortSlot(ctx, &PortsvcPortSlotRecord{
		PortGroupID: group.ID, Port: 10000, Name: "redis", Kind: "cache", Protocol: "redis", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcPortSlot(ctx, &PortsvcPortSlotRecord{
		PortGroupID: group.ID, Port: 10000, Name: "mysql", Kind: "db", Protocol: "mysql", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.Error(t, err)
}

func TestPortGroupRepositoryLinksLoadRepositoryAndPreserveSlotID(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	now := time.Now().UTC()

	installation, err := store.UpsertGithubInstallation(ctx, GithubInstallationRecord{
		GithubInstallationID: 42, AccountID: 7, AccountLogin: "acme", AccountType: "Organization",
		RepositorySelection: "all", Permissions: `{"metadata":"read"}`, Status: "active", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	require.NoError(t, store.ReplaceGithubRepositories(ctx, installation.ID, []GithubRepositoryRecord{{
		GithubRepositoryID: 99, OwnerLogin: "acme", Name: "api", FullName: "acme/api",
		HTMLURL: "https://github.com/acme/api", Visibility: "private", Private: true,
	}}, now))
	repositories, err := store.ListGithubRepositories(ctx, "", "")
	require.NoError(t, err)
	require.Len(t, repositories, 1)

	group, err := store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		PortPrefix: 1000, RuntimeMode: "dind", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	slotID := "018f2f9c-2222-7000-8000-000000000001"
	require.NoError(t, store.AddPortsvcPortSlots(ctx, []PortsvcPortSlotRecord{{
		ID: slotID, PortGroupID: group.ID, Port: 10000, Name: "api", Kind: "app", Protocol: "http",
		Status: "running", CreatedAt: now, UpdatedAt: now,
	}}))
	require.NoError(t, store.AddPortsvcPortGroupRepositoryLinks(ctx, []PortsvcPortGroupRepositoryLinkRecord{{
		PortGroupID: group.ID, PortSlotID: slotID, RepositoryID: repositories[0].ID,
		RelationType: "source", Required: true, CreatedAt: now, UpdatedAt: now,
	}}))

	children, err := store.ListPortsvcPortGroupChildrenByGroupIDs(ctx, []string{group.ID})
	require.NoError(t, err)
	require.Len(t, children, 1)
	require.Len(t, children[0].Slots, 1)
	require.Equal(t, slotID, children[0].Slots[0].ID)
	require.Len(t, children[0].RepositoryLinks, 1)
	require.Equal(t, "acme/api", children[0].RepositoryLinks[0].Repository.FullName)
}

func newPortsvcTestDB(t *testing.T, ctx context.Context) *bun.DB {
	t.Helper()
	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	require.NoError(t, err)
	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
	require.NoError(t, dbschema.Apply(ctx, db, GithubSchema(), GithubIndexesSchema()))
	require.NoError(t, dbschema.Apply(ctx, db, PortsvcSchema(), PortsvcIndexesSchema()))
	return db
}
