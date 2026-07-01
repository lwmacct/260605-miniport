package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/identity"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/requestctx"

	"github.com/lwmacct/260605-miniport/internal/service"
)

type portsvcHandler struct {
	config   Config
	services Services
}

func RegisterPortsvc(api huma.API, config Config, services Services) {
	handler := portsvcHandler{
		config:   config,
		services: services,
	}
	huma.Register(api, huma.Operation{OperationID: "list-hosts", Method: http.MethodGet, Path: "/hosts", Summary: "List hosts", Tags: []string{"hosts"}}, handler.listHosts)
	huma.Register(api, huma.Operation{OperationID: "create-host", Method: http.MethodPost, Path: "/hosts", Summary: "Create host", Tags: []string{"hosts"}}, handler.createHost)
	huma.Register(api, huma.Operation{OperationID: "update-host", Method: http.MethodPut, Path: "/hosts/{id}", Summary: "Update host", Tags: []string{"hosts"}}, handler.updateHost)
	huma.Register(api, huma.Operation{OperationID: "delete-host", Method: http.MethodDelete, Path: "/hosts/{id}", Summary: "Delete host", Tags: []string{"hosts"}}, handler.deleteHost)

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
	if _, err := h.actor(ctx, input.Session); err != nil {
		return nil, err
	}
	hosts, err := h.services.Portsvc.ListHosts(ctx, service.HostListParams{Query: input.Query, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[[]HostDTO]{Body: ToHostDTOs(hosts)}, nil
}

func (h portsvcHandler) createHost(ctx context.Context, input *HostBodyInputDTO) (*BodyDTO[HostDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	host, err := h.services.Portsvc.CreateHost(ctx, actor, ToHostPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[HostDTO]{Body: ToHostDTO(*host)}, nil
}

func (h portsvcHandler) updateHost(ctx context.Context, input *HostUpdateInputDTO) (*BodyDTO[HostDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	host, err := h.services.Portsvc.UpdateHost(ctx, actor, input.ID, ToHostPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[HostDTO]{Body: ToHostDTO(*host)}, nil
}

func (h portsvcHandler) deleteHost(ctx context.Context, input *HostInputDTO) (*BodyDTO[DeleteDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.services.Portsvc.DeleteHost(ctx, actor, input.ID); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h portsvcHandler) listPortGroups(ctx context.Context, input *PortGroupListInputDTO) (*BodyDTO[[]PortGroupDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	groups, err := h.services.Portsvc.ListPortGroups(ctx, service.PortGroupListParams{Actor: actor, OwnerSubject: input.OwnerSubject, Query: input.Query, Sort: input.Sort, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[[]PortGroupDTO]{Body: ToPortGroupDTOs(groups)}, nil
}

func (h portsvcHandler) getPortGroup(ctx context.Context, input *PortGroupInputDTO) (*BodyDTO[PortGroupDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	group, err := h.services.Portsvc.GetPortGroup(ctx, actor, input.ID)
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortGroupDTO]{Body: ToPortGroupDTO(*group)}, nil
}

func (h portsvcHandler) createPortGroup(ctx context.Context, input *PortGroupBodyInputDTO) (*BodyDTO[PortGroupDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	group, err := h.services.Portsvc.CreatePortGroup(ctx, actor, ToPortGroupPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortGroupDTO]{Body: ToPortGroupDTO(*group)}, nil
}

func (h portsvcHandler) updatePortGroup(ctx context.Context, input *PortGroupUpdateInputDTO) (*BodyDTO[PortGroupDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	group, err := h.services.Portsvc.UpdatePortGroup(ctx, actor, input.ID, ToPortGroupPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortGroupDTO]{Body: ToPortGroupDTO(*group)}, nil
}

func (h portsvcHandler) deletePortGroup(ctx context.Context, input *PortGroupInputDTO) (*BodyDTO[DeleteDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.services.Portsvc.DeletePortGroup(ctx, actor, input.ID); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h portsvcHandler) createPortSlot(ctx context.Context, input *PortSlotBodyInputDTO) (*BodyDTO[PortSlotDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	slot, err := h.services.Portsvc.CreatePortSlot(ctx, actor, input.ID, ToPortSlotPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortSlotDTO]{Body: ToPortSlotDTO(*slot)}, nil
}

func (h portsvcHandler) updatePortSlot(ctx context.Context, input *PortSlotUpdateInputDTO) (*BodyDTO[PortSlotDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	slot, err := h.services.Portsvc.UpdatePortSlot(ctx, actor, input.ID, ToPortSlotPayload(input.Body))
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[PortSlotDTO]{Body: ToPortSlotDTO(*slot)}, nil
}

func (h portsvcHandler) deletePortSlot(ctx context.Context, input *PortGroupInputDTO) (*BodyDTO[DeleteDTO], error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.services.Portsvc.DeletePortSlot(ctx, actor, input.ID); err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &BodyDTO[DeleteDTO]{Body: DeleteDTO{Deleted: true}}, nil
}

func (h portsvcHandler) exportPortGroups(ctx context.Context, input *PortGroupListInputDTO) (*CSVOutputDTO, error) {
	actor, err := h.actor(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	body, err := h.services.Portsvc.ExportPortGroupsCSV(ctx, service.PortGroupListParams{Actor: actor, OwnerSubject: input.OwnerSubject, Query: input.Query, Sort: input.Sort, Status: input.Status})
	if err != nil {
		return nil, utilPortsvcAPIError(err)
	}
	return &CSVOutputDTO{
		ContentType:        "text/csv; charset=utf-8",
		ContentDisposition: `attachment; filename="miniport-port-groups.csv"`,
		Body:               body,
	}, nil
}

func (h portsvcHandler) actor(ctx context.Context, sessionID string) (service.PortsvcActor, error) {
	if sessionID == "" || h.config.Identity == nil {
		return service.PortsvcActor{}, huma.Error401Unauthorized("unauthorized")
	}
	request, ok := requestctx.RequestFromContext(ctx)
	if !ok {
		return service.PortsvcActor{}, huma.Error401Unauthorized("unauthorized")
	}
	principal, err := h.config.Identity.CurrentPrincipal(ctx, sessionID, request)
	if err != nil || principal == nil {
		return service.PortsvcActor{}, huma.Error401Unauthorized("unauthorized")
	}
	if !principal.Active() || principal.Status == identity.StatusDisabled {
		return service.PortsvcActor{}, huma.Error401Unauthorized("unauthorized")
	}
	actor := ToPortsvcActor(principal)
	if actor.OwnerSubject == "" {
		return service.PortsvcActor{}, huma.Error401Unauthorized("unauthorized")
	}
	return actor, nil
}
