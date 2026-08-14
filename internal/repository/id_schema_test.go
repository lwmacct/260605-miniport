package repository

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAllSchemaModelsUseUUIDPrimaryKeys(t *testing.T) {
	models := append([]any{}, GithubSchema()...)
	models = append(models, PortsvcSchema()...)

	for _, model := range models {
		modelType := reflect.TypeOf(model)
		field, ok := modelType.Elem().FieldByName("ID")
		require.True(t, ok, "%s must define an ID field", modelType.Elem().Name())
		require.Equal(t, reflect.TypeFor[string](), field.Type, "%s ID must be string-backed UUID", modelType.Elem().Name())

		tags := strings.Split(field.Tag.Get("bun"), ",")
		require.Contains(t, tags, "pk", "%s ID must be the primary key", modelType.Elem().Name())
		require.Contains(t, tags, "type:uuid", "%s ID must use the UUID database type", modelType.Elem().Name())
	}
}

func TestAddPortsvcPortSlotsRegeneratesNonUUID7IDs(t *testing.T) {
	ctx := t.Context()
	db := newPortsvcTestDB(t, ctx)
	store := NewStore(db)
	now := time.Now().UTC()

	group, err := store.CreatePortsvcPortGroup(ctx, &PortsvcPortGroupRecord{
		PortPrefix: 1000, RuntimeMode: "dind", Status: "available", CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)

	err = store.AddPortsvcPortSlots(ctx, []PortsvcPortSlotRecord{{
		ID: "legacy-id", PortGroupID: group.ID, Port: 10000, Name: "api", Kind: "app", Protocol: "http",
		Status: "running", CreatedAt: now, UpdatedAt: now,
	}})
	require.NoError(t, err)
	children, err := store.ListPortsvcPortGroupChildrenByGroupIDs(ctx, []string{group.ID})
	require.NoError(t, err)
	require.Len(t, children, 1)
	require.Len(t, children[0].Slots, 1)
	require.True(t, IsUUID7(children[0].Slots[0].ID))
}
