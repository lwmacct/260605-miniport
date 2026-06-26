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
	huma.Register(api, huma.Operation{OperationID: "list-hosts", Method: http.MethodGet, Path: "/hosts", Summary: "List hosts", Tags: []string{"inventory"}}, handler.listHosts)
	huma.Register(api, huma.Operation{OperationID: "create-host", Method: http.MethodPost, Path: "/hosts", Summary: "Create host", Tags: []string{"inventory"}}, handler.createHost)
	huma.Register(api, huma.Operation{OperationID: "update-host", Method: http.MethodPut, Path: "/hosts/{id}", Summary: "Update host", Tags: []string{"inventory"}}, handler.updateHost)
	huma.Register(api, huma.Operation{OperationID: "delete-host", Method: http.MethodDelete, Path: "/hosts/{id}", Summary: "Delete host", Tags: []string{"inventory"}}, handler.deleteHost)
	huma.Register(api, huma.Operation{OperationID: "list-port-groups", Method: http.MethodGet, Path: "/port-groups", Summary: "List port groups", Tags: []string{"inventory"}}, handler.listPortGroups)
	huma.Register(api, huma.Operation{OperationID: "get-port-group", Method: http.MethodGet, Path: "/port-groups/{id}", Summary: "Get port group", Tags: []string{"inventory"}}, handler.getPortGroup)
	huma.Register(api, huma.Operation{OperationID: "create-port-group", Method: http.MethodPost, Path: "/port-groups", Summary: "Create port group", Tags: []string{"inventory"}}, handler.createPortGroup)
	huma.Register(api, huma.Operation{OperationID: "update-port-group", Method: http.MethodPut, Path: "/port-groups/{id}", Summary: "Update port group", Tags: []string{"inventory"}}, handler.updatePortGroup)
	huma.Register(api, huma.Operation{OperationID: "delete-port-group", Method: http.MethodDelete, Path: "/port-groups/{id}", Summary: "Delete port group", Tags: []string{"inventory"}}, handler.deletePortGroup)
	huma.Register(api, huma.Operation{OperationID: "batch-update-port-groups", Method: http.MethodPost, Path: "/port-groups/batch-update", Summary: "Batch update port groups", Tags: []string{"inventory"}}, handler.batchUpdatePortGroups)
	huma.Register(api, huma.Operation{OperationID: "batch-delete-port-groups", Method: http.MethodPost, Path: "/port-groups/batch-delete", Summary: "Batch delete port groups", Tags: []string{"inventory"}}, handler.batchDeletePortGroups)
	huma.Register(api, huma.Operation{OperationID: "export-port-groups", Method: http.MethodGet, Path: "/exports/port-groups.csv", Summary: "Export port groups", Tags: []string{"inventory"}}, handler.exportPortGroups)
}

func (h inventoryHandler) listHosts(ctx context.Context, input *HostListInputDTO) (*BodyDTO[[]HostDTO], error) {
	if _, err := h.auth.session(ctx, input.Session); err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	hosts, err := h.services.Inventory.ListHosts(ctx, service.HostListParams{Environment: input.Environment, Query: input.Query, Sort: input.Sort})
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[[]HostDTO]{Body: ToHostDTOs(hosts)}, nil
}

func (h inventoryHandler) createHost(ctx context.Context, input *HostBodyInputDTO) (*BodyDTO[HostDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	host, err := h.services.Inventory.CreateHost(ctx, input.Body)
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[HostDTO]{Body: ToHostDTO(*host)}, nil
}

func (h inventoryHandler) updateHost(ctx context.Context, input *HostUpdateInputDTO) (*BodyDTO[HostDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	host, err := h.services.Inventory.UpdateHost(ctx, input.ID, input.Body)
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[HostDTO]{Body: ToHostDTO(*host)}, nil
}

func (h inventoryHandler) deleteHost(ctx context.Context, input *HostInputDTO) (*BodyDTO[DeleteDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	if err := h.services.Inventory.DeleteHost(ctx, input.ID); err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h inventoryHandler) listPortGroups(ctx context.Context, input *PortGroupListInputDTO) (*BodyDTO[[]PortGroupDTO], error) {
	if _, err := h.auth.session(ctx, input.Session); err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	groups, err := h.services.Inventory.ListPortGroups(ctx, service.PortGroupListParams{HostID: input.HostID, Query: input.Query, Sort: input.Sort, Status: input.Status})
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[[]PortGroupDTO]{Body: ToPortGroupDTOs(groups)}, nil
}

func (h inventoryHandler) getPortGroup(ctx context.Context, input *PortGroupInputDTO) (*BodyDTO[PortGroupDTO], error) {
	if _, err := h.auth.session(ctx, input.Session); err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	group, err := h.services.Inventory.GetPortGroup(ctx, input.ID)
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[PortGroupDTO]{Body: ToPortGroupDTO(*group)}, nil
}

func (h inventoryHandler) createPortGroup(ctx context.Context, input *PortGroupBodyInputDTO) (*BodyDTO[PortGroupDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	group, err := h.services.Inventory.CreatePortGroup(ctx, input.Body)
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[PortGroupDTO]{Body: ToPortGroupDTO(*group)}, nil
}

func (h inventoryHandler) updatePortGroup(ctx context.Context, input *PortGroupUpdateInputDTO) (*BodyDTO[PortGroupDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	group, err := h.services.Inventory.UpdatePortGroup(ctx, input.ID, input.Body)
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[PortGroupDTO]{Body: ToPortGroupDTO(*group)}, nil
}

func (h inventoryHandler) deletePortGroup(ctx context.Context, input *PortGroupInputDTO) (*BodyDTO[DeleteDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	if err := h.services.Inventory.DeletePortGroup(ctx, input.ID); err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h inventoryHandler) batchUpdatePortGroups(ctx context.Context, input *PortGroupBatchUpdateInputDTO) (*BodyDTO[[]PortGroupDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	groups, err := h.services.Inventory.UpdatePortGroups(ctx, input.Body)
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[[]PortGroupDTO]{Body: ToPortGroupDTOs(groups)}, nil
}

func (h inventoryHandler) batchDeletePortGroups(ctx context.Context, input *PortGroupBatchDeleteInputDTO) (*BodyDTO[DeleteDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	if err := h.services.Inventory.DeletePortGroups(ctx, input.Body); err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h inventoryHandler) exportPortGroups(ctx context.Context, input *PortGroupListInputDTO) (*CSVOutputDTO, error) {
	if _, err := h.auth.session(ctx, input.Session); err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	body, err := h.services.Inventory.ExportPortGroupsCSV(ctx, service.PortGroupListParams{HostID: input.HostID, Query: input.Query, Sort: input.Sort, Status: input.Status})
	if err != nil {
		return nil, utilInventoryAPIError(err)
	}
	return &CSVOutputDTO{
		ContentType:        "text/csv; charset=utf-8",
		ContentDisposition: `attachment; filename="miniport-port-groups.csv"`,
		Body:               body,
	}, nil
}

func (h inventoryHandler) requireAdmin(ctx context.Context, sessionID string) (*AuthUserDTO, error) {
	user, err := h.auth.session(ctx, sessionID)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	if user.User == nil || !user.User.Admin {
		return nil, huma.Error403Forbidden("forbidden")
	}
	return user.User, nil
}
