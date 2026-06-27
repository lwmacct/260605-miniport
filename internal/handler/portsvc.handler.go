package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260605-miniport/internal/service"
)

type portsvcHandler struct {
	auth     authHandler
	services Services
}

func RegisterPortsvc(api huma.API, config Config, services Services) {
	handler := portsvcHandler{
		auth:     authHandler{config: config, services: services},
		services: services,
	}
	huma.Register(api, huma.Operation{OperationID: "list-services", Method: http.MethodGet, Path: "/services", Summary: "List services", Tags: []string{"services"}}, handler.listServices)
	huma.Register(api, huma.Operation{OperationID: "get-service", Method: http.MethodGet, Path: "/services/{id}", Summary: "Get service", Tags: []string{"services"}}, handler.getService)
	huma.Register(api, huma.Operation{OperationID: "create-service", Method: http.MethodPost, Path: "/services", Summary: "Create service", Tags: []string{"services"}}, handler.createService)
	huma.Register(api, huma.Operation{OperationID: "update-service", Method: http.MethodPut, Path: "/services/{id}", Summary: "Update service", Tags: []string{"services"}}, handler.updateService)
	huma.Register(api, huma.Operation{OperationID: "delete-service", Method: http.MethodDelete, Path: "/services/{id}", Summary: "Delete service", Tags: []string{"services"}}, handler.deleteService)
	huma.Register(api, huma.Operation{OperationID: "batch-delete-services", Method: http.MethodPost, Path: "/services/batch-delete", Summary: "Batch delete services", Tags: []string{"services"}}, handler.batchDeleteServices)
	huma.Register(api, huma.Operation{OperationID: "list-port-allocations", Method: http.MethodGet, Path: "/port-allocations", Summary: "List port allocations", Tags: []string{"port-allocations"}}, handler.listPortAllocations)
	huma.Register(api, huma.Operation{OperationID: "create-port-allocation", Method: http.MethodPost, Path: "/port-allocations", Summary: "Create port allocation", Tags: []string{"port-allocations"}}, handler.createPortAllocation)
	huma.Register(api, huma.Operation{OperationID: "update-port-allocation", Method: http.MethodPut, Path: "/port-allocations/{id}", Summary: "Update port allocation", Tags: []string{"port-allocations"}}, handler.updatePortAllocation)
	huma.Register(api, huma.Operation{OperationID: "delete-port-allocation", Method: http.MethodDelete, Path: "/port-allocations/{id}", Summary: "Delete port allocation", Tags: []string{"port-allocations"}}, handler.deletePortAllocation)
	huma.Register(api, huma.Operation{OperationID: "export-services", Method: http.MethodGet, Path: "/services/export.csv", Summary: "Export services", Tags: []string{"services"}}, handler.exportServices)
}

func (h portsvcHandler) listServices(ctx context.Context, input *ServiceListInputDTO) (*BodyDTO[[]ServiceDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	services, err := h.services.Portsvc.ListServices(ctx, service.ServiceListParams{Actor: actor, UserID: input.UserID, Query: input.Query, Sort: input.Sort, Status: input.Status, ProjectName: input.ProjectName})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[[]ServiceDTO]{Body: ToServiceDTOs(services)}, nil
}

func (h portsvcHandler) getService(ctx context.Context, input *ServiceInputDTO) (*BodyDTO[ServiceDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	item, err := h.services.Portsvc.GetService(ctx, actor, input.ID)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[ServiceDTO]{Body: ToServiceDTO(*item)}, nil
}

func (h portsvcHandler) createService(ctx context.Context, input *ServiceBodyInputDTO) (*BodyDTO[ServiceDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	item, err := h.services.Portsvc.CreateService(ctx, actor, ToServicePayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[ServiceDTO]{Body: ToServiceDTO(*item)}, nil
}

func (h portsvcHandler) updateService(ctx context.Context, input *ServiceUpdateInputDTO) (*BodyDTO[ServiceDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	item, err := h.services.Portsvc.UpdateService(ctx, actor, input.ID, ToServicePayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[ServiceDTO]{Body: ToServiceDTO(*item)}, nil
}

func (h portsvcHandler) deleteService(ctx context.Context, input *ServiceInputDTO) (*BodyDTO[DeleteDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.services.Portsvc.DeleteService(ctx, actor, input.ID); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h portsvcHandler) batchDeleteServices(ctx context.Context, input *ServiceBatchDeleteInputDTO) (*BodyDTO[DeleteDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.services.Portsvc.DeleteServices(ctx, ToServiceBatchDeleteInput(actor, input.Body.IDs)); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h portsvcHandler) listPortAllocations(ctx context.Context, input *PortAllocationListInputDTO) (*BodyDTO[[]PortAllocationDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	groups, err := h.services.Portsvc.ListPortAllocations(ctx, service.PortAllocationListParams{Actor: actor, UserID: input.UserID, Sort: input.Sort, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[[]PortAllocationDTO]{Body: ToPortAllocationDTOs(groups)}, nil
}

func (h portsvcHandler) createPortAllocation(ctx context.Context, input *PortAllocationBodyInputDTO) (*BodyDTO[PortAllocationDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	group, err := h.services.Portsvc.CreatePortAllocation(ctx, actor, ToPortAllocationPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortAllocationDTO]{Body: ToPortAllocationDTO(*group)}, nil
}

func (h portsvcHandler) updatePortAllocation(ctx context.Context, input *PortAllocationUpdateInputDTO) (*BodyDTO[PortAllocationDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	group, err := h.services.Portsvc.UpdatePortAllocation(ctx, actor, input.ID, ToPortAllocationPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortAllocationDTO]{Body: ToPortAllocationDTO(*group)}, nil
}

func (h portsvcHandler) deletePortAllocation(ctx context.Context, input *PortAllocationInputDTO) (*BodyDTO[DeleteDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.services.Portsvc.DeletePortAllocation(ctx, actor, input.ID); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h portsvcHandler) exportServices(ctx context.Context, input *ServiceListInputDTO) (*CSVOutputDTO, error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	body, err := h.services.Portsvc.ExportServicesCSV(ctx, service.ServiceListParams{Actor: actor, UserID: input.UserID, Query: input.Query, Sort: input.Sort, Status: input.Status, ProjectName: input.ProjectName})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &CSVOutputDTO{
		ContentType:        "text/csv; charset=utf-8",
		ContentDisposition: `attachment; filename="miniport-services.csv"`,
		Body:               body,
	}, nil
}

func (h portsvcHandler) actor(ctx context.Context, sessionID string) (service.PortsvcActor, error) {
	session, err := h.auth.session(ctx, sessionID)
	if err != nil {
		return service.PortsvcActor{}, huma.Error401Unauthorized("unauthorized")
	}
	return ToPortsvcActor(session), nil
}
