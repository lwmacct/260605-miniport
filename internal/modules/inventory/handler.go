package inventory

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *Service
}

func RegisterAPI(api huma.API, service *Service) {
	handler := &Handler{service: service}
	huma.Register(api, huma.Operation{
		OperationID: "list-hosts",
		Method:      http.MethodGet,
		Path:        "/hosts",
		Summary:     "List IP hosts",
		Tags:        []string{"inventory"},
	}, handler.listHosts)

	huma.Register(api, huma.Operation{
		OperationID: "create-host",
		Method:      http.MethodPost,
		Path:        "/hosts",
		Summary:     "Create IP host",
		Tags:        []string{"inventory"},
	}, handler.createHost)

	huma.Register(api, huma.Operation{
		OperationID: "update-host",
		Method:      http.MethodPut,
		Path:        "/hosts/{id}",
		Summary:     "Update IP host",
		Tags:        []string{"inventory"},
	}, handler.updateHost)

	huma.Register(api, huma.Operation{
		OperationID: "delete-host",
		Method:      http.MethodDelete,
		Path:        "/hosts/{id}",
		Summary:     "Delete IP host",
		Tags:        []string{"inventory"},
	}, handler.deleteHost)

	huma.Register(api, huma.Operation{
		OperationID: "list-port-groups",
		Method:      http.MethodGet,
		Path:        "/port-groups",
		Summary:     "List port groups",
		Tags:        []string{"inventory"},
	}, handler.listPortGroups)

	huma.Register(api, huma.Operation{
		OperationID: "get-port-group",
		Method:      http.MethodGet,
		Path:        "/port-groups/{id}",
		Summary:     "Get port group",
		Tags:        []string{"inventory"},
	}, handler.getPortGroup)

	huma.Register(api, huma.Operation{
		OperationID: "create-port-group",
		Method:      http.MethodPost,
		Path:        "/port-groups",
		Summary:     "Create port group",
		Tags:        []string{"inventory"},
	}, handler.createPortGroup)

	huma.Register(api, huma.Operation{
		OperationID: "update-port-group",
		Method:      http.MethodPut,
		Path:        "/port-groups/{id}",
		Summary:     "Update port group",
		Tags:        []string{"inventory"},
	}, handler.updatePortGroup)

	huma.Register(api, huma.Operation{
		OperationID: "delete-port-group",
		Method:      http.MethodDelete,
		Path:        "/port-groups/{id}",
		Summary:     "Delete port group",
		Tags:        []string{"inventory"},
	}, handler.deletePortGroup)
}

func (h *Handler) listHosts(ctx context.Context, _ *struct{}) (*hostListOutput, error) {
	hosts, err := h.service.ListHosts(ctx)
	if err != nil {
		return nil, err
	}
	return &hostListOutput{Body: hosts}, nil
}

func (h *Handler) createHost(ctx context.Context, input *hostBodyInput) (*hostOutput, error) {
	host, err := h.service.CreateHost(ctx, input.Body)
	if err != nil {
		return nil, err
	}
	return &hostOutput{Body: *host}, nil
}

func (h *Handler) updateHost(ctx context.Context, input *hostUpdateInput) (*hostOutput, error) {
	host, err := h.service.UpdateHost(ctx, input.ID, input.Body)
	if err != nil {
		return nil, err
	}
	return &hostOutput{Body: *host}, nil
}

func (h *Handler) deleteHost(ctx context.Context, input *hostInput) (*deleteOutput, error) {
	if err := h.service.DeleteHost(ctx, input.ID); err != nil {
		return nil, err
	}
	out := &deleteOutput{}
	out.Body.Deleted = true
	return out, nil
}

func (h *Handler) listPortGroups(ctx context.Context, _ *struct{}) (*portGroupListOutput, error) {
	groups, err := h.service.ListPortGroups(ctx)
	if err != nil {
		return nil, err
	}
	return &portGroupListOutput{Body: groups}, nil
}

func (h *Handler) getPortGroup(ctx context.Context, input *portGroupInput) (*portGroupOutput, error) {
	group, err := h.service.GetPortGroupView(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return &portGroupOutput{Body: *group}, nil
}

func (h *Handler) createPortGroup(ctx context.Context, input *portGroupBodyInput) (*portGroupOutput, error) {
	group, err := h.service.CreatePortGroup(ctx, input.Body)
	if err != nil {
		return nil, err
	}
	return &portGroupOutput{Body: *group}, nil
}

func (h *Handler) updatePortGroup(ctx context.Context, input *portGroupUpdateInput) (*portGroupOutput, error) {
	group, err := h.service.UpdatePortGroup(ctx, input.ID, input.Body)
	if err != nil {
		return nil, err
	}
	return &portGroupOutput{Body: *group}, nil
}

func (h *Handler) deletePortGroup(ctx context.Context, input *portGroupInput) (*deleteOutput, error) {
	if err := h.service.DeletePortGroup(ctx, input.ID); err != nil {
		return nil, err
	}
	out := &deleteOutput{}
	out.Body.Deleted = true
	return out, nil
}
