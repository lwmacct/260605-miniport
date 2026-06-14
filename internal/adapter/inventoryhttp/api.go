package inventoryhttp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"

	"github.com/lwmacct/260605-miniport/internal/domain/inventory"
)

type hostListInput struct {
	Environment string `query:"environment" example:"dev"`
	Query       string `query:"q" example:"172.22.11"`
	Sort        string `query:"sort" example:"ip"`
}

type hostListOutput struct {
	Body []inventory.Host `json:"body"`
}

type hostOutput struct {
	Body inventory.Host `json:"body"`
}

type hostInput struct {
	ID int64 `path:"id" example:"1"`
}

type hostBodyInput struct {
	Body inventory.HostPayload
}

type hostUpdateInput struct {
	ID   int64 `path:"id" example:"1"`
	Body inventory.HostPayload
}

type portGroupListOutput struct {
	Body []inventory.PortGroupView `json:"body"`
}

type portGroupListInput struct {
	HostID int64  `query:"hostId" example:"1"`
	Query  string `query:"q" example:"order-service"`
	Sort   string `query:"sort" example:"host_port"`
	Status string `query:"status" example:"running"`
}

type portGroupOutput struct {
	Body inventory.PortGroupView `json:"body"`
}

type portGroupInput struct {
	ID int64 `path:"id" example:"1"`
}

type portGroupBodyInput struct {
	Body inventory.PortGroupPayload
}

type portGroupUpdateInput struct {
	ID   int64 `path:"id" example:"1"`
	Body inventory.PortGroupPayload
}

type deleteOutput struct {
	Body struct {
		Deleted bool `json:"deleted" example:"true"`
	}
}

func Register(api huma.API, service *inventory.Service) {
	huma.Register(api, huma.Operation{
		OperationID: "list-hosts",
		Method:      http.MethodGet,
		Path:        "/hosts",
		Summary:     "List hosts",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *hostListInput) (*hostListOutput, error) {
		hosts, err := service.ListHosts(ctx, inventory.HostListParams{
			Environment: input.Environment,
			Query:       input.Query,
			Sort:        input.Sort,
		})
		if err != nil {
			return nil, err
		}
		return &hostListOutput{Body: hosts}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-host",
		Method:      http.MethodPost,
		Path:        "/hosts",
		Summary:     "Create host",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *hostBodyInput) (*hostOutput, error) {
		host, err := service.CreateHost(ctx, input.Body)
		if err != nil {
			return nil, err
		}
		return &hostOutput{Body: *host}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-host",
		Method:      http.MethodPut,
		Path:        "/hosts/{id}",
		Summary:     "Update host",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *hostUpdateInput) (*hostOutput, error) {
		host, err := service.UpdateHost(ctx, input.ID, input.Body)
		if err != nil {
			return nil, err
		}
		return &hostOutput{Body: *host}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-host",
		Method:      http.MethodDelete,
		Path:        "/hosts/{id}",
		Summary:     "Delete host",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *hostInput) (*deleteOutput, error) {
		if err := service.DeleteHost(ctx, input.ID); err != nil {
			return nil, err
		}
		out := &deleteOutput{}
		out.Body.Deleted = true
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-port-groups",
		Method:      http.MethodGet,
		Path:        "/port-groups",
		Summary:     "List port groups",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *portGroupListInput) (*portGroupListOutput, error) {
		groups, err := service.ListPortGroups(ctx, inventory.PortGroupListParams{
			HostID: input.HostID,
			Query:  input.Query,
			Sort:   input.Sort,
			Status: input.Status,
		})
		if err != nil {
			return nil, err
		}
		return &portGroupListOutput{Body: groups}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-port-group",
		Method:      http.MethodGet,
		Path:        "/port-groups/{id}",
		Summary:     "Get port group",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *portGroupInput) (*portGroupOutput, error) {
		group, err := service.GetPortGroup(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return &portGroupOutput{Body: *group}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-port-group",
		Method:      http.MethodPost,
		Path:        "/port-groups",
		Summary:     "Create port group",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *portGroupBodyInput) (*portGroupOutput, error) {
		group, err := service.CreatePortGroup(ctx, input.Body)
		if err != nil {
			return nil, err
		}
		return &portGroupOutput{Body: *group}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-port-group",
		Method:      http.MethodPut,
		Path:        "/port-groups/{id}",
		Summary:     "Update port group",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *portGroupUpdateInput) (*portGroupOutput, error) {
		group, err := service.UpdatePortGroup(ctx, input.ID, input.Body)
		if err != nil {
			return nil, err
		}
		return &portGroupOutput{Body: *group}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-port-group",
		Method:      http.MethodDelete,
		Path:        "/port-groups/{id}",
		Summary:     "Delete port group",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *portGroupInput) (*deleteOutput, error) {
		if err := service.DeletePortGroup(ctx, input.ID); err != nil {
			return nil, err
		}
		out := &deleteOutput{}
		out.Body.Deleted = true
		return out, nil
	})
}

func RegisterExportRoutes(router chi.Router, service *inventory.Service) {
	router.Get("/exports/port-groups.csv", func(w http.ResponseWriter, r *http.Request) {
		hostID, _ := strconv.ParseInt(r.URL.Query().Get("hostId"), 10, 64)
		body, err := service.ExportPortGroupsCSV(r.Context(), inventory.PortGroupListParams{
			HostID: hostID,
			Query:  r.URL.Query().Get("q"),
			Sort:   r.URL.Query().Get("sort"),
			Status: r.URL.Query().Get("status"),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", `attachment; filename="miniport-port-groups.csv"`)
		_, _ = w.Write(body)
	})
}
