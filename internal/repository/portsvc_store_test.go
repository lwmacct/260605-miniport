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

func TestListPortsvcServicesSortsByName(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	ownerSubject := "018f2f9c-1111-7000-8000-000000000001"
	now := time.Now().UTC()

	_, err := store.CreatePortsvcService(ctx, &PortsvcServiceRecord{
		OwnerSubject: ownerSubject, Name: "zeta", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcService(ctx, &PortsvcServiceRecord{
		OwnerSubject: ownerSubject, Name: "alpha", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)

	services, err := store.ListPortsvcServices(ctx, PortsvcServiceListFilter{OwnerSubject: ownerSubject})
	require.NoError(t, err)
	require.Len(t, services, 2)
	require.Equal(t, "alpha", services[0].Name)
	require.Equal(t, "zeta", services[1].Name)
}

func TestPortAllocationsAreUniquePerUser(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	first := "018f2f9c-1111-7000-8000-000000000001"
	second := "018f2f9c-1111-7000-8000-000000000002"
	now := time.Now().UTC()

	_, err := store.CreatePortsvcPortAllocation(ctx, &PortsvcPortAllocationRecord{
		OwnerSubject: first, PortStart: 10000, PortEnd: 10009, Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcPortAllocation(ctx, &PortsvcPortAllocationRecord{
		OwnerSubject: second, PortStart: 10000, PortEnd: 10009, Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcPortAllocation(ctx, &PortsvcPortAllocationRecord{
		OwnerSubject: first, PortStart: 10000, PortEnd: 10009, Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.Error(t, err)
}

func TestServicePortAllocationIsUnique(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	ownerSubject := "018f2f9c-1111-7000-8000-000000000001"
	now := time.Now().UTC()
	port, err := store.CreatePortsvcPortAllocation(ctx, &PortsvcPortAllocationRecord{
		OwnerSubject: ownerSubject, PortStart: 10000, PortEnd: 10009, Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcService(ctx, &PortsvcServiceRecord{
		OwnerSubject: ownerSubject, PortAllocationID: port.ID, Name: "first", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcService(ctx, &PortsvcServiceRecord{
		OwnerSubject: ownerSubject, PortAllocationID: port.ID, Name: "second", Status: "running", CreatedAt: now, UpdatedAt: now,
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
