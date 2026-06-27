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
	user := createPortsvcTestUser(t, ctx, store, "member")
	now := time.Now().UTC()

	_, err := store.CreatePortsvcService(ctx, &PortsvcServiceRecord{
		UserID: user.ID, Name: "zeta", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcService(ctx, &PortsvcServiceRecord{
		UserID: user.ID, Name: "alpha", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)

	services, err := store.ListPortsvcServices(ctx, PortsvcServiceListFilter{UserID: user.ID})
	require.NoError(t, err)
	require.Len(t, services, 2)
	require.Equal(t, "alpha", services[0].Name)
	require.Equal(t, "zeta", services[1].Name)
}

func TestPortAllocationsAreUniquePerUser(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	first := createPortsvcTestUser(t, ctx, store, "first")
	second := createPortsvcTestUser(t, ctx, store, "second")
	now := time.Now().UTC()

	_, err := store.CreatePortsvcPortAllocation(ctx, &PortsvcPortAllocationRecord{
		UserID: first.ID, PortStart: 10000, PortEnd: 10009, Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcPortAllocation(ctx, &PortsvcPortAllocationRecord{
		UserID: second.ID, PortStart: 10000, PortEnd: 10009, Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcPortAllocation(ctx, &PortsvcPortAllocationRecord{
		UserID: first.ID, PortStart: 10000, PortEnd: 10009, Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.Error(t, err)
}

func TestServicePortAllocationIsUnique(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	user := createPortsvcTestUser(t, ctx, store, "member")
	now := time.Now().UTC()
	port, err := store.CreatePortsvcPortAllocation(ctx, &PortsvcPortAllocationRecord{
		UserID: user.ID, PortStart: 10000, PortEnd: 10009, Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcService(ctx, &PortsvcServiceRecord{
		UserID: user.ID, PortAllocationID: port.ID, Name: "first", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = store.CreatePortsvcService(ctx, &PortsvcServiceRecord{
		UserID: user.ID, PortAllocationID: port.ID, Name: "second", Status: "running", CreatedAt: now, UpdatedAt: now,
	})
	require.Error(t, err)
}

func createPortsvcTestUser(t *testing.T, ctx context.Context, store *Store, username string) *UserModel {
	t.Helper()
	user, err := store.CreateUser(ctx, username, username)
	require.NoError(t, err)
	return user
}

func newPortsvcTestDB(t *testing.T, ctx context.Context) *bun.DB {
	t.Helper()
	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	require.NoError(t, err)
	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
	models := append([]any{}, UserSchema()...)
	models = append(models, PortsvcSchema()...)
	require.NoError(t, dbschema.Apply(ctx, db, models, PortsvcIndexesSchema()))
	return db
}
