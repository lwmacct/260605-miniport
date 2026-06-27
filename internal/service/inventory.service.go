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

func (s *InventoryService) ListPortGroups(ctx context.Context, params PortGroupListParams) ([]PortGroupView, error) {
	groups, err := s.store.ListInventoryPortGroups(ctx, repository.InventoryPortGroupListFilter{
		UserID:      utilVisibleUserID(params.Actor, params.UserID),
		Admin:       params.Actor.Admin,
		Query:       params.Query,
		Sort:        params.Sort,
		Status:      params.Status,
		ProjectName: params.ProjectName,
		DindIP:      params.DindIP,
	})
	if err != nil {
		return nil, err
	}
	return s.buildPortGroupViews(ctx, groups)
}

func (s *InventoryService) GetPortGroup(ctx context.Context, actor InventoryActor, id int64) (*PortGroupView, error) {
	group, err := s.store.FetchInventoryPortGroupWithHostByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !actor.Admin && group.UserID != actor.UserID {
		return nil, utilInventoryNotFound("allocation not found")
	}
	views, err := s.buildPortGroupViews(ctx, []InventoryPortGroup{*group})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (s *InventoryService) CreatePortGroup(ctx context.Context, actor InventoryActor, payload PortGroupPayload) (*PortGroupView, error) {
	var out *PortGroupView
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := NewInventoryService(txStore)
		group, groupErr := utilInventoryPortGroupFromPayload(ctx, tx.store, actor, nil, payload)
		if groupErr != nil {
			return groupErr
		}
		now := utilNowUTC()
		group.CreatedAt = now
		group.UpdatedAt = now
		created, createErr := tx.store.CreateInventoryPortGroup(ctx, group)
		if createErr != nil {
			return createErr
		}
		group.ID = created.ID
		group.Username = created.Username
		if replaceErr := utilReplacePortGroupChildren(ctx, tx.store, group.ID, payload, now, *group); replaceErr != nil {
			return replaceErr
		}
		view, viewErr := tx.GetPortGroup(ctx, InventoryActor{UserID: group.UserID, Admin: true}, group.ID)
		if viewErr != nil {
			return viewErr
		}
		out = view
		return nil
	})
	return out, err
}

func (s *InventoryService) UpdatePortGroup(ctx context.Context, actor InventoryActor, id int64, payload PortGroupPayload) (*PortGroupView, error) {
	var out *PortGroupView
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := NewInventoryService(txStore)
		current, err := tx.store.FetchInventoryPortGroupByID(ctx, id)
		if err != nil {
			return err
		}
		if !actor.Admin && current.UserID != actor.UserID {
			return utilInventoryNotFound("allocation not found")
		}
		group, groupErr := utilInventoryPortGroupFromPayload(ctx, tx.store, actor, current, payload)
		if groupErr != nil {
			return groupErr
		}
		group.ID = id
		group.UpdatedAt = utilNowUTC()
		if _, updateErr := tx.store.UpdateInventoryPortGroup(ctx, id, group); updateErr != nil {
			return updateErr
		}
		if replaceErr := utilReplacePortGroupChildren(ctx, tx.store, id, payload, group.UpdatedAt, *group); replaceErr != nil {
			return replaceErr
		}
		view, viewErr := tx.GetPortGroup(ctx, InventoryActor{UserID: group.UserID, Admin: true}, id)
		if viewErr != nil {
			return viewErr
		}
		out = view
		return nil
	})
	return out, err
}

func (s *InventoryService) DeletePortGroup(ctx context.Context, actor InventoryActor, id int64) error {
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := NewInventoryService(txStore)
		count, err := tx.store.CountInventoryPortGroupsByIDs(ctx, []int64{id}, actor.UserID, actor.Admin)
		if err != nil {
			return err
		}
		if count == 0 {
			return utilInventoryNotFound("allocation not found")
		}
		return tx.store.DeleteInventoryPortGroups(ctx, []int64{id}, actor.UserID, actor.Admin)
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
		count, countErr := tx.store.CountInventoryPortGroupsByIDs(ctx, normalized.IDs, normalized.Actor.UserID, normalized.Actor.Admin)
		if countErr != nil {
			return countErr
		}
		if count != len(normalized.IDs) {
			return utilInventoryNotFound("one or more allocations were not found")
		}
		if _, updateErr := tx.store.UpdateInventoryPortGroupsBatch(ctx, normalized.IDs, normalized.Status, normalized.Owner, normalized.Tags, utilNowUTC(), normalized.Actor.UserID, normalized.Actor.Admin); updateErr != nil {
			return updateErr
		}
		groups, listErr := tx.store.ListInventoryPortGroupsByIDs(ctx, normalized.IDs, normalized.Actor.UserID, normalized.Actor.Admin)
		if listErr != nil {
			return listErr
		}
		views, buildErr := tx.buildPortGroupViews(ctx, groups)
		if buildErr != nil {
			return buildErr
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
		count, err := tx.store.CountInventoryPortGroupsByIDs(ctx, normalized.IDs, normalized.Actor.UserID, normalized.Actor.Admin)
		if err != nil {
			return err
		}
		if count != len(normalized.IDs) {
			return utilInventoryNotFound("one or more allocations were not found")
		}
		return tx.store.DeleteInventoryPortGroups(ctx, normalized.IDs, normalized.Actor.UserID, normalized.Actor.Admin)
	})
}

func (s *InventoryService) ExportPortGroupsCSV(ctx context.Context, params PortGroupListParams) ([]byte, error) {
	groups, err := s.ListPortGroups(ctx, params)
	if err != nil {
		return nil, err
	}
	records := [][]string{{
		"username", "name", "status", "port_start", "port_end", "dind_ip", "dind_container",
		"owner", "tags", "projects", "dependencies", "repositories", "ports", "notes",
	}}
	for _, group := range groups {
		records = append(records, []string{
			group.Username,
			group.Name,
			group.Status,
			strconv.Itoa(group.PortStart),
			strconv.Itoa(group.PortEnd),
			group.DindIP,
			group.DindContainer,
			group.Owner,
			group.Tags,
			utilPortGroupProjects(group.Projects),
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
	childrenList, err := s.store.ListInventoryPortGroupChildrenByPortGroupIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	children := make(map[int64]repository.InventoryPortGroupChildrenRecord, len(childrenList))
	for _, child := range childrenList {
		children[child.PortGroupID] = child
	}
	for idx := range views {
		child := children[views[idx].ID]
		views[idx].Slots = child.Slots
		views[idx].Projects = child.Projects
		views[idx].Components = child.Components
		views[idx].Repositories = child.Repositories
	}
	return views, nil
}
