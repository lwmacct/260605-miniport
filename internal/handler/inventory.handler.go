package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260605-miniport/internal/service"
)

type inventoryHandler struct {
	auth     authHandler
	services Services
}

func RegisterInventory(api huma.API, config Config, services Services) {
	handler := inventoryHandler{
		auth:     authHandler{config: config, services: services},
		services: services,
	}
	huma.Register(api, huma.Operation{OperationID: "list-inventory-allocations", Method: http.MethodGet, Path: "/inventory/allocations", Summary: "List port allocations", Tags: []string{"inventory"}}, handler.listAllocations)
	huma.Register(api, huma.Operation{OperationID: "get-inventory-allocation", Method: http.MethodGet, Path: "/inventory/allocations/{id}", Summary: "Get port allocation", Tags: []string{"inventory"}}, handler.getAllocation)
	huma.Register(api, huma.Operation{OperationID: "create-inventory-allocation", Method: http.MethodPost, Path: "/inventory/allocations", Summary: "Create port allocation", Tags: []string{"inventory"}}, handler.createAllocation)
	huma.Register(api, huma.Operation{OperationID: "update-inventory-allocation", Method: http.MethodPut, Path: "/inventory/allocations/{id}", Summary: "Update port allocation", Tags: []string{"inventory"}}, handler.updateAllocation)
	huma.Register(api, huma.Operation{OperationID: "delete-inventory-allocation", Method: http.MethodDelete, Path: "/inventory/allocations/{id}", Summary: "Delete port allocation", Tags: []string{"inventory"}}, handler.deleteAllocation)
	huma.Register(api, huma.Operation{OperationID: "batch-update-inventory-allocations", Method: http.MethodPost, Path: "/inventory/allocations/batch-update", Summary: "Batch update port allocations", Tags: []string{"inventory"}}, handler.batchUpdateAllocations)
	huma.Register(api, huma.Operation{OperationID: "batch-delete-inventory-allocations", Method: http.MethodPost, Path: "/inventory/allocations/batch-delete", Summary: "Batch delete port allocations", Tags: []string{"inventory"}}, handler.batchDeleteAllocations)
	huma.Register(api, huma.Operation{OperationID: "export-inventory-allocations", Method: http.MethodGet, Path: "/inventory/exports/allocations.csv", Summary: "Export port allocations", Tags: []string{"inventory"}}, handler.exportAllocations)
}

func (h inventoryHandler) listAllocations(ctx context.Context, input *AllocationListInputDTO) (*BodyDTO[[]AllocationDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	groups, err := h.services.Inventory.ListPortGroups(ctx, service.PortGroupListParams{Actor: actor, UserID: input.UserID, Query: input.Query, Sort: input.Sort, Status: input.Status, ProjectName: input.ProjectName, DindIP: input.DindIP})
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[[]AllocationDTO]{Body: ToAllocationDTOs(groups)}, nil
}

func (h inventoryHandler) getAllocation(ctx context.Context, input *AllocationInputDTO) (*BodyDTO[AllocationDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	group, err := h.services.Inventory.GetPortGroup(ctx, actor, input.ID)
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[AllocationDTO]{Body: ToAllocationDTO(*group)}, nil
}

func (h inventoryHandler) createAllocation(ctx context.Context, input *AllocationBodyInputDTO) (*BodyDTO[AllocationDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	group, err := h.services.Inventory.CreatePortGroup(ctx, actor, ToPortGroupPayload(input.Body))
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[AllocationDTO]{Body: ToAllocationDTO(*group)}, nil
}

func (h inventoryHandler) updateAllocation(ctx context.Context, input *AllocationUpdateInputDTO) (*BodyDTO[AllocationDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	group, err := h.services.Inventory.UpdatePortGroup(ctx, actor, input.ID, ToPortGroupPayload(input.Body))
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[AllocationDTO]{Body: ToAllocationDTO(*group)}, nil
}

func (h inventoryHandler) deleteAllocation(ctx context.Context, input *AllocationInputDTO) (*BodyDTO[DeleteDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.services.Inventory.DeletePortGroup(ctx, actor, input.ID); err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h inventoryHandler) batchUpdateAllocations(ctx context.Context, input *AllocationBatchUpdateInputDTO) (*BodyDTO[[]AllocationDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	groups, err := h.services.Inventory.UpdatePortGroups(ctx, ToBatchUpdateInput(actor, input.Body))
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[[]AllocationDTO]{Body: ToAllocationDTOs(groups)}, nil
}

func (h inventoryHandler) batchDeleteAllocations(ctx context.Context, input *AllocationBatchDeleteInputDTO) (*BodyDTO[DeleteDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.services.Inventory.DeletePortGroups(ctx, ToBatchDeleteInput(actor, input.Body)); err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h inventoryHandler) exportAllocations(ctx context.Context, input *AllocationListInputDTO) (*CSVOutputDTO, error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	body, err := h.services.Inventory.ExportPortGroupsCSV(ctx, service.PortGroupListParams{Actor: actor, UserID: input.UserID, Query: input.Query, Sort: input.Sort, Status: input.Status, ProjectName: input.ProjectName, DindIP: input.DindIP})
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &CSVOutputDTO{
		ContentType:        "text/csv; charset=utf-8",
		ContentDisposition: `attachment; filename="miniport-allocations.csv"`,
		Body:               body,
	}, nil
}

func (h inventoryHandler) actor(ctx context.Context, sessionID string) (service.InventoryActor, error) {
	session, err := h.auth.session(ctx, sessionID)
	if err != nil {
		return service.InventoryActor{}, huma.Error401Unauthorized("unauthorized")
	}
	return ToInventoryActor(session), nil
}
