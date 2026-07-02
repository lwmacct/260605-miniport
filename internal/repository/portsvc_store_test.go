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
	ownerSubject := "018f2f9c-1111-7000-8000-000000000001"
	now := time.Now().UTC()

	_, err := store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		OwnerSubject: ownerSubject, PortPrefix: 1002, RuntimeMode: "dind", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		OwnerSubject: ownerSubject, PortPrefix: 1000, RuntimeMode: "dind", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)

	groups, err := store.ListPortsvcPortGroups(ctx, PortsvcPortGroupListFilter{OwnerSubject: ownerSubject})
	require.NoError(t, err)
	require.Len(t, groups, 2)
	require.Equal(t, 1000, groups[0].PortPrefix)
	require.Equal(t, 1002, groups[1].PortPrefix)
}

func TestPortGroupsAreUniquePerUser(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	first := "018f2f9c-1111-7000-8000-000000000001"
	second := "018f2f9c-1111-7000-8000-000000000002"
	now := time.Now().UTC()

	_, err := store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		OwnerSubject: first, PortPrefix: 1000, RuntimeMode: "dind", Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		OwnerSubject: second, PortPrefix: 1000, RuntimeMode: "dind", Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		OwnerSubject: first, PortPrefix: 1000, RuntimeMode: "dind", Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.Error(t, err)
}

func TestPortSlotsAreUniqueInsideGroup(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	ownerSubject := "018f2f9c-1111-7000-8000-000000000001"
	now := time.Now().UTC()
	group, err := store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		OwnerSubject: ownerSubject, PortPrefix: 1000, RuntimeMode: "dind", Status: "available", CreatedAt: now, UpdatedAt: now,
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

func newPortsvcTestDB(t *testing.T, ctx context.Context) *bun.DB {
	t.Helper()
	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	require.NoError(t, err)
	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
	require.NoError(t, dbschema.Apply(ctx, db, PortsvcSchema(), PortsvcIndexesSchema()))
	return db
}
