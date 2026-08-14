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
	h := portsvcHandler{services: services}
	g := huma.NewGroup(api, "/console")
	console := []string{"Console"}
	huma.Register(g, huma.Operation{OperationID: "console-list-hosts", Method: http.MethodGet, Path: "/hosts", Summary: "List hosts", Tags: console}, h.listHosts)
	huma.Register(g, huma.Operation{OperationID: "console-create-hosts", Method: http.MethodPost, Path: "/hosts", DefaultStatus: http.StatusCreated, Summary: "Create hosts", Tags: console}, h.createHosts)
	huma.Register(g, huma.Operation{OperationID: "console-update-hosts", Method: http.MethodPut, Path: "/hosts", Summary: "Update hosts", Tags: console}, h.updateHosts)
	huma.Register(g, huma.Operation{OperationID: "console-delete-hosts", Method: http.MethodDelete, Path: "/hosts", Summary: "Delete hosts", Tags: console}, h.deleteHosts)

	huma.Register(g, huma.Operation{OperationID: "console-list-dependency-assets", Method: http.MethodGet, Path: "/dependency-assets", Summary: "List dependency assets", Tags: console}, h.listDependencyAssets)
	huma.Register(g, huma.Operation{OperationID: "console-create-dependency-assets", Method: http.MethodPost, Path: "/dependency-assets", DefaultStatus: http.StatusCreated, Summary: "Create dependency assets", Tags: console}, h.createDependencyAssets)
	huma.Register(g, huma.Operation{OperationID: "console-update-dependency-assets", Method: http.MethodPut, Path: "/dependency-assets", Summary: "Update dependency assets", Tags: console}, h.updateDependencyAssets)
	huma.Register(g, huma.Operation{OperationID: "console-delete-dependency-assets", Method: http.MethodDelete, Path: "/dependency-assets", Summary: "Delete dependency assets", Tags: console}, h.deleteDependencyAssets)

	huma.Register(g, huma.Operation{OperationID: "console-list-service-groups", Method: http.MethodGet, Path: "/service-groups", Summary: "List service groups", Tags: console}, h.listServiceGroups)
	huma.Register(g, huma.Operation{OperationID: "console-create-service-groups", Method: http.MethodPost, Path: "/service-groups", DefaultStatus: http.StatusCreated, Summary: "Create service groups", Tags: console}, h.createServiceGroups)
	huma.Register(g, huma.Operation{OperationID: "console-update-service-groups", Method: http.MethodPut, Path: "/service-groups", Summary: "Update service groups", Tags: console}, h.updateServiceGroups)
	huma.Register(g, huma.Operation{OperationID: "console-delete-service-groups", Method: http.MethodDelete, Path: "/service-groups", Summary: "Delete service groups", Tags: console}, h.deleteServiceGroups)

	huma.Register(g, huma.Operation{OperationID: "console-list-port-groups", Method: http.MethodGet, Path: "/port-groups", Summary: "List port groups", Tags: console}, h.listPortGroups)
	huma.Register(g, huma.Operation{OperationID: "console-create-port-groups", Method: http.MethodPost, Path: "/port-groups", DefaultStatus: http.StatusCreated, Summary: "Create port groups", Tags: console}, h.createPortGroups)
	huma.Register(g, huma.Operation{OperationID: "console-update-port-groups", Method: http.MethodPut, Path: "/port-groups", Summary: "Update port groups", Tags: console}, h.updatePortGroups)
	huma.Register(g, huma.Operation{OperationID: "console-delete-port-groups", Method: http.MethodDelete, Path: "/port-groups", Summary: "Delete port groups", Tags: console}, h.deletePortGroups)
	huma.Register(g, huma.Operation{
		OperationID: "console-export-port-groups",
		Method:      http.MethodGet,
		Path:        "/port-groups/export.csv",
		Summary:     "Export port groups",
		Tags:        console,
		Responses: map[string]*huma.Response{
			"200": {
				Description: "CSV export",
				Content: map[string]*huma.MediaType{
					"text/csv": {Schema: &huma.Schema{Type: "string"}},
				},
			},
		},
	}, h.exportPortGroups)
}

func (h portsvcHandler) listHosts(ctx context.Context, input *HostListInputDTO) (*BodyDTO[HostListDTO], error) {
	items, err := h.services.Portsvc.ListHosts(ctx, service.HostListParams{Query: input.Query, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[HostListDTO]{Body: HostListDTO{Items: ToHostDTOs(items)}}, nil
}

func (h portsvcHandler) createHosts(ctx context.Context, input *HostCreateInputDTO) (*BodyDTO[HostBatchDTO], error) {
	payloads := make([]service.HostPayload, 0, len(input.Body.Items))
	for _, item := range input.Body.Items {
		payloads = append(payloads, ToHostCreatePayload(item))
	}
	items, err := h.services.Portsvc.CreateHosts(ctx, payloads)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[HostBatchDTO]{Body: HostBatchDTO{Items: ToHostDTOs(items)}}, nil
}

func (h portsvcHandler) updateHosts(ctx context.Context, input *HostUpdateInputDTO) (*BodyDTO[HostBatchDTO], error) {
	items := make([]service.HostUpdateInput, 0, len(input.Body.Items))
	for _, item := range input.Body.Items {
		items = append(items, ToHostUpdatePayload(item))
	}
	values, err := h.services.Portsvc.UpdateHosts(ctx, items)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[HostBatchDTO]{Body: HostBatchDTO{Items: ToHostDTOs(values)}}, nil
}

func (h portsvcHandler) deleteHosts(ctx context.Context, input *HostDeleteInputDTO) (*ActionOutputDTO, error) {
	if err := h.services.Portsvc.DeleteHosts(ctx, input.Body.IDs); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &ActionOutputDTO{Body: ActionDTO{OK: true}}, nil
}

func (h portsvcHandler) listDependencyAssets(ctx context.Context, input *DependencyAssetListInputDTO) (*BodyDTO[DependencyAssetListDTO], error) {
	items, err := h.services.Portsvc.ListDependencyAssets(ctx, service.DependencyAssetListParams{Query: input.Query, AssetKind: input.AssetKind, AssetType: input.AssetType, Provider: input.Provider, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DependencyAssetListDTO]{Body: DependencyAssetListDTO{Items: ToDependencyAssetDTOs(items)}}, nil
}

func (h portsvcHandler) createDependencyAssets(ctx context.Context, input *DependencyAssetCreateInputDTO) (*BodyDTO[DependencyAssetBatchDTO], error) {
	payloads := make([]service.DependencyAssetPayload, 0, len(input.Body.Items))
	for _, item := range input.Body.Items {
		payloads = append(payloads, ToDependencyAssetCreatePayload(item))
	}
	items, err := h.services.Portsvc.CreateDependencyAssets(ctx, payloads)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DependencyAssetBatchDTO]{Body: DependencyAssetBatchDTO{Items: ToDependencyAssetDTOs(items)}}, nil
}

func (h portsvcHandler) updateDependencyAssets(ctx context.Context, input *DependencyAssetUpdateInputDTO) (*BodyDTO[DependencyAssetBatchDTO], error) {
	items := make([]service.DependencyAssetUpdateInput, 0, len(input.Body.Items))
	for _, item := range input.Body.Items {
		items = append(items, ToDependencyAssetUpdatePayload(item))
	}
	values, err := h.services.Portsvc.UpdateDependencyAssets(ctx, items)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DependencyAssetBatchDTO]{Body: DependencyAssetBatchDTO{Items: ToDependencyAssetDTOs(values)}}, nil
}

func (h portsvcHandler) deleteDependencyAssets(ctx context.Context, input *DependencyAssetDeleteInputDTO) (*ActionOutputDTO, error) {
	if err := h.services.Portsvc.DeleteDependencyAssets(ctx, input.Body.IDs); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &ActionOutputDTO{Body: ActionDTO{OK: true}}, nil
}

func (h portsvcHandler) listServiceGroups(ctx context.Context, input *ServiceGroupListInputDTO) (*BodyDTO[ServiceGroupListDTO], error) {
	items, err := h.services.Portsvc.ListServiceGroups(ctx, service.ServiceGroupListParams{Query: input.Query, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[ServiceGroupListDTO]{Body: ServiceGroupListDTO{Items: ToServiceGroupDTOs(items)}}, nil
}

func (h portsvcHandler) createServiceGroups(ctx context.Context, input *ServiceGroupCreateInputDTO) (*BodyDTO[ServiceGroupBatchDTO], error) {
	payloads := make([]service.ServiceGroupPayload, 0, len(input.Body.Items))
	for _, item := range input.Body.Items {
		payloads = append(payloads, ToServiceGroupCreatePayload(item))
	}
	items, err := h.services.Portsvc.CreateServiceGroups(ctx, payloads)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[ServiceGroupBatchDTO]{Body: ServiceGroupBatchDTO{Items: ToServiceGroupDTOs(items)}}, nil
}

func (h portsvcHandler) updateServiceGroups(ctx context.Context, input *ServiceGroupUpdateInputDTO) (*BodyDTO[ServiceGroupBatchDTO], error) {
	items := make([]service.ServiceGroupUpdateInput, 0, len(input.Body.Items))
	for _, item := range input.Body.Items {
		items = append(items, ToServiceGroupUpdatePayload(item))
	}
	values, err := h.services.Portsvc.UpdateServiceGroups(ctx, items)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[ServiceGroupBatchDTO]{Body: ServiceGroupBatchDTO{Items: ToServiceGroupDTOs(values)}}, nil
}

func (h portsvcHandler) deleteServiceGroups(ctx context.Context, input *ServiceGroupDeleteInputDTO) (*ActionOutputDTO, error) {
	if err := h.services.Portsvc.DeleteServiceGroups(ctx, input.Body.IDs); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &ActionOutputDTO{Body: ActionDTO{OK: true}}, nil
}

func (h portsvcHandler) listPortGroups(ctx context.Context, input *PortGroupListInputDTO) (*BodyDTO[PortGroupListDTO], error) {
	items, err := h.services.Portsvc.ListPortGroups(ctx, service.PortGroupListParams{Query: input.Query, Sort: input.Sort, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortGroupListDTO]{Body: PortGroupListDTO{Items: ToPortGroupDTOs(items)}}, nil
}

func (h portsvcHandler) createPortGroups(ctx context.Context, input *PortGroupCreateInputDTO) (*BodyDTO[PortGroupBatchDTO], error) {
	payloads := make([]service.PortGroupPayload, 0, len(input.Body.Items))
	for _, item := range input.Body.Items {
		payloads = append(payloads, ToPortGroupCreatePayload(item))
	}
	items, err := h.services.Portsvc.CreatePortGroups(ctx, payloads)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortGroupBatchDTO]{Body: PortGroupBatchDTO{Items: ToPortGroupDTOs(items)}}, nil
}

func (h portsvcHandler) updatePortGroups(ctx context.Context, input *PortGroupUpdateInputDTO) (*BodyDTO[PortGroupBatchDTO], error) {
	items := make([]service.PortGroupUpdateInput, 0, len(input.Body.Items))
	for _, item := range input.Body.Items {
		items = append(items, ToPortGroupUpdatePayload(item))
	}
	values, err := h.services.Portsvc.UpdatePortGroups(ctx, items)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortGroupBatchDTO]{Body: PortGroupBatchDTO{Items: ToPortGroupDTOs(values)}}, nil
}

func (h portsvcHandler) deletePortGroups(ctx context.Context, input *PortGroupDeleteInputDTO) (*ActionOutputDTO, error) {
	if err := h.services.Portsvc.DeletePortGroups(ctx, input.Body.IDs); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &ActionOutputDTO{Body: ActionDTO{OK: true}}, nil
}

func (h portsvcHandler) exportPortGroups(ctx context.Context, input *PortGroupListInputDTO) (*CSVOutputDTO, error) {
	body, err := h.services.Portsvc.ExportPortGroupsCSV(ctx, service.PortGroupListParams{Query: input.Query, Sort: input.Sort, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &CSVOutputDTO{ContentType: "text/csv; charset=utf-8", ContentDisposition: `attachment; filename="miniport-port-groups.csv"`, Body: body}, nil
}
