package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260605-miniport/internal/service"
)

type portsvcHandler struct {
	services Services
}

func RegisterPortsvc(api huma.API, services Services) {
	handler := portsvcHandler{
		services: services,
	}
	huma.Register(api, huma.Operation{OperationID: "list-hosts", Method: http.MethodGet, Path: "/hosts", Summary: "List hosts", Tags: []string{"hosts"}}, handler.listHosts)
	huma.Register(api, huma.Operation{OperationID: "create-host", Method: http.MethodPost, Path: "/hosts", Summary: "Create host", Tags: []string{"hosts"}}, handler.createHost)
	huma.Register(api, huma.Operation{OperationID: "update-host", Method: http.MethodPut, Path: "/hosts/{id}", Summary: "Update host", Tags: []string{"hosts"}}, handler.updateHost)
	huma.Register(api, huma.Operation{OperationID: "delete-host", Method: http.MethodDelete, Path: "/hosts/{id}", Summary: "Delete host", Tags: []string{"hosts"}}, handler.deleteHost)

	huma.Register(api, huma.Operation{OperationID: "list-dependency-assets", Method: http.MethodGet, Path: "/dependency-assets", Summary: "List dependency assets", Tags: []string{"dependency-assets"}}, handler.listDependencyAssets)
	huma.Register(api, huma.Operation{OperationID: "create-dependency-asset", Method: http.MethodPost, Path: "/dependency-assets", Summary: "Create dependency asset", Tags: []string{"dependency-assets"}}, handler.createDependencyAsset)
	huma.Register(api, huma.Operation{OperationID: "update-dependency-asset", Method: http.MethodPut, Path: "/dependency-assets/{id}", Summary: "Update dependency asset", Tags: []string{"dependency-assets"}}, handler.updateDependencyAsset)
	huma.Register(api, huma.Operation{OperationID: "delete-dependency-asset", Method: http.MethodDelete, Path: "/dependency-assets/{id}", Summary: "Delete dependency asset", Tags: []string{"dependency-assets"}}, handler.deleteDependencyAsset)

	huma.Register(api, huma.Operation{OperationID: "list-service-groups", Method: http.MethodGet, Path: "/service-groups", Summary: "List service groups", Tags: []string{"service-groups"}}, handler.listServiceGroups)
	huma.Register(api, huma.Operation{OperationID: "get-service-group", Method: http.MethodGet, Path: "/service-groups/{id}", Summary: "Get service group", Tags: []string{"service-groups"}}, handler.getServiceGroup)
	huma.Register(api, huma.Operation{OperationID: "create-service-group", Method: http.MethodPost, Path: "/service-groups", Summary: "Create service group", Tags: []string{"service-groups"}}, handler.createServiceGroup)
	huma.Register(api, huma.Operation{OperationID: "update-service-group", Method: http.MethodPut, Path: "/service-groups/{id}", Summary: "Update service group", Tags: []string{"service-groups"}}, handler.updateServiceGroup)
	huma.Register(api, huma.Operation{OperationID: "delete-service-group", Method: http.MethodDelete, Path: "/service-groups/{id}", Summary: "Delete service group", Tags: []string{"service-groups"}}, handler.deleteServiceGroup)

	huma.Register(api, huma.Operation{OperationID: "list-port-groups", Method: http.MethodGet, Path: "/port-groups", Summary: "List port groups", Tags: []string{"port-groups"}}, handler.listPortGroups)
	huma.Register(api, huma.Operation{OperationID: "get-port-group", Method: http.MethodGet, Path: "/port-groups/{id}", Summary: "Get port group", Tags: []string{"port-groups"}}, handler.getPortGroup)
	huma.Register(api, huma.Operation{OperationID: "create-port-group", Method: http.MethodPost, Path: "/port-groups", Summary: "Create port group", Tags: []string{"port-groups"}}, handler.createPortGroup)
	huma.Register(api, huma.Operation{OperationID: "update-port-group", Method: http.MethodPut, Path: "/port-groups/{id}", Summary: "Update port group", Tags: []string{"port-groups"}}, handler.updatePortGroup)
	huma.Register(api, huma.Operation{OperationID: "delete-port-group", Method: http.MethodDelete, Path: "/port-groups/{id}", Summary: "Delete port group", Tags: []string{"port-groups"}}, handler.deletePortGroup)
	huma.Register(api, huma.Operation{OperationID: "create-port-slot", Method: http.MethodPost, Path: "/port-groups/{id}/slots", Summary: "Create port slot", Tags: []string{"port-slots"}}, handler.createPortSlot)
	huma.Register(api, huma.Operation{OperationID: "update-port-slot", Method: http.MethodPut, Path: "/port-slots/{id}", Summary: "Update port slot", Tags: []string{"port-slots"}}, handler.updatePortSlot)
	huma.Register(api, huma.Operation{OperationID: "delete-port-slot", Method: http.MethodDelete, Path: "/port-slots/{id}", Summary: "Delete port slot", Tags: []string{"port-slots"}}, handler.deletePortSlot)
	huma.Register(api, huma.Operation{OperationID: "export-port-groups", Method: http.MethodGet, Path: "/port-groups/export.csv", Summary: "Export port groups", Tags: []string{"port-groups"}}, handler.exportPortGroups)
}

func (h portsvcHandler) listHosts(ctx context.Context, input *HostListInputDTO) (*BodyDTO[[]HostDTO], error) {
	hosts, err := h.services.Portsvc.ListHosts(ctx, service.HostListParams{Query: input.Query, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[[]HostDTO]{Body: ToHostDTOs(hosts)}, nil
}

func (h portsvcHandler) createHost(ctx context.Context, input *HostBodyInputDTO) (*BodyDTO[HostDTO], error) {
	host, err := h.services.Portsvc.CreateHost(ctx, ToHostPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[HostDTO]{Body: ToHostDTO(*host)}, nil
}

func (h portsvcHandler) updateHost(ctx context.Context, input *HostUpdateInputDTO) (*BodyDTO[HostDTO], error) {
	host, err := h.services.Portsvc.UpdateHost(ctx, input.ID, ToHostPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[HostDTO]{Body: ToHostDTO(*host)}, nil
}

func (h portsvcHandler) deleteHost(ctx context.Context, input *HostInputDTO) (*BodyDTO[DeleteDTO], error) {
	if err := h.services.Portsvc.DeleteHost(ctx, input.ID); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h portsvcHandler) listDependencyAssets(ctx context.Context, input *DependencyAssetListInputDTO) (*BodyDTO[[]DependencyAssetDTO], error) {
	assets, err := h.services.Portsvc.ListDependencyAssets(ctx, service.DependencyAssetListParams{
		Query:     input.Query,
		AssetKind: input.AssetKind,
		AssetType: input.AssetType,
		Provider:  input.Provider,
		Status:    input.Status,
	})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[[]DependencyAssetDTO]{Body: ToDependencyAssetDTOs(assets)}, nil
}

func (h portsvcHandler) createDependencyAsset(ctx context.Context, input *DependencyAssetBodyInputDTO) (*BodyDTO[DependencyAssetDTO], error) {
	asset, err := h.services.Portsvc.CreateDependencyAsset(ctx, ToDependencyAssetPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DependencyAssetDTO]{Body: ToDependencyAssetDTO(*asset)}, nil
}

func (h portsvcHandler) updateDependencyAsset(ctx context.Context, input *DependencyAssetUpdateInputDTO) (*BodyDTO[DependencyAssetDTO], error) {
	asset, err := h.services.Portsvc.UpdateDependencyAsset(ctx, input.ID, ToDependencyAssetPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DependencyAssetDTO]{Body: ToDependencyAssetDTO(*asset)}, nil
}

func (h portsvcHandler) deleteDependencyAsset(ctx context.Context, input *DependencyAssetInputDTO) (*BodyDTO[DeleteDTO], error) {
	if err := h.services.Portsvc.DeleteDependencyAsset(ctx, input.ID); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h portsvcHandler) listServiceGroups(ctx context.Context, input *ServiceGroupListInputDTO) (*BodyDTO[[]ServiceGroupDTO], error) {
	groups, err := h.services.Portsvc.ListServiceGroups(ctx, service.ServiceGroupListParams{
		Query:  input.Query,
		Status: input.Status,
	})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[[]ServiceGroupDTO]{Body: ToServiceGroupDTOs(groups)}, nil
}

func (h portsvcHandler) getServiceGroup(ctx context.Context, input *ServiceGroupInputDTO) (*BodyDTO[ServiceGroupDTO], error) {
	group, err := h.services.Portsvc.GetServiceGroup(ctx, input.ID)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[ServiceGroupDTO]{Body: ToServiceGroupDTO(*group)}, nil
}

func (h portsvcHandler) createServiceGroup(ctx context.Context, input *ServiceGroupBodyInputDTO) (*BodyDTO[ServiceGroupDTO], error) {
	group, err := h.services.Portsvc.CreateServiceGroup(ctx, ToServiceGroupPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[ServiceGroupDTO]{Body: ToServiceGroupDTO(*group)}, nil
}

func (h portsvcHandler) updateServiceGroup(ctx context.Context, input *ServiceGroupUpdateInputDTO) (*BodyDTO[ServiceGroupDTO], error) {
	group, err := h.services.Portsvc.UpdateServiceGroup(ctx, input.ID, ToServiceGroupPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[ServiceGroupDTO]{Body: ToServiceGroupDTO(*group)}, nil
}

func (h portsvcHandler) deleteServiceGroup(ctx context.Context, input *ServiceGroupInputDTO) (*BodyDTO[DeleteDTO], error) {
	if err := h.services.Portsvc.DeleteServiceGroup(ctx, input.ID); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h portsvcHandler) listPortGroups(ctx context.Context, input *PortGroupListInputDTO) (*BodyDTO[[]PortGroupDTO], error) {
	groups, err := h.services.Portsvc.ListPortGroups(ctx, service.PortGroupListParams{Query: input.Query, Sort: input.Sort, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[[]PortGroupDTO]{Body: ToPortGroupDTOs(groups)}, nil
}

func (h portsvcHandler) getPortGroup(ctx context.Context, input *PortGroupInputDTO) (*BodyDTO[PortGroupDTO], error) {
	group, err := h.services.Portsvc.GetPortGroup(ctx, input.ID)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortGroupDTO]{Body: ToPortGroupDTO(*group)}, nil
}

func (h portsvcHandler) createPortGroup(ctx context.Context, input *PortGroupBodyInputDTO) (*BodyDTO[PortGroupDTO], error) {
	group, err := h.services.Portsvc.CreatePortGroup(ctx, ToPortGroupPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortGroupDTO]{Body: ToPortGroupDTO(*group)}, nil
}

func (h portsvcHandler) updatePortGroup(ctx context.Context, input *PortGroupUpdateInputDTO) (*BodyDTO[PortGroupDTO], error) {
	group, err := h.services.Portsvc.UpdatePortGroup(ctx, input.ID, ToPortGroupPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortGroupDTO]{Body: ToPortGroupDTO(*group)}, nil
}

func (h portsvcHandler) deletePortGroup(ctx context.Context, input *PortGroupInputDTO) (*BodyDTO[DeleteDTO], error) {
	if err := h.services.Portsvc.DeletePortGroup(ctx, input.ID); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h portsvcHandler) createPortSlot(ctx context.Context, input *PortSlotBodyInputDTO) (*BodyDTO[PortSlotDTO], error) {
	slot, err := h.services.Portsvc.CreatePortSlot(ctx, input.ID, ToPortSlotPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortSlotDTO]{Body: ToPortSlotDTO(*slot)}, nil
}

func (h portsvcHandler) updatePortSlot(ctx context.Context, input *PortSlotUpdateInputDTO) (*BodyDTO[PortSlotDTO], error) {
	slot, err := h.services.Portsvc.UpdatePortSlot(ctx, input.ID, ToPortSlotPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortSlotDTO]{Body: ToPortSlotDTO(*slot)}, nil
}

func (h portsvcHandler) deletePortSlot(ctx context.Context, input *PortGroupInputDTO) (*BodyDTO[DeleteDTO], error) {
	if err := h.services.Portsvc.DeletePortSlot(ctx, input.ID); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h portsvcHandler) exportPortGroups(ctx context.Context, input *PortGroupListInputDTO) (*CSVOutputDTO, error) {
	body, err := h.services.Portsvc.ExportPortGroupsCSV(ctx, service.PortGroupListParams{Query: input.Query, Sort: input.Sort, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &CSVOutputDTO{
		ContentType:        "text/csv; charset=utf-8",
		ContentDisposition: `attachment; filename="miniport-port-groups.csv"`,
		Body:               body,
	}, nil
}
