package service

import (
	"context"
	"strconv"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

type InventoryService struct {
	store *repository.Store
}

func NewInventoryService(store *repository.Store) *InventoryService {
	if store == nil {
		panic("NewInventoryService: store is nil")
	}
	return &InventoryService{store: store}
}

func (s *InventoryService) ListHosts(ctx context.Context, params HostListParams) ([]InventoryHost, error) {
	return s.store.ListInventoryHosts(ctx, repository.InventoryHostListFilter(params))
}

func (s *InventoryService) CreateHost(ctx context.Context, payload HostPayload) (*InventoryHost, error) {
	host := utilInventoryHostFromPayload(payload)
	if err := validateHost(host); err != nil {
		return nil, err
	}
	now := utilNowUTC()
	host.CreatedAt = now
	host.UpdatedAt = now
	if _, err := s.store.CreateInventoryHost(ctx, host); err != nil {
		return nil, err
	}
	return host, nil
}

func (s *InventoryService) UpdateHost(ctx context.Context, id int64, payload HostPayload) (*InventoryHost, error) {
	host := utilInventoryHostFromPayload(payload)
	host.ID = id
	if err := validateHost(host); err != nil {
		return nil, err
	}
	host.UpdatedAt = utilNowUTC()
	out, err := s.store.UpdateInventoryHost(ctx, id, host)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *InventoryService) DeleteHost(ctx context.Context, id int64) error {
	count, err := s.store.CountInventoryPortGroupsByHostID(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return utilBadInventoryRequest("host still has port groups")
	}
	deleted, err := s.store.DeleteInventoryHost(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return utilInventoryNotFound("host not found")
	}
	return nil
}

func (s *InventoryService) ListPortGroups(ctx context.Context, params PortGroupListParams) ([]PortGroupView, error) {
	groups, err := s.store.ListInventoryPortGroups(ctx, repository.InventoryPortGroupListFilter(params))
	if err != nil {
		return nil, err
	}
	return s.buildPortGroupViews(ctx, groups)
}

func (s *InventoryService) GetPortGroup(ctx context.Context, id int64) (*PortGroupView, error) {
	group, err := s.store.FetchInventoryPortGroupWithHostByID(ctx, id)
	if err != nil {
		return nil, err
	}
	views, err := s.buildPortGroupViews(ctx, []InventoryPortGroup{*group})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (s *InventoryService) CreatePortGroup(ctx context.Context, payload PortGroupPayload) (*PortGroupView, error) {
	var out *PortGroupView
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := NewInventoryService(txStore)
		group, err := utilInventoryPortGroupFromPayload(ctx, tx.store, 0, payload)
		if err != nil {
			return err
		}
		now := utilNowUTC()
		group.CreatedAt = now
		group.UpdatedAt = now
		if _, err := tx.store.CreateInventoryPortGroup(ctx, group); err != nil {
			return err
		}
		if err := tx.replacePortGroupChildren(ctx, group.ID, payload, now); err != nil {
			return err
		}
		view, err := tx.GetPortGroup(ctx, group.ID)
		if err != nil {
			return err
		}
		out = view
		return nil
	})
	return out, err
}

func (s *InventoryService) UpdatePortGroup(ctx context.Context, id int64, payload PortGroupPayload) (*PortGroupView, error) {
	var out *PortGroupView
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := NewInventoryService(txStore)
		if _, err := tx.store.FetchInventoryPortGroupByID(ctx, id); err != nil {
			return err
		}
		group, err := utilInventoryPortGroupFromPayload(ctx, tx.store, id, payload)
		if err != nil {
			return err
		}
		group.ID = id
		group.UpdatedAt = utilNowUTC()
		updated, err := tx.store.UpdateInventoryPortGroup(ctx, id, group)
		if err != nil {
			return err
		}
		_ = updated
		if err := tx.replacePortGroupChildren(ctx, id, payload, group.UpdatedAt); err != nil {
			return err
		}
		view, err := tx.GetPortGroup(ctx, id)
		if err != nil {
			return err
		}
		out = view
		return nil
	})
	return out, err
}

func (s *InventoryService) DeletePortGroup(ctx context.Context, id int64) error {
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := NewInventoryService(txStore)
		count, err := tx.store.CountInventoryPortGroupsByIDs(ctx, []int64{id})
		if err != nil {
			return err
		}
		if count == 0 {
			return utilInventoryNotFound("port group not found")
		}
		return tx.store.DeleteInventoryPortGroups(ctx, []int64{id})
	})
}

func (s *InventoryService) UpdatePortGroups(ctx context.Context, input PortGroupBatchUpdateInput) ([]PortGroupView, error) {
	normalized, err := utilNormalizeBatchUpdateInput(input)
	if err != nil {
		return nil, err
	}
	var out []PortGroupView
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := NewInventoryService(txStore)
		count, err := tx.store.CountInventoryPortGroupsByIDs(ctx, normalized.IDs)
		if err != nil {
			return err
		}
		if count != len(normalized.IDs) {
			return utilInventoryNotFound("one or more port groups were not found")
		}
		if _, err := tx.store.UpdateInventoryPortGroupsBatch(ctx, normalized.IDs, normalized.Status, normalized.Owner, normalized.Tags, utilNowUTC()); err != nil {
			return err
		}
		groups, err := tx.store.ListInventoryPortGroupsByIDs(ctx, normalized.IDs)
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

func (s *InventoryService) DeletePortGroups(ctx context.Context, input PortGroupBatchDeleteInput) error {
	normalized, err := utilNormalizeBatchDeleteInput(input)
	if err != nil {
		return err
	}
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := NewInventoryService(txStore)
		count, err := tx.store.CountInventoryPortGroupsByIDs(ctx, normalized.IDs)
		if err != nil {
			return err
		}
		if count != len(normalized.IDs) {
			return utilInventoryNotFound("one or more port groups were not found")
		}
		return tx.store.DeleteInventoryPortGroups(ctx, normalized.IDs)
	})
}

func (s *InventoryService) ExportPortGroupsCSV(ctx context.Context, params PortGroupListParams) ([]byte, error) {
	groups, err := s.ListPortGroups(ctx, params)
	if err != nil {
		return nil, err
	}
	records := [][]string{{
		"host_ip", "host_name", "environment", "service_name", "status", "port_start", "port_end",
		"container_name", "dind_host", "owner", "tags", "components", "repositories", "slots", "notes",
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
			utilPortGroupComponents(group.Components),
			utilPortGroupRepositories(group.Repositories),
			utilPortGroupSlots(group.Slots),
			group.Notes,
		})
	}
	return utilCSVBytes(records)
}

func (s *InventoryService) buildPortGroupViews(ctx context.Context, groups []InventoryPortGroup) ([]PortGroupView, error) {
	views := make([]PortGroupView, len(groups))
	if len(groups) == 0 {
		return views, nil
	}
	ids := make([]int64, len(groups))
	for idx, group := range groups {
		ids[idx] = group.ID
		views[idx] = PortGroupView{InventoryPortGroup: group}
	}
	children, err := s.store.FetchInventoryPortGroupChildrenByPortGroupIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	viewByID := make(map[int64]*PortGroupView, len(views))
	for idx := range views {
		viewByID[views[idx].ID] = &views[idx]
	}
	for _, slot := range children.Slots {
		viewByID[slot.PortGroupID].Slots = append(viewByID[slot.PortGroupID].Slots, slot)
	}
	for _, component := range children.Components {
		viewByID[component.PortGroupID].Components = append(viewByID[component.PortGroupID].Components, component)
	}
	for _, repo := range children.Repositories {
		viewByID[repo.PortGroupID].Repositories = append(viewByID[repo.PortGroupID].Repositories, repo)
	}
	return views, nil
}
