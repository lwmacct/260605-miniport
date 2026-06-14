package inventory

import (
	"context"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/uptrace/bun"
)

type Service struct {
	store *store
}

func NewService(db *bun.DB) *Service {
	return &Service{store: newStore(db)}
}

func (svc *Service) ListHosts(ctx context.Context, params HostListParams) ([]Host, error) {
	return svc.store.listHosts(ctx, params)
}

func (svc *Service) CreateHost(ctx context.Context, payload HostPayload) (*Host, error) {
	host := hostFromPayload(payload)
	if err := validateHost(host); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	host.CreatedAt = now
	host.UpdatedAt = now
	if err := svc.store.createHost(ctx, host); err != nil {
		return nil, err
	}
	return host, nil
}

func (svc *Service) UpdateHost(ctx context.Context, id int64, payload HostPayload) (*Host, error) {
	host := hostFromPayload(payload)
	host.ID = id
	if err := validateHost(host); err != nil {
		return nil, err
	}

	host.UpdatedAt = time.Now().UTC()
	if err := svc.store.updateHost(ctx, id, host); err != nil {
		return nil, err
	}
	return svc.store.getHost(ctx, id)
}

func (svc *Service) DeleteHost(ctx context.Context, id int64) error {
	return svc.store.deleteHost(ctx, id)
}

func (svc *Service) ListPortGroups(ctx context.Context, params PortGroupListParams) ([]PortGroupView, error) {
	groups, err := svc.store.listPortGroups(ctx, params)
	if err != nil {
		return nil, err
	}
	return svc.store.buildPortGroupViews(ctx, groups)
}

func (svc *Service) GetPortGroup(ctx context.Context, id int64) (*PortGroupView, error) {
	group, err := svc.store.getPortGroupWithHost(ctx, id)
	if err != nil {
		return nil, err
	}
	views, err := svc.store.buildPortGroupViews(ctx, []PortGroup{*group})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (svc *Service) CreatePortGroup(ctx context.Context, payload PortGroupPayload) (*PortGroupView, error) {
	var out *PortGroupView
	err := svc.store.runInTx(ctx, func(ctx context.Context, tx *store) error {
		group, err := portGroupFromPayload(ctx, tx, 0, payload)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		group.CreatedAt = now
		group.UpdatedAt = now
		if err := tx.createPortGroup(ctx, group); err != nil {
			return err
		}
		if err := replacePortGroupChildren(ctx, tx, group.ID, payload, now); err != nil {
			return err
		}

		view, err := tx.getPortGroupView(ctx, group.ID)
		if err != nil {
			return err
		}
		out = view
		return nil
	})
	return out, err
}

func (svc *Service) UpdatePortGroup(ctx context.Context, id int64, payload PortGroupPayload) (*PortGroupView, error) {
	var out *PortGroupView
	err := svc.store.runInTx(ctx, func(ctx context.Context, tx *store) error {
		if _, err := tx.getPortGroup(ctx, id); err != nil {
			return err
		}

		group, err := portGroupFromPayload(ctx, tx, id, payload)
		if err != nil {
			return err
		}

		group.ID = id
		group.UpdatedAt = time.Now().UTC()
		if err := tx.updatePortGroup(ctx, id, group); err != nil {
			return err
		}
		if err := replacePortGroupChildren(ctx, tx, id, payload, group.UpdatedAt); err != nil {
			return err
		}

		view, err := tx.getPortGroupView(ctx, id)
		if err != nil {
			return err
		}
		out = view
		return nil
	})
	return out, err
}

func (svc *Service) DeletePortGroup(ctx context.Context, id int64) error {
	return svc.store.runInTx(ctx, func(ctx context.Context, tx *store) error {
		return tx.deletePortGroup(ctx, id)
	})
}

func (svc *Service) UpdatePortGroups(ctx context.Context, input PortGroupBatchUpdateInput) ([]PortGroupView, error) {
	normalized, err := normalizeBatchUpdateInput(input)
	if err != nil {
		return nil, err
	}

	var out []PortGroupView
	err = svc.store.runInTx(ctx, func(ctx context.Context, tx *store) error {
		count, err := tx.countPortGroupsByIDs(ctx, normalized.IDs)
		if err != nil {
			return err
		}
		if count != len(normalized.IDs) {
			return huma.Error404NotFound("one or more port groups were not found")
		}

		if err := tx.batchUpdatePortGroups(ctx, normalized, time.Now().UTC()); err != nil {
			return err
		}

		groups, err := tx.listPortGroupsByIDs(ctx, normalized.IDs)
		if err != nil {
			return err
		}
		views, err := tx.buildPortGroupViews(ctx, groups)
		if err != nil {
			return err
		}
		out = views
		return nil
	})
	return out, err
}

func (svc *Service) DeletePortGroups(ctx context.Context, input PortGroupBatchDeleteInput) error {
	normalized, err := normalizeBatchDeleteInput(input)
	if err != nil {
		return err
	}

	return svc.store.runInTx(ctx, func(ctx context.Context, tx *store) error {
		count, err := tx.countPortGroupsByIDs(ctx, normalized.IDs)
		if err != nil {
			return err
		}
		if count != len(normalized.IDs) {
			return huma.Error404NotFound("one or more port groups were not found")
		}
		return tx.deletePortGroups(ctx, normalized.IDs)
	})
}

func (svc *Service) ExportPortGroupsCSV(ctx context.Context, params PortGroupListParams) ([]byte, error) {
	groups, err := svc.ListPortGroups(ctx, params)
	if err != nil {
		return nil, err
	}

	records := [][]string{{
		"host_ip",
		"host_name",
		"environment",
		"service_name",
		"status",
		"port_start",
		"port_end",
		"container_name",
		"dind_host",
		"owner",
		"tags",
		"components",
		"repositories",
		"slots",
		"notes",
	}}

	for _, group := range groups {
		hostIP := ""
		hostName := ""
		hostEnvironment := ""
		if group.Host != nil {
			hostIP = group.Host.IP
			hostName = group.Host.Name
			hostEnvironment = group.Host.Environment
		}

		records = append(records, []string{
			hostIP,
			hostName,
			hostEnvironment,
			group.ServiceName,
			group.Status,
			strconv.Itoa(group.PortStart),
			strconv.Itoa(group.PortEnd),
			group.ContainerName,
			group.DindHost,
			group.Owner,
			group.Tags,
			formatPortGroupComponents(group.Components),
			formatPortGroupRepositories(group.Repositories),
			formatPortGroupSlots(group.Slots),
			group.Notes,
		})
	}

	return csvBytes(records)
}
