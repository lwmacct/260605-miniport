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

func TestListInventoryPortGroupsSortsByPort(t *testing.T) {
	ctx := t.Context()
	db := newInventoryTestDB(t, ctx)
	store := NewStore(db)
	user := createInventoryTestUser(t, ctx, store, "member")
	now := time.Now().UTC()

	_, err := store.CreateInventoryPortGroup(ctx, &InventoryPortGroupRecord{
		UserID: user.ID, PortStart: 10020, PortEnd: 10029, Name: "beta", Status: "used", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreateInventoryPortGroup(ctx, &InventoryPortGroupRecord{
		UserID: user.ID, PortStart: 10000, PortEnd: 10009, Name: "alpha", Status: "used", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)

	groups, err := store.ListInventoryPortGroups(ctx, InventoryPortGroupListFilter{UserID: user.ID})
	require.NoError(t, err)
	require.Len(t, groups, 2)
	require.Equal(t, 10000, groups[0].PortStart)
	require.Equal(t, 10020, groups[1].PortStart)
}

func TestListInventoryPortGroupsSearchesProjectFields(t *testing.T) {
	ctx := t.Context()
	db := newInventoryTestDB(t, ctx)
	store := NewStore(db)
	user := createInventoryTestUser(t, ctx, store, "member")
	now := time.Now().UTC()

	group, err := store.CreateInventoryPortGroup(ctx, &InventoryPortGroupRecord{
		UserID: user.ID, PortStart: 10000, PortEnd: 10009, Name: "worker", Status: "used", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	require.NoError(t, store.AddInventoryProjects(ctx, []InventoryProjectRecord{{
		PortGroupID: group.ID, Name: "searchable-project", CreatedAt: now, UpdatedAt: now,
	}}))

	groups, err := store.ListInventoryPortGroups(ctx, InventoryPortGroupListFilter{UserID: user.ID, Query: "searchable"})
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, group.ID, groups[0].ID)
}

func TestInventoryPortGroupsAreUniquePerUser(t *testing.T) {
	ctx := t.Context()
	db := newInventoryTestDB(t, ctx)
	store := NewStore(db)
	first := createInventoryTestUser(t, ctx, store, "first")
	second := createInventoryTestUser(t, ctx, store, "second")
	now := time.Now().UTC()

	_, err := store.CreateInventoryPortGroup(ctx, &InventoryPortGroupRecord{
		UserID: first.ID, PortStart: 10000, PortEnd: 10009, Name: "first", Status: "used", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreateInventoryPortGroup(ctx, &InventoryPortGroupRecord{
		UserID: second.ID, PortStart: 10000, PortEnd: 10009, Name: "second", Status: "used", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreateInventoryPortGroup(ctx, &InventoryPortGroupRecord{
		UserID: first.ID, PortStart: 10000, PortEnd: 10009, Name: "duplicate", Status: "used", CreatedAt: now, UpdatedAt: now,
	})
	require.Error(t, err)
}

func createInventoryTestUser(t *testing.T, ctx context.Context, store *Store, username string) *UserModel {
	t.Helper()
	user, err := store.CreateUser(ctx, username, username)
	require.NoError(t, err)
	return user
}

func newInventoryTestDB(t *testing.T, ctx context.Context) *bun.DB {
	t.Helper()
	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	require.NoError(t, err)
	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
	models := append([]any{}, UserSchema()...)
	models = append(models, InventorySchema()...)
	require.NoError(t, dbschema.Apply(ctx, db, models, InventoryIndexesSchema()))
	return db
}
