package service

import (
	"context"
	"errors"
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

func (s *PortsvcService) ListServices(ctx context.Context, params ServiceListParams) ([]ServiceView, error) {
	services, err := s.store.ListPortsvcServices(ctx, repository.PortsvcServiceListFilter{
		OwnerSubject: utilVisibleOwnerSubject(params.Actor, params.OwnerSubject),
		Admin:        params.Actor.Admin,
		Query:        params.Query,
		Sort:         params.Sort,
		Status:       params.Status,
		ProjectName:  params.ProjectName,
	})
	if err != nil {
		return nil, err
	}
	return s.buildServiceViews(ctx, services)
}

func (s *PortsvcService) GetService(ctx context.Context, actor PortsvcActor, id string) (*ServiceView, error) {
	service, err := s.store.FetchPortsvcServiceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !actor.Admin && service.OwnerSubject != actor.OwnerSubject {
		return nil, utilPortsvcNotFound("service not found")
	}
	views, err := s.buildServiceViews(ctx, []PortsvcServiceRecord{*service})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (s *PortsvcService) CreateService(ctx context.Context, actor PortsvcActor, payload ServicePayload) (*ServiceView, error) {
	service, err := utilNormalizeService(ctx, s.store, actor, nil, payload)
	if err != nil {
		return nil, err
	}
	if principalErr := utilResolveActivePrincipal(ctx, s.directory, service.OwnerSubject); principalErr != nil {
		return nil, principalErr
	}
	var createdID string
	now := utilNowUTC()
	service.CreatedAt = now
	service.UpdatedAt = now
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore, directory: s.directory}
		created, createErr := tx.store.CreatePortsvcService(ctx, service)
		if createErr != nil {
			return createErr
		}
		createdID = created.ID
		if childErr := tx.replaceServiceChildren(ctx, created.ID, created.OwnerSubject, payload, now); childErr != nil {
			return childErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetService(ctx, PortsvcActor{OwnerSubject: service.OwnerSubject, Admin: true}, createdID)
}

func (s *PortsvcService) UpdateService(ctx context.Context, actor PortsvcActor, id string, payload ServicePayload) (*ServiceView, error) {
	current, err := s.store.FetchPortsvcServiceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !actor.Admin && current.OwnerSubject != actor.OwnerSubject {
		return nil, utilPortsvcNotFound("service not found")
	}
	service, err := utilNormalizeService(ctx, s.store, actor, current, payload)
	if err != nil {
		return nil, err
	}
	if principalErr := utilResolveActivePrincipal(ctx, s.directory, service.OwnerSubject); principalErr != nil {
		return nil, principalErr
	}
	service.UpdatedAt = utilNowUTC()
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := &PortsvcService{store: txStore, directory: s.directory}
		updated, updateErr := tx.store.UpdatePortsvcService(ctx, id, service)
		if updateErr != nil {
			return updateErr
		}
		if childErr := tx.replaceServiceChildren(ctx, id, updated.OwnerSubject, payload, service.UpdatedAt); childErr != nil {
			return childErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetService(ctx, PortsvcActor{OwnerSubject: service.OwnerSubject, Admin: true}, id)
}

func (s *PortsvcService) DeleteService(ctx context.Context, actor PortsvcActor, id string) error {
	return s.DeleteServices(ctx, ServiceBatchDeleteInput{Actor: actor, IDs: []string{id}})
}

func (s *PortsvcService) DeleteServices(ctx context.Context, input ServiceBatchDeleteInput) error {
	ids, err := utilNormalizeIDs(input.IDs)
	if err != nil {
		return err
	}
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		count, countErr := txStore.CountPortsvcServicesByIDs(ctx, ids, input.Actor.OwnerSubject, input.Actor.Admin)
		if countErr != nil {
			return countErr
		}
		if count != len(ids) {
			return utilPortsvcNotFound("one or more services were not found")
		}
		return txStore.DeletePortsvcServices(ctx, ids, input.Actor.OwnerSubject, input.Actor.Admin)
	})
}

func (s *PortsvcService) ListPortAllocations(ctx context.Context, params PortAllocationListParams) ([]PortAllocation, error) {
	groups, err := s.store.ListPortsvcPortAllocations(ctx, repository.PortsvcPortAllocationListFilter{
		OwnerSubject: utilVisibleOwnerSubject(params.Actor, params.OwnerSubject),
		Admin:        params.Actor.Admin,
		Sort:         params.Sort,
		Status:       params.Status,
	})
	if err != nil {
		return nil, err
	}
	if err := s.attachPortAllocationOwnerNames(ctx, groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *PortsvcService) GetPortAllocation(ctx context.Context, actor PortsvcActor, id string) (*PortAllocation, error) {
	group, err := s.store.FetchPortsvcPortAllocationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !actor.Admin && group.OwnerSubject != actor.OwnerSubject {
		return nil, utilPortsvcNotFound("port allocation not found")
	}
	groups := []PortAllocation{*group}
	if err := s.attachPortAllocationOwnerNames(ctx, groups); err != nil {
		return nil, err
	}
	return &groups[0], nil
}

func (s *PortsvcService) CreatePortAllocation(ctx context.Context, actor PortsvcActor, payload PortAllocationPayload) (*PortAllocation, error) {
	group, err := utilNormalizePortAllocation(ctx, s.store, actor, nil, payload)
	if err != nil {
		return nil, err
	}
	if principalErr := utilResolveActivePrincipal(ctx, s.directory, group.OwnerSubject); principalErr != nil {
		return nil, principalErr
	}
	var out *PortAllocation
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		now := utilNowUTC()
		group.CreatedAt = now
		group.UpdatedAt = now
		created, createErr := txStore.CreatePortsvcPortAllocation(ctx, group)
		if createErr != nil {
			return createErr
		}
		out = created
		return nil
	})
	if err != nil {
		return nil, err
	}
	groups := []PortAllocation{*out}
	if err := s.attachPortAllocationOwnerNames(ctx, groups); err != nil {
		return nil, err
	}
	return &groups[0], nil
}

func (s *PortsvcService) UpdatePortAllocation(ctx context.Context, actor PortsvcActor, id string, payload PortAllocationPayload) (*PortAllocation, error) {
	current, err := s.store.FetchPortsvcPortAllocationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !actor.Admin && current.OwnerSubject != actor.OwnerSubject {
		return nil, utilPortsvcNotFound("port allocation not found")
	}
	group, err := utilNormalizePortAllocation(ctx, s.store, actor, current, payload)
	if err != nil {
		return nil, err
	}
	if principalErr := utilResolveActivePrincipal(ctx, s.directory, group.OwnerSubject); principalErr != nil {
		return nil, principalErr
	}
	var out *PortAllocation
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		group.UpdatedAt = utilNowUTC()
		updated, updateErr := txStore.UpdatePortsvcPortAllocation(ctx, id, group)
		if updateErr != nil {
			return updateErr
		}
		out = updated
		return nil
	})
	if err != nil {
		return nil, err
	}
	groups := []PortAllocation{*out}
	if err := s.attachPortAllocationOwnerNames(ctx, groups); err != nil {
		return nil, err
	}
	return &groups[0], nil
}

func (s *PortsvcService) DeletePortAllocation(ctx context.Context, actor PortsvcActor, id string) error {
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		count, err := txStore.CountPortsvcPortAllocationsByIDs(ctx, []string{id}, actor.OwnerSubject, actor.Admin)
		if err != nil {
			return err
		}
		if count == 0 {
			return utilPortsvcNotFound("port allocation not found")
		}
		return txStore.DeletePortsvcPortAllocations(ctx, []string{id}, actor.OwnerSubject, actor.Admin)
	})
}

func (s *PortsvcService) ExportServicesCSV(ctx context.Context, params ServiceListParams) ([]byte, error) {
	services, err := s.ListServices(ctx, params)
	if err != nil {
		return nil, err
	}
	records := [][]string{{
		"ownerName", "name", "project_name", "status", "port_range", "dind_ip", "dind_container",
		"owner", "tags", "repositories", "dependencies", "notes",
	}}
	for _, service := range services {
		records = append(records, []string{
			service.OwnerName,
			service.Name,
			service.ProjectName,
			service.Status,
			utilPortRange(service.PortAllocation),
			service.DindIP,
			service.DindContainer,
			service.Owner,
			service.Tags,
			utilServiceRepositories(service.Repositories),
			utilServiceDependencies(service.Dependencies),
			service.Notes,
		})
	}
	return utilCSVBytes(records)
}

func (s *PortsvcService) buildServiceViews(ctx context.Context, services []PortsvcServiceRecord) ([]ServiceView, error) {
	views := make([]ServiceView, len(services))
	if len(services) == 0 {
		return views, nil
	}
	ids := make([]string, len(services))
	for idx, service := range services {
		ids[idx] = service.ID
		views[idx] = ServiceView{PortsvcServiceRecord: service}
	}
	childrenList, err := s.store.ListPortsvcServiceChildrenByServiceIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	children := make(map[string]repository.PortsvcServiceChildrenRecord, len(childrenList))
	for _, child := range childrenList {
		children[child.ServiceID] = child
	}
	for idx := range views {
		child := children[views[idx].ID]
		views[idx].Repositories = child.Repositories
		views[idx].Dependencies = child.Dependencies
	}
	if err := s.attachServiceOwnerNames(ctx, views); err != nil {
		return nil, err
	}
	return views, nil
}

func (s *PortsvcService) replaceServiceChildren(ctx context.Context, serviceID string, ownerSubject string, payload ServicePayload, now time.Time) error {
	if err := s.store.ReplacePortsvcServiceChildren(ctx, serviceID); err != nil {
		return err
	}
	repositoryLinks := make([]repository.PortsvcServiceRepositoryRecord, 0, len(payload.Repositories))
	for _, item := range payload.Repositories {
		repo, ok, err := s.ensureRepository(ctx, ownerSubject, item, now)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		repositoryLinks = append(repositoryLinks, repository.PortsvcServiceRepositoryRecord{
			ServiceID:    serviceID,
			RepositoryID: repo.ID,
			Role:         strings.TrimSpace(item.Role),
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}
	if err := s.store.AddPortsvcServiceRepositories(ctx, repositoryLinks); err != nil {
		return err
	}

	dependencyLinks := make([]repository.PortsvcServiceDependencyRecord, 0, len(payload.Dependencies))
	for _, item := range payload.Dependencies {
		dependency, ok, err := s.ensureDependency(ctx, ownerSubject, item, now)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		dependencyLinks = append(dependencyLinks, repository.PortsvcServiceDependencyRecord{
			ServiceID:    serviceID,
			DependencyID: dependency.ID,
			Role:         strings.TrimSpace(item.Role),
			Notes:        strings.TrimSpace(item.Notes),
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}
	return s.store.AddPortsvcServiceDependencies(ctx, dependencyLinks)
}

func (s *PortsvcService) ensureRepository(ctx context.Context, ownerSubject string, payload RepositoryPayload, now time.Time) (*RepositoryRef, bool, error) {
	if payload.ID != "" {
		return &RepositoryRef{ID: payload.ID}, true, nil
	}
	name := strings.TrimSpace(payload.Name)
	url := strings.TrimSpace(payload.URL)
	if name == "" && url == "" {
		return nil, false, nil
	}
	if url != "" {
		existing, err := s.store.FetchPortsvcRepositoryByOwnerAndURL(ctx, ownerSubject, url)
		if err == nil {
			return existing, true, nil
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, false, err
		}
	}
	kind := strings.TrimSpace(payload.Kind)
	if kind == "" {
		kind = defaultRepositoryKind
	}
	created, err := s.store.CreatePortsvcRepository(ctx, &repository.PortsvcRepositoryRecord{
		OwnerSubject: ownerSubject,
		Name:         name,
		URL:          url,
		Kind:         kind,
		Notes:        strings.TrimSpace(payload.Notes),
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, false, err
	}
	return created, true, nil
}

func (s *PortsvcService) attachServiceOwnerNames(ctx context.Context, views []ServiceView) error {
	subjects := make([]string, 0, len(views)*2)
	for _, view := range views {
		subjects = append(subjects, view.OwnerSubject)
		if view.PortAllocation != nil {
			subjects = append(subjects, view.PortAllocation.OwnerSubject)
		}
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
		if views[idx].PortAllocation != nil {
			views[idx].PortAllocation.OwnerName = names[views[idx].PortAllocation.OwnerSubject]
			if views[idx].PortAllocation.OwnerName == "" {
				views[idx].PortAllocation.OwnerName = views[idx].PortAllocation.OwnerSubject
			}
		}
	}
	return nil
}

func (s *PortsvcService) attachPortAllocationOwnerNames(ctx context.Context, groups []PortAllocation) error {
	subjects := make([]string, 0, len(groups))
	for _, group := range groups {
		subjects = append(subjects, group.OwnerSubject)
	}
	names, err := s.ownerNames(ctx, subjects)
	if err != nil {
		return err
	}
	for idx := range groups {
		groups[idx].OwnerName = names[groups[idx].OwnerSubject]
		if groups[idx].OwnerName == "" {
			groups[idx].OwnerName = groups[idx].OwnerSubject
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

func (s *PortsvcService) ensureDependency(ctx context.Context, ownerSubject string, payload DependencyPayload, now time.Time) (*Dependency, bool, error) {
	if payload.ID != "" {
		return &Dependency{ID: payload.ID}, true, nil
	}
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return nil, false, nil
	}
	itemType := strings.TrimSpace(payload.Type)
	if itemType == "" {
		itemType = defaultDependencyType
	}
	version := strings.TrimSpace(payload.Version)
	existing, err := s.store.FetchPortsvcDependencyByIdentity(ctx, ownerSubject, name, itemType, version)
	if err == nil {
		return existing, true, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, false, err
	}
	created, err := s.store.CreatePortsvcDependency(ctx, &repository.PortsvcDependencyRecord{
		OwnerSubject: ownerSubject,
		Name:         name,
		Type:         itemType,
		URL:          strings.TrimSpace(payload.URL),
		Version:      version,
		Notes:        strings.TrimSpace(payload.Notes),
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, false, err
	}
	return created, true, nil
}
