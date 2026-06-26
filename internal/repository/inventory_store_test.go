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

func TestListInventoryPortGroupsSortsByHostPort(t *testing.T) {
	ctx := t.Context()
	db := newInventoryTestDB(t, ctx)
	store := NewStore(db)
	now := time.Now().UTC()

	firstHost, err := store.CreateInventoryHost(ctx, &InventoryHostRecord{IP: "10.0.0.2", Name: "beta", CreatedAt: now, UpdatedAt: now})
	require.NoError(t, err)
	secondHost, err := store.CreateInventoryHost(ctx, &InventoryHostRecord{IP: "10.0.0.1", Name: "alpha", CreatedAt: now, UpdatedAt: now})
	require.NoError(t, err)

	_, err = store.CreateInventoryPortGroup(ctx, &InventoryPortGroupRecord{
		HostID: firstHost.ID, PortStart: 2000, PortEnd: 2009, ServiceName: "beta-service", Status: "used", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreateInventoryPortGroup(ctx, &InventoryPortGroupRecord{
		HostID: secondHost.ID, PortStart: 1000, PortEnd: 1009, ServiceName: "alpha-service", Status: "used", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)

	groups, err := store.ListInventoryPortGroups(ctx, InventoryPortGroupListFilter{Sort: "host_port"})
	require.NoError(t, err)
	require.Len(t, groups, 2)
	require.Equal(t, "10.0.0.1", groups[0].Host.IP)
	require.Equal(t, 1000, groups[0].PortStart)
	require.Equal(t, "10.0.0.2", groups[1].Host.IP)
	require.Equal(t, 2000, groups[1].PortStart)
}

func TestListInventoryPortGroupsSearchesHostFields(t *testing.T) {
	ctx := t.Context()
	db := newInventoryTestDB(t, ctx)
	store := NewStore(db)
	now := time.Now().UTC()

	host, err := store.CreateInventoryHost(ctx, &InventoryHostRecord{IP: "10.0.0.3", Name: "searchable-host", CreatedAt: now, UpdatedAt: now})
	require.NoError(t, err)
	_, err = store.CreateInventoryPortGroup(ctx, &InventoryPortGroupRecord{
		HostID: host.ID, PortStart: 3000, PortEnd: 3009, ServiceName: "worker", Status: "used", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)

	groups, err := store.ListInventoryPortGroups(ctx, InventoryPortGroupListFilter{Query: "searchable"})
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, host.ID, groups[0].HostID)
	require.Equal(t, "searchable-host", groups[0].Host.Name)
}

func newInventoryTestDB(t *testing.T, ctx context.Context) *bun.DB {
	t.Helper()
	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	require.NoError(t, err)
	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
	require.NoError(t, dbschema.Apply(ctx, db, InventorySchema(), InventoryIndexesSchema()))
	return db
}
