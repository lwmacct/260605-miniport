package service

import (
	"context"
	"strings"
	"time"

	"github.com/lwmacct/260605-miniport/internal/repository"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/identity"
)

type PortsvcService struct {
	store     *repository.Store
	directory identity.Directory
}

func NewPortsvcService(store *repository.Store, directory identity.Directory) *PortsvcService {
	if store == nil {
		panic("NewPortsvcService: store is nil")
	}
	if directory == nil {
		panic("NewPortsvcService: directory is nil")
	}
	return &PortsvcService{store: store, directory: directory}
}

func (s *PortsvcService) ListHosts(ctx context.Context, params HostListParams) ([]Host, error) {
	return s.store.ListPortsvcHosts(ctx, repository.PortsvcHostListFilter{Query: params.Query, Status: params.Status})
}

func (s *PortsvcService) CreateHost(ctx context.Context, _ PortsvcActor, payload HostPayload) (*Host, error) {
	host, err := utilNormalizeHost(payload)
	if err != nil {
		return nil, err
	}
	now := utilNowUTC()
	host.CreatedAt = now
	host.UpdatedAt = now
	return s.store.CreatePortsvcHost(ctx, host)
}

func (s *PortsvcService) UpdateHost(ctx context.Context, _ PortsvcActor, id string, payload HostPayload) (*Host, error) {
	host, err := utilNormalizeHost(payload)
	if err != nil {
		return nil, err
	}
	host.UpdatedAt = utilNowUTC()
	out, err := s.store.UpdatePortsvcHost(ctx, id, host)
	if err != nil {
		return nil, utilWrapNotFound(err, "host not found")
	}
	return out, nil
}

func (s *PortsvcService) DeleteHost(ctx context.Context, _ PortsvcActor, id string) error {
	if err := s.store.DeletePortsvcHost(ctx, id); err != nil {
		return utilWrapNotFound(err, "host not found")
	}
	return nil
}

func (s *PortsvcService) ListPortGroups(ctx context.Context, params PortGroupListParams) ([]PortGroupView, error) {
	groups, err := s.store.ListPortsvcPortGroups(ctx, repository.PortsvcPortGroupListFilter{
		OwnerSubject: utilVisibleOwnerSubject(params.Actor, params.OwnerSubject),
		Admin:        params.Actor.Admin,
		Query:        params.Query,
		Sort:         params.Sort,
		Status:       params.Status,
	})
	if err != nil {
		return nil, err
	}
	return s.buildPortGroupViews(ctx, groups)
}

func (s *PortsvcService) GetPortGroup(ctx context.Context, actor PortsvcActor, id string) (*PortGroupView, error) {
	group, err := s.store.FetchPortsvcPortGroupByID(ctx, id)
	if err != nil {
		return nil, utilWrapNotFound(err, "port group not found")
	}
	if !actor.Admin && group.OwnerSubject != actor.OwnerSubject {
		return nil, utilPortsvcNotFound("port group not found")
	}
	views, err := s.buildPortGroupViews(ctx, []PortGroup{*group})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (s *PortsvcService) CreatePortGroup(ctx context.Context, actor PortsvcActor, payload PortGroupPayload) (*PortGroupView, error) {
	group, err := utilNormalizePortGroup(ctx, s.store, actor, nil, payload)
	if err != nil {
		return nil, err
	}
	if principalErr := utilResolveActivePrincipal(ctx, s.directory, group.OwnerSubject); principalErr != nil {
		return nil, principalErr
	}
	var createdID string
	now := utilNowUTC()
	group.CreatedAt = now
	group.UpdatedAt = now
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore, directory: s.directory}
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
		return nil, err
	}
	return s.GetPortGroup(ctx, PortsvcActor{OwnerSubject: group.OwnerSubject, Admin: true}, createdID)
}

func (s *PortsvcService) UpdatePortGroup(ctx context.Context, actor PortsvcActor, id string, payload PortGroupPayload) (*PortGroupView, error) {
	current, err := s.store.FetchPortsvcPortGroupByID(ctx, id)
	if err != nil {
		return nil, utilWrapNotFound(err, "port group not found")
	}
	if !actor.Admin && current.OwnerSubject != actor.OwnerSubject {
		return nil, utilPortsvcNotFound("port group not found")
	}
	group, err := utilNormalizePortGroup(ctx, s.store, actor, current, payload)
	if err != nil {
		return nil, err
	}
	if principalErr := utilResolveActivePrincipal(ctx, s.directory, group.OwnerSubject); principalErr != nil {
		return nil, principalErr
	}
	group.UpdatedAt = utilNowUTC()
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore, directory: s.directory}
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
		return nil, err
	}
	return s.GetPortGroup(ctx, PortsvcActor{OwnerSubject: group.OwnerSubject, Admin: true}, id)
}

func (s *PortsvcService) DeletePortGroup(ctx context.Context, actor PortsvcActor, id string) error {
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		count, err := txStore.CountPortsvcPortGroupsByIDs(ctx, []string{id}, actor.OwnerSubject, actor.Admin)
		if err != nil {
			return err
		}
		if count == 0 {
			return utilPortsvcNotFound("port group not found")
		}
		return txStore.DeletePortsvcPortGroups(ctx, []string{id}, actor.OwnerSubject, actor.Admin)
	})
}

func (s *PortsvcService) CreatePortSlot(ctx context.Context, actor PortsvcActor, groupID string, payload PortSlotPayload) (*PortSlot, error) {
	group, err := s.store.FetchPortsvcPortGroupByID(ctx, groupID)
	if err != nil {
		return nil, utilWrapNotFound(err, "port group not found")
	}
	if !actor.Admin && group.OwnerSubject != actor.OwnerSubject {
		return nil, utilPortsvcNotFound("port group not found")
	}
	slot, err := utilNormalizePortSlot(*group, payload)
	if err != nil {
		return nil, err
	}
	now := utilNowUTC()
	slot.CreatedAt = now
	slot.UpdatedAt = now
	return s.store.CreatePortsvcPortSlot(ctx, slot)
}

func (s *PortsvcService) UpdatePortSlot(ctx context.Context, actor PortsvcActor, id string, payload PortSlotPayload) (*PortSlot, error) {
	current, err := s.store.FetchPortsvcPortSlotByID(ctx, id)
	if err != nil {
		return nil, utilWrapNotFound(err, "port slot not found")
	}
	group, err := s.store.FetchPortsvcPortGroupByID(ctx, current.PortGroupID)
	if err != nil {
		return nil, err
	}
	if !actor.Admin && group.OwnerSubject != actor.OwnerSubject {
		return nil, utilPortsvcNotFound("port slot not found")
	}
	slot, err := utilNormalizePortSlot(*group, payload)
	if err != nil {
		return nil, err
	}
	slot.PortGroupID = current.PortGroupID
	slot.UpdatedAt = utilNowUTC()
	return s.store.UpdatePortsvcPortSlot(ctx, id, slot)
}

func (s *PortsvcService) DeletePortSlot(ctx context.Context, actor PortsvcActor, id string) error {
	current, err := s.store.FetchPortsvcPortSlotByID(ctx, id)
	if err != nil {
		return utilWrapNotFound(err, "port slot not found")
	}
	group, err := s.store.FetchPortsvcPortGroupByID(ctx, current.PortGroupID)
	if err != nil {
		return err
	}
	if !actor.Admin && group.OwnerSubject != actor.OwnerSubject {
		return utilPortsvcNotFound("port slot not found")
	}
	if err := s.store.DeletePortsvcPortSlot(ctx, id); err != nil {
		return utilWrapNotFound(err, "port slot not found")
	}
	return nil
}

func (s *PortsvcService) ExportPortGroupsCSV(ctx context.Context, params PortGroupListParams) ([]byte, error) {
	groups, err := s.ListPortGroups(ctx, params)
	if err != nil {
		return nil, err
	}
	records := [][]string{{
		"ownerName", "project_name", "status", "port_range", "runtime_mode", "runtime_name", "service_ip",
		"host_name", "host_ip", "project_owner", "tags", "slots", "repositories", "dependencies", "notes",
	}}
	for _, group := range groups {
		hostName := ""
		hostIP := ""
		if group.Host != nil {
			hostName = group.Host.Name
			hostIP = group.Host.IP
		}
		records = append(records, []string{
			group.OwnerName,
			group.ProjectName,
			group.Status,
			utilPortRange(group.PortGroup),
			group.RuntimeMode,
			group.RuntimeName,
			group.ServiceIP,
			hostName,
			hostIP,
			group.ProjectOwner,
			group.Tags,
			utilPortSlots(group.Slots),
			utilRepositories(group.Repositories),
			utilDependencies(group.Dependencies),
			group.Notes,
		})
	}
	return utilCSVBytes(records)
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
		views[idx].Repositories = child.Repositories
		views[idx].Dependencies = child.Dependencies
	}
	if err := s.attachPortGroupOwnerNames(ctx, views); err != nil {
		return nil, err
	}
	return views, nil
}

func (s *PortsvcService) replacePortGroupChildren(ctx context.Context, group PortGroup, payload PortGroupPayload, now time.Time) error {
	if err := s.store.ReplacePortsvcPortGroupChildren(ctx, group.ID); err != nil {
		return err
	}
	slots, err := utilNormalizePortSlots(group, payload.Slots)
	if err != nil {
		return err
	}
	for idx := range slots {
		slots[idx].CreatedAt = now
		slots[idx].UpdatedAt = now
	}
	if err := s.store.AddPortsvcPortSlots(ctx, slots); err != nil {
		return err
	}

	repos := make([]repository.PortsvcRepositoryRecord, 0, len(payload.Repositories))
	for _, item := range payload.Repositories {
		name := strings.TrimSpace(item.Name)
		url := strings.TrimSpace(item.URL)
		if name == "" && url == "" {
			continue
		}
		if url == "" {
			return utilBadPortsvcRequest("repository url is required")
		}
		kind := strings.TrimSpace(item.Kind)
		if kind == "" {
			kind = defaultRepositoryKind
		}
		repos = append(repos, repository.PortsvcRepositoryRecord{
			OwnerSubject: group.OwnerSubject,
			PortGroupID:  group.ID,
			Name:         name,
			URL:          url,
			Kind:         kind,
			Notes:        strings.TrimSpace(item.Notes),
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}
	if err := s.store.AddPortsvcRepositories(ctx, repos); err != nil {
		return err
	}

	deps := make([]repository.PortsvcDependencyRecord, 0, len(payload.Dependencies))
	for _, item := range payload.Dependencies {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		itemType := strings.TrimSpace(item.Type)
		if itemType == "" {
			itemType = defaultDependencyType
		}
		deps = append(deps, repository.PortsvcDependencyRecord{
			OwnerSubject: group.OwnerSubject,
			PortGroupID:  group.ID,
			Name:         name,
			Type:         itemType,
			URL:          strings.TrimSpace(item.URL),
			Version:      strings.TrimSpace(item.Version),
			Notes:        strings.TrimSpace(item.Notes),
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}
	return s.store.AddPortsvcDependencies(ctx, deps)
}

func (s *PortsvcService) attachPortGroupOwnerNames(ctx context.Context, views []PortGroupView) error {
	subjects := make([]string, 0, len(views))
	for _, view := range views {
		subjects = append(subjects, view.OwnerSubject)
	}
	names, err := s.ownerNames(ctx, subjects)
	if err != nil {
		return err
	}
	for idx := range views {
		views[idx].OwnerName = names[views[idx].OwnerSubject]
		if views[idx].OwnerName == "" {
			views[idx].OwnerName = views[idx].OwnerSubject
		}
	}
	return nil
}

func (s *PortsvcService) ownerNames(ctx context.Context, subjects []string) (map[string]string, error) {
	out := make(map[string]string, len(subjects))
	unique := make([]string, 0, len(subjects))
	for _, subject := range subjects {
		subject = strings.TrimSpace(subject)
		if subject == "" {
			continue
		}
		if _, ok := out[subject]; ok {
			continue
		}
		out[subject] = subject
		unique = append(unique, subject)
	}
	principals, err := s.directory.Principals(ctx, unique)
	if err != nil {
		return nil, err
	}
	for subject, principal := range principals {
		out[subject] = utilPrincipalName(principal, subject)
	}
	return out, nil
}
