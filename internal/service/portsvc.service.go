package service

import (
	"context"
	"errors"
	"time"

	"github.com/lwmacct/260630-go-hsr-shared/pkg/idgen"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

type PortsvcService struct {
	store *repository.Store
}

func NewPortsvcService(store *repository.Store) *PortsvcService {
	if store == nil {
		panic("NewPortsvcService: store is nil")
	}
	return &PortsvcService{store: store}
}

func (s *PortsvcService) CreateHosts(ctx context.Context, payloads []HostPayload) ([]Host, error) {
	items := make([]Host, 0, len(payloads))
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, payload := range payloads {
			item, err := tx.CreateHost(ctx, payload)
			if err != nil {
				return err
			}
			items = append(items, *item)
		}
		return nil
	})
	return items, err
}

func (s *PortsvcService) UpdateHosts(ctx context.Context, inputs []HostUpdateInput) ([]Host, error) {
	items := make([]Host, 0, len(inputs))
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, input := range inputs {
			item, err := tx.UpdateHost(ctx, input.ID, input.Payload)
			if err != nil {
				return err
			}
			items = append(items, *item)
		}
		return nil
	})
	return items, err
}

func (s *PortsvcService) DeleteHosts(ctx context.Context, ids []string) error {
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, id := range ids {
			if err := tx.DeleteHost(ctx, id); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *PortsvcService) ListHosts(ctx context.Context, params HostListParams) ([]Host, error) {
	return s.store.ListPortsvcHosts(ctx, repository.PortsvcHostListFilter{Query: params.Query, Status: params.Status})
}

func (s *PortsvcService) CreateHost(ctx context.Context, payload HostPayload) (*Host, error) {
	host, err := utilNormalizeHost(payload)
	if err != nil {
		return nil, err
	}
	now := utilNowUTC()
	host.CreatedAt = now
	host.UpdatedAt = now
	created, err := s.store.CreatePortsvcHost(ctx, host)
	if err != nil {
		return nil, utilWrapConflict(err, "host name already exists")
	}
	return created, nil
}

func (s *PortsvcService) UpdateHost(ctx context.Context, id string, payload HostPayload) (*Host, error) {
	host, err := utilNormalizeHost(payload)
	if err != nil {
		return nil, err
	}
	host.UpdatedAt = utilNowUTC()
	out, err := s.store.UpdatePortsvcHost(ctx, id, host)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, utilPortsvcNotFound("host not found")
		}
		return nil, utilWrapConflict(err, "host name already exists")
	}
	return out, nil
}

func (s *PortsvcService) DeleteHost(ctx context.Context, id string) error {
	if err := s.store.DeletePortsvcHost(ctx, id); err != nil {
		return utilWrapNotFound(err, "host not found")
	}
	return nil
}

func (s *PortsvcService) ListDependencyAssets(ctx context.Context, params DependencyAssetListParams) ([]DependencyAsset, error) {
	return s.store.ListPortsvcDependencyAssets(ctx, repository.PortsvcDependencyAssetListFilter{
		Query: params.Query, AssetKind: params.AssetKind, AssetType: params.AssetType,
		Provider: params.Provider, Status: params.Status,
	})
}

func (s *PortsvcService) CreateDependencyAssets(ctx context.Context, payloads []DependencyAssetPayload) ([]DependencyAsset, error) {
	items := make([]DependencyAsset, 0, len(payloads))
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, payload := range payloads {
			item, err := tx.CreateDependencyAsset(ctx, payload)
			if err != nil {
				return err
			}
			items = append(items, *item)
		}
		return nil
	})
	return items, err
}

func (s *PortsvcService) UpdateDependencyAssets(ctx context.Context, inputs []DependencyAssetUpdateInput) ([]DependencyAsset, error) {
	items := make([]DependencyAsset, 0, len(inputs))
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, input := range inputs {
			item, err := tx.UpdateDependencyAsset(ctx, input.ID, input.Payload)
			if err != nil {
				return err
			}
			items = append(items, *item)
		}
		return nil
	})
	return items, err
}

func (s *PortsvcService) DeleteDependencyAssets(ctx context.Context, ids []string) error {
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, id := range ids {
			if err := tx.DeleteDependencyAsset(ctx, id); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *PortsvcService) CreateDependencyAsset(ctx context.Context, payload DependencyAssetPayload) (*DependencyAsset, error) {
	asset, err := utilNormalizeDependencyAsset(payload)
	if err != nil {
		return nil, err
	}
	now := utilNowUTC()
	asset.CreatedAt = now
	asset.UpdatedAt = now
	created, err := s.store.CreatePortsvcDependencyAsset(ctx, asset)
	if err != nil {
		return nil, utilWrapConflict(err, "dependency asset name and kind already exist")
	}
	return created, nil
}

func (s *PortsvcService) UpdateDependencyAsset(ctx context.Context, id string, payload DependencyAssetPayload) (*DependencyAsset, error) {
	_, err := s.store.FetchPortsvcDependencyAssetByID(ctx, id)
	if err != nil {
		return nil, utilWrapNotFound(err, "dependency asset not found")
	}
	asset, err := utilNormalizeDependencyAsset(payload)
	if err != nil {
		return nil, err
	}
	asset.UpdatedAt = utilNowUTC()
	updated, err := s.store.UpdatePortsvcDependencyAsset(ctx, id, asset)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, utilPortsvcNotFound("dependency asset not found")
		}
		return nil, utilWrapConflict(err, "dependency asset name and kind already exist")
	}
	return updated, nil
}

func (s *PortsvcService) DeleteDependencyAsset(ctx context.Context, id string) error {
	if err := s.store.DeletePortsvcDependencyAsset(ctx, id); err != nil {
		return utilWrapNotFound(err, "dependency asset not found")
	}
	return nil
}

func (s *PortsvcService) ListServiceGroups(ctx context.Context, params ServiceGroupListParams) ([]ServiceGroupView, error) {
	groups, err := s.store.ListPortsvcServiceGroups(ctx, repository.PortsvcServiceGroupListFilter{
		Query: params.Query, Status: params.Status,
	})
	if err != nil {
		return nil, err
	}
	return s.buildServiceGroupViews(ctx, groups)
}

func (s *PortsvcService) CreateServiceGroups(ctx context.Context, payloads []ServiceGroupPayload) ([]ServiceGroupView, error) {
	ids := make([]string, 0, len(payloads))
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, payload := range payloads {
			item, err := tx.CreateServiceGroup(ctx, payload)
			if err != nil {
				return err
			}
			ids = append(ids, item.ID)
		}
		return nil
	})
	if err != nil {
		return nil, utilWrapConflict(err, "service group name already exists")
	}
	return s.serviceGroupViewsByIDs(ctx, ids)
}

func (s *PortsvcService) UpdateServiceGroups(ctx context.Context, inputs []ServiceGroupUpdateInput) ([]ServiceGroupView, error) {
	ids := make([]string, 0, len(inputs))
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, input := range inputs {
			item, err := tx.UpdateServiceGroup(ctx, input.ID, input.Payload)
			if err != nil {
				return err
			}
			ids = append(ids, item.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.serviceGroupViewsByIDs(ctx, ids)
}

func (s *PortsvcService) DeleteServiceGroups(ctx context.Context, ids []string) error {
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, id := range ids {
			if err := tx.DeleteServiceGroup(ctx, id); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *PortsvcService) GetServiceGroup(ctx context.Context, id string) (*ServiceGroupView, error) {
	group, err := s.store.FetchPortsvcServiceGroupByID(ctx, id)
	if err != nil {
		return nil, utilWrapNotFound(err, "service group not found")
	}
	views, err := s.buildServiceGroupViews(ctx, []ServiceGroup{*group})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (s *PortsvcService) CreateServiceGroup(ctx context.Context, payload ServiceGroupPayload) (*ServiceGroupView, error) {
	group, err := utilNormalizeServiceGroup(payload)
	if err != nil {
		return nil, err
	}
	var createdID string
	now := utilNowUTC()
	group.CreatedAt = now
	group.UpdatedAt = now
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		created, createErr := tx.store.CreatePortsvcServiceGroup(ctx, group)
		if createErr != nil {
			return createErr
		}
		createdID = created.ID
		if childErr := tx.replaceServiceGroupPortGroups(ctx, *created, payload, now); childErr != nil {
			return childErr
		}
		return nil
	})
	if err != nil {
		return nil, utilWrapConflict(err, "service group name already exists")
	}
	return s.GetServiceGroup(ctx, createdID)
}

func (s *PortsvcService) UpdateServiceGroup(ctx context.Context, id string, payload ServiceGroupPayload) (*ServiceGroupView, error) {
	_, err := s.store.FetchPortsvcServiceGroupByID(ctx, id)
	if err != nil {
		return nil, utilWrapNotFound(err, "service group not found")
	}
	group, err := utilNormalizeServiceGroup(payload)
	if err != nil {
		return nil, err
	}
	group.UpdatedAt = utilNowUTC()
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		updated, updateErr := tx.store.UpdatePortsvcServiceGroup(ctx, id, group)
		if updateErr != nil {
			return updateErr
		}
		if childErr := tx.replaceServiceGroupPortGroups(ctx, *updated, payload, group.UpdatedAt); childErr != nil {
			return childErr
		}
		return nil
	})
	if err != nil {
		return nil, utilWrapConflict(err, "service group name already exists")
	}
	return s.GetServiceGroup(ctx, id)
}

func (s *PortsvcService) DeleteServiceGroup(ctx context.Context, id string) error {
	if err := s.store.DeletePortsvcServiceGroup(ctx, id); err != nil {
		return utilWrapNotFound(err, "service group not found")
	}
	return nil
}

func (s *PortsvcService) ListPortGroups(ctx context.Context, params PortGroupListParams) ([]PortGroupView, error) {
	groups, err := s.store.ListPortsvcPortGroups(ctx, repository.PortsvcPortGroupListFilter{
		Query: params.Query, Sort: params.Sort, Status: params.Status,
	})
	if err != nil {
		return nil, err
	}
	return s.buildPortGroupViews(ctx, groups)
}

func (s *PortsvcService) CreatePortGroups(ctx context.Context, payloads []PortGroupPayload) ([]PortGroupView, error) {
	ids := make([]string, 0, len(payloads))
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, payload := range payloads {
			item, err := tx.CreatePortGroup(ctx, payload)
			if err != nil {
				return err
			}
			ids = append(ids, item.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.portGroupViewsByIDs(ctx, ids)
}

func (s *PortsvcService) UpdatePortGroups(ctx context.Context, inputs []PortGroupUpdateInput) ([]PortGroupView, error) {
	ids := make([]string, 0, len(inputs))
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, input := range inputs {
			item, err := tx.UpdatePortGroup(ctx, input.ID, input.Payload)
			if err != nil {
				return err
			}
			ids = append(ids, item.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.portGroupViewsByIDs(ctx, ids)
}

func (s *PortsvcService) DeletePortGroups(ctx context.Context, ids []string) error {
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		for _, id := range ids {
			if err := tx.DeletePortGroup(ctx, id); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *PortsvcService) GetPortGroup(ctx context.Context, id string) (*PortGroupView, error) {
	group, err := s.store.FetchPortsvcPortGroupByID(ctx, id)
	if err != nil {
		return nil, utilWrapNotFound(err, "port group not found")
	}
	views, err := s.buildPortGroupViews(ctx, []PortGroup{*group})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (s *PortsvcService) CreatePortGroup(ctx context.Context, payload PortGroupPayload) (*PortGroupView, error) {
	group, err := utilNormalizePortGroup(ctx, s.store, nil, payload)
	if err != nil {
		return nil, err
	}
	var createdID string
	now := utilNowUTC()
	group.CreatedAt = now
	group.UpdatedAt = now
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		created, createErr := tx.store.CreatePortsvcPortGroup(ctx, group)
		if createErr != nil {
			return createErr
		}
		createdID = created.ID
		if childErr := tx.replacePortGroupChildren(ctx, *created, payload, now); childErr != nil {
			return childErr
		}
		return nil
	})
	if err != nil {
		return nil, utilWrapConflict(err, "port group is already allocated")
	}
	return s.GetPortGroup(ctx, createdID)
}

func (s *PortsvcService) UpdatePortGroup(ctx context.Context, id string, payload PortGroupPayload) (*PortGroupView, error) {
	current, err := s.store.FetchPortsvcPortGroupByID(ctx, id)
	if err != nil {
		return nil, utilWrapNotFound(err, "port group not found")
	}
	group, err := utilNormalizePortGroup(ctx, s.store, current, payload)
	if err != nil {
		return nil, err
	}
	group.UpdatedAt = utilNowUTC()
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore}
		updated, updateErr := tx.store.UpdatePortsvcPortGroup(ctx, id, group)
		if updateErr != nil {
			return updateErr
		}
		if childErr := tx.replacePortGroupChildren(ctx, *updated, payload, group.UpdatedAt); childErr != nil {
			return childErr
		}
		return nil
	})
	if err != nil {
		return nil, utilWrapConflict(err, "port group is already allocated")
	}
	return s.GetPortGroup(ctx, id)
}

func (s *PortsvcService) DeletePortGroup(ctx context.Context, id string) error {
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		count, err := txStore.CountPortsvcPortGroupsByIDs(ctx, []string{id})
		if err != nil {
			return err
		}
		if count == 0 {
			return utilPortsvcNotFound("port group not found")
		}
		return txStore.DeletePortsvcPortGroups(ctx, []string{id})
	})
}

func (s *PortsvcService) ExportPortGroupsCSV(ctx context.Context, params PortGroupListParams) ([]byte, error) {
	groups, err := s.ListPortGroups(ctx, params)
	if err != nil {
		return nil, err
	}
	records := [][]string{{
		"environment_name", "status", "port_prefix", "runtime_mode", "runtime_name", "service_ip",
		"host_name", "host_ip", "environment_owner", "tags", "slots", "asset_links", "repository_links", "notes",
	}}
	for _, group := range groups {
		hostName := ""
		hostIP := ""
		if group.Host != nil {
			hostName = group.Host.Name
			hostIP = group.Host.IP
		}
		records = append(records, []string{
			group.EnvironmentName,
			group.Status,
			utilPortGroupLabel(group.PortGroup),
			group.RuntimeMode,
			group.RuntimeName,
			group.ServiceIP,
			hostName,
			hostIP,
			group.EnvironmentOwner,
			group.Tags,
			utilPortSlots(group.Slots),
			utilAssetLinks(group.AssetLinks),
			utilRepositoryLinks(group.RepositoryLinks),
			group.Notes,
		})
	}
	return utilCSVBytes(records)
}

func (s *PortsvcService) buildServiceGroupViews(ctx context.Context, groups []ServiceGroup) ([]ServiceGroupView, error) {
	views := make([]ServiceGroupView, len(groups))
	if len(groups) == 0 {
		return views, nil
	}
	ids := make([]string, len(groups))
	for idx, group := range groups {
		ids[idx] = group.ID
		views[idx] = ServiceGroupView{ServiceGroup: group}
	}
	childrenList, err := s.store.ListPortsvcServiceGroupChildrenByGroupIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	children := make(map[string]repository.PortsvcServiceGroupChildrenRecord, len(childrenList))
	for _, child := range childrenList {
		children[child.ServiceGroupID] = child
	}
	for idx := range views {
		child := children[views[idx].ID]
		views[idx].PortGroups = child.PortGroups
	}
	return views, nil
}

func (s *PortsvcService) serviceGroupViewsByIDs(ctx context.Context, ids []string) ([]ServiceGroupView, error) {
	groups := make([]ServiceGroup, 0, len(ids))
	for _, id := range ids {
		group, err := s.store.FetchPortsvcServiceGroupByID(ctx, id)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *group)
	}
	return s.buildServiceGroupViews(ctx, groups)
}

func (s *PortsvcService) replaceServiceGroupPortGroups(ctx context.Context, group ServiceGroup, payload ServiceGroupPayload, now time.Time) error {
	portGroups, err := utilNormalizeServiceGroupPortGroups(ctx, s.store, group, payload.PortGroups)
	if err != nil {
		return err
	}
	if err := s.store.ReplacePortsvcServiceGroupPortGroups(ctx, group.ID); err != nil {
		return err
	}
	for idx := range portGroups {
		portGroups[idx].CreatedAt = now
		portGroups[idx].UpdatedAt = now
	}
	return s.store.AddPortsvcServiceGroupPortGroups(ctx, portGroups)
}

func (s *PortsvcService) buildPortGroupViews(ctx context.Context, groups []PortGroup) ([]PortGroupView, error) {
	views := make([]PortGroupView, len(groups))
	if len(groups) == 0 {
		return views, nil
	}
	ids := make([]string, len(groups))
	for idx, group := range groups {
		ids[idx] = group.ID
		views[idx] = PortGroupView{PortGroup: group}
	}
	childrenList, err := s.store.ListPortsvcPortGroupChildrenByGroupIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	children := make(map[string]repository.PortsvcPortGroupChildrenRecord, len(childrenList))
	for _, child := range childrenList {
		children[child.PortGroupID] = child
	}
	for idx := range views {
		child := children[views[idx].ID]
		views[idx].Slots = child.Slots
		views[idx].AssetLinks = child.AssetLinks
		views[idx].RepositoryLinks = child.RepositoryLinks
	}
	return views, nil
}

func (s *PortsvcService) portGroupViewsByIDs(ctx context.Context, ids []string) ([]PortGroupView, error) {
	groups := make([]PortGroup, 0, len(ids))
	for _, id := range ids {
		group, err := s.store.FetchPortsvcPortGroupByID(ctx, id)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *group)
	}
	return s.buildPortGroupViews(ctx, groups)
}

func (s *PortsvcService) replacePortGroupChildren(ctx context.Context, group PortGroup, payload PortGroupPayload, now time.Time) error {
	slots, err := utilNormalizePortSlots(group, payload.Slots)
	if err != nil {
		return err
	}
	links, err := utilNormalizeAssetLinks(ctx, s.store, group, payload.AssetLinks)
	if err != nil {
		return err
	}
	repositoryLinks, err := utilNormalizeRepositoryLinks(ctx, s.store, group, payload.RepositoryLinks)
	if err != nil {
		return err
	}
	if err := s.store.ReplacePortsvcPortGroupChildren(ctx, group.ID); err != nil {
		return err
	}
	slotIDs := make(map[string]string, len(slots))
	for idx := range slots {
		originalID := slots[idx].ID
		if !repository.IsUUID7(originalID) {
			slots[idx].ID = idgen.NewUUID7()
		}
		if originalID != "" {
			slotIDs[originalID] = slots[idx].ID
		}
		slots[idx].CreatedAt = now
		slots[idx].UpdatedAt = now
	}
	if err := s.store.AddPortsvcPortSlots(ctx, slots); err != nil {
		return err
	}
	for idx := range links {
		if linkID, ok := slotIDs[links[idx].PortSlotID]; ok {
			links[idx].PortSlotID = linkID
		}
		links[idx].CreatedAt = now
		links[idx].UpdatedAt = now
	}
	if err := s.store.AddPortsvcPortGroupAssetLinks(ctx, links); err != nil {
		return err
	}
	for idx := range repositoryLinks {
		repositoryLinks[idx].CreatedAt = now
		repositoryLinks[idx].UpdatedAt = now
	}
	return s.store.AddPortsvcPortGroupRepositoryLinks(ctx, repositoryLinks)
}
