package service

import (
	"context"
	"errors"
	"strings"
	"time"

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

func (s *PortsvcService) ListServices(ctx context.Context, params ServiceListParams) ([]ServiceView, error) {
	services, err := s.store.ListPortsvcServices(ctx, repository.PortsvcServiceListFilter{
		UserID:      utilVisibleUserID(params.Actor, params.UserID),
		Admin:       params.Actor.Admin,
		Query:       params.Query,
		Sort:        params.Sort,
		Status:      params.Status,
		ProjectName: params.ProjectName,
	})
	if err != nil {
		return nil, err
	}
	return s.buildServiceViews(ctx, services)
}

func (s *PortsvcService) GetService(ctx context.Context, actor PortsvcActor, id int64) (*ServiceView, error) {
	service, err := s.store.FetchPortsvcServiceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !actor.Admin && service.UserID != actor.UserID {
		return nil, utilPortsvcNotFound("service not found")
	}
	views, err := s.buildServiceViews(ctx, []PortsvcServiceRecord{*service})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (s *PortsvcService) CreateService(ctx context.Context, actor PortsvcActor, payload ServicePayload) (*ServiceView, error) {
	var out *ServiceView
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := NewPortsvcService(txStore)
		service, serviceErr := utilNormalizeService(ctx, tx.store, actor, nil, payload)
		if serviceErr != nil {
			return serviceErr
		}
		now := utilNowUTC()
		service.CreatedAt = now
		service.UpdatedAt = now
		created, createErr := tx.store.CreatePortsvcService(ctx, service)
		if createErr != nil {
			return createErr
		}
		if childErr := tx.replaceServiceChildren(ctx, created.ID, created.UserID, payload, now); childErr != nil {
			return childErr
		}
		view, viewErr := tx.GetService(ctx, PortsvcActor{UserID: created.UserID, Admin: true}, created.ID)
		if viewErr != nil {
			return viewErr
		}
		out = view
		return nil
	})
	return out, err
}

func (s *PortsvcService) UpdateService(ctx context.Context, actor PortsvcActor, id int64, payload ServicePayload) (*ServiceView, error) {
	var out *ServiceView
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		tx := NewPortsvcService(txStore)
		current, err := tx.store.FetchPortsvcServiceByID(ctx, id)
		if err != nil {
			return err
		}
		if !actor.Admin && current.UserID != actor.UserID {
			return utilPortsvcNotFound("service not found")
		}
		service, serviceErr := utilNormalizeService(ctx, tx.store, actor, current, payload)
		if serviceErr != nil {
			return serviceErr
		}
		service.UpdatedAt = utilNowUTC()
		updated, updateErr := tx.store.UpdatePortsvcService(ctx, id, service)
		if updateErr != nil {
			return updateErr
		}
		if childErr := tx.replaceServiceChildren(ctx, id, updated.UserID, payload, service.UpdatedAt); childErr != nil {
			return childErr
		}
		view, viewErr := tx.GetService(ctx, PortsvcActor{UserID: updated.UserID, Admin: true}, id)
		if viewErr != nil {
			return viewErr
		}
		out = view
		return nil
	})
	return out, err
}

func (s *PortsvcService) DeleteService(ctx context.Context, actor PortsvcActor, id int64) error {
	return s.DeleteServices(ctx, ServiceBatchDeleteInput{Actor: actor, IDs: []int64{id}})
}

func (s *PortsvcService) DeleteServices(ctx context.Context, input ServiceBatchDeleteInput) error {
	ids, err := utilNormalizeIDs(input.IDs)
	if err != nil {
		return err
	}
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		count, countErr := txStore.CountPortsvcServicesByIDs(ctx, ids, input.Actor.UserID, input.Actor.Admin)
		if countErr != nil {
			return countErr
		}
		if count != len(ids) {
			return utilPortsvcNotFound("one or more services were not found")
		}
		return txStore.DeletePortsvcServices(ctx, ids, input.Actor.UserID, input.Actor.Admin)
	})
}

func (s *PortsvcService) ListPortAllocations(ctx context.Context, params PortAllocationListParams) ([]PortAllocation, error) {
	return s.store.ListPortsvcPortAllocations(ctx, repository.PortsvcPortAllocationListFilter{
		UserID: utilVisibleUserID(params.Actor, params.UserID),
		Admin:  params.Actor.Admin,
		Sort:   params.Sort,
		Status: params.Status,
	})
}

func (s *PortsvcService) GetPortAllocation(ctx context.Context, actor PortsvcActor, id int64) (*PortAllocation, error) {
	group, err := s.store.FetchPortsvcPortAllocationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !actor.Admin && group.UserID != actor.UserID {
		return nil, utilPortsvcNotFound("port allocation not found")
	}
	return group, nil
}

func (s *PortsvcService) CreatePortAllocation(ctx context.Context, actor PortsvcActor, payload PortAllocationPayload) (*PortAllocation, error) {
	var out *PortAllocation
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		group, groupErr := utilNormalizePortAllocation(ctx, txStore, actor, nil, payload)
		if groupErr != nil {
			return groupErr
		}
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
	return out, err
}

func (s *PortsvcService) UpdatePortAllocation(ctx context.Context, actor PortsvcActor, id int64, payload PortAllocationPayload) (*PortAllocation, error) {
	var out *PortAllocation
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		current, err := txStore.FetchPortsvcPortAllocationByID(ctx, id)
		if err != nil {
			return err
		}
		if !actor.Admin && current.UserID != actor.UserID {
			return utilPortsvcNotFound("port allocation not found")
		}
		group, groupErr := utilNormalizePortAllocation(ctx, txStore, actor, current, payload)
		if groupErr != nil {
			return groupErr
		}
		group.UpdatedAt = utilNowUTC()
		updated, updateErr := txStore.UpdatePortsvcPortAllocation(ctx, id, group)
		if updateErr != nil {
			return updateErr
		}
		out = updated
		return nil
	})
	return out, err
}

func (s *PortsvcService) DeletePortAllocation(ctx context.Context, actor PortsvcActor, id int64) error {
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		count, err := txStore.CountPortsvcPortAllocationsByIDs(ctx, []int64{id}, actor.UserID, actor.Admin)
		if err != nil {
			return err
		}
		if count == 0 {
			return utilPortsvcNotFound("port allocation not found")
		}
		return txStore.DeletePortsvcPortAllocations(ctx, []int64{id}, actor.UserID, actor.Admin)
	})
}

func (s *PortsvcService) ExportServicesCSV(ctx context.Context, params ServiceListParams) ([]byte, error) {
	services, err := s.ListServices(ctx, params)
	if err != nil {
		return nil, err
	}
	records := [][]string{{
		"username", "name", "project_name", "status", "port_range", "dind_ip", "dind_container",
		"owner", "tags", "repositories", "dependencies", "notes",
	}}
	for _, service := range services {
		records = append(records, []string{
			service.Username,
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
	ids := make([]int64, len(services))
	for idx, service := range services {
		ids[idx] = service.ID
		views[idx] = ServiceView{PortsvcServiceRecord: service}
	}
	childrenList, err := s.store.ListPortsvcServiceChildrenByServiceIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	children := make(map[int64]repository.PortsvcServiceChildrenRecord, len(childrenList))
	for _, child := range childrenList {
		children[child.ServiceID] = child
	}
	for idx := range views {
		child := children[views[idx].ID]
		views[idx].Repositories = child.Repositories
		views[idx].Dependencies = child.Dependencies
	}
	return views, nil
}

func (s *PortsvcService) replaceServiceChildren(ctx context.Context, serviceID int64, userID int64, payload ServicePayload, now time.Time) error {
	if err := s.store.ReplacePortsvcServiceChildren(ctx, serviceID); err != nil {
		return err
	}
	repositoryLinks := make([]repository.PortsvcServiceRepositoryRecord, 0, len(payload.Repositories))
	for _, item := range payload.Repositories {
		repo, ok, err := s.ensureRepository(ctx, userID, item, now)
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
		dependency, ok, err := s.ensureDependency(ctx, userID, item, now)
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

func (s *PortsvcService) ensureRepository(ctx context.Context, userID int64, payload RepositoryPayload, now time.Time) (*RepositoryRef, bool, error) {
	if payload.ID > 0 {
		return &RepositoryRef{ID: payload.ID}, true, nil
	}
	name := strings.TrimSpace(payload.Name)
	url := strings.TrimSpace(payload.URL)
	if name == "" && url == "" {
		return nil, false, nil
	}
	if url != "" {
		existing, err := s.store.FetchPortsvcRepositoryByUserAndURL(ctx, userID, url)
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
		UserID:    userID,
		Name:      name,
		URL:       url,
		Kind:      kind,
		Notes:     strings.TrimSpace(payload.Notes),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, false, err
	}
	return created, true, nil
}

func (s *PortsvcService) ensureDependency(ctx context.Context, userID int64, payload DependencyPayload, now time.Time) (*Dependency, bool, error) {
	if payload.ID > 0 {
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
	existing, err := s.store.FetchPortsvcDependencyByIdentity(ctx, userID, name, itemType, version)
	if err == nil {
		return existing, true, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, false, err
	}
	created, err := s.store.CreatePortsvcDependency(ctx, &repository.PortsvcDependencyRecord{
		UserID:    userID,
		Name:      name,
		Type:      itemType,
		URL:       strings.TrimSpace(payload.URL),
		Version:   version,
		Notes:     strings.TrimSpace(payload.Notes),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, false, err
	}
	return created, true, nil
}
