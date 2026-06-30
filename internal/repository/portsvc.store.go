package repository

import (
	"context"
	"strings"

	"github.com/lwmacct/260630-go-hsr-shared/pkg/idgen"
	"github.com/uptrace/bun"
)

func (s *Store) ListPortsvcServices(ctx context.Context, params PortsvcServiceListFilter) ([]PortsvcServiceRecord, error) {
	var rows []ServicesModel
	query := s.db.NewSelect().Model(&rows).Relation("PortAllocation")
	if !params.Admin {
		query = query.Where("service.owner_subject = ?", params.OwnerSubject)
	} else if params.OwnerSubject != "" {
		query = query.Where("service.owner_subject = ?", params.OwnerSubject)
	}
	if status := utilCompactString(params.Status); status != "" {
		query = query.Where("service.status = ?", status)
	}
	if projectName := utilCompactString(params.ProjectName); projectName != "" {
		query = query.Where("LOWER(service.project_name) LIKE ?", utilSearchPattern(projectName))
	}
	if keyword := utilCompactString(params.Query); keyword != "" {
		pattern := utilSearchPattern(keyword)
		query = query.Where(
			utilJoinSearchClauses([]string{
				"service.name",
				"service.project_name",
				"service.dind_ip",
				"service.dind_container",
				"service.owner",
				"service.tags",
				"service.notes",
			}),
			utilJoinSearchArgs(pattern, 7)...,
		)
	}
	switch strings.ToLower(params.Sort) {
	case "name":
		query = query.Order("service.name ASC", "service.id ASC")
	case "status":
		query = query.Order("service.status ASC", "service.name ASC")
	case "project":
		query = query.Order("service.project_name ASC", "service.name ASC")
	case "updated_desc":
		query = query.Order("service.updated_at DESC", "service.name ASC")
	default:
		query = query.Order("service.name ASC", "service.id ASC")
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]PortsvcServiceRecord, 0, len(rows))
	for idx := range rows {
		out = append(out, *utilPortsvcServiceRecordFromModel(&rows[idx]))
	}
	return out, nil
}

func (s *Store) FetchPortsvcServiceByID(ctx context.Context, id string) (*PortsvcServiceRecord, error) {
	row := new(ServicesModel)
	err := s.db.NewSelect().Model(row).Relation("PortAllocation").Where("service.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcServiceRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcService(ctx context.Context, service *PortsvcServiceRecord) (*PortsvcServiceRecord, error) {
	row := &ServicesModel{
		ID:               idgen.NewUUID7(),
		OwnerSubject:     service.OwnerSubject,
		PortAllocationID: service.PortAllocationID,
		Name:             service.Name,
		ProjectName:      service.ProjectName,
		DindIP:           service.DindIP,
		DindContainer:    service.DindContainer,
		Status:           service.Status,
		Owner:            service.Owner,
		Tags:             service.Tags,
		Notes:            service.Notes,
		CreatedAt:        service.CreatedAt,
		UpdatedAt:        service.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return s.FetchPortsvcServiceByID(ctx, row.ID)
}

func (s *Store) UpdatePortsvcService(ctx context.Context, id string, service *PortsvcServiceRecord) (*PortsvcServiceRecord, error) {
	row := &ServicesModel{
		ID:               id,
		OwnerSubject:     service.OwnerSubject,
		PortAllocationID: service.PortAllocationID,
		Name:             service.Name,
		ProjectName:      service.ProjectName,
		DindIP:           service.DindIP,
		DindContainer:    service.DindContainer,
		Status:           service.Status,
		Owner:            service.Owner,
		Tags:             service.Tags,
		Notes:            service.Notes,
		UpdatedAt:        service.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("owner_subject", "port_allocation_id", "name", "project_name", "dind_ip", "dind_container", "status", "owner", "tags", "notes", "updated_at").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, ErrNotFound
	}
	return s.FetchPortsvcServiceByID(ctx, id)
}

func (s *Store) DeletePortsvcServices(ctx context.Context, ids []string, ownerSubject string, admin bool) error {
	query := s.db.NewDelete().Model((*ServicesModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("owner_subject = ?", ownerSubject)
	}
	_, err := query.Exec(ctx)
	return err
}

func (s *Store) CountPortsvcServicesByIDs(ctx context.Context, ids []string, ownerSubject string, admin bool) (int, error) {
	query := s.db.NewSelect().Model((*ServicesModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("owner_subject = ?", ownerSubject)
	}
	return query.Count(ctx)
}

func (s *Store) ListPortsvcServiceChildrenByServiceIDs(ctx context.Context, ids []string) ([]PortsvcServiceChildrenRecord, error) {
	children := make(map[string]PortsvcServiceChildrenRecord, len(ids))
	for _, id := range ids {
		children[id] = PortsvcServiceChildrenRecord{ServiceID: id}
	}

	var repoRows []struct {
		bun.BaseModel `bun:"table:service_repositories,alias:sr"`
		RepositoriesModel

		ServiceID string `bun:"service_id"`
	}
	if err := s.db.NewSelect().
		Model(&repoRows).
		Column("sr.service_id").
		ColumnExpr("repository_ref.*").
		Join("JOIN repositories AS repository_ref ON repository_ref.id = sr.repository_id").
		Where("sr.service_id IN (?)", bun.List(ids)).
		Order("repository_ref.name ASC").
		Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range repoRows {
		record := *utilPortsvcRepositoryRecordFromModel(&repoRows[idx].RepositoriesModel)
		child := children[repoRows[idx].ServiceID]
		child.Repositories = append(child.Repositories, record)
		children[repoRows[idx].ServiceID] = child
	}

	var depRows []struct {
		bun.BaseModel `bun:"table:service_dependencies,alias:sd"`
		DependenciesModel

		ServiceID string `bun:"service_id"`
	}
	if err := s.db.NewSelect().
		Model(&depRows).
		Column("sd.service_id").
		ColumnExpr("dependency.*").
		Join("JOIN dependencies AS dependency ON dependency.id = sd.dependency_id").
		Where("sd.service_id IN (?)", bun.List(ids)).
		Order("dependency.name ASC").
		Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range depRows {
		record := *utilPortsvcDependencyRecordFromModel(&depRows[idx].DependenciesModel)
		child := children[depRows[idx].ServiceID]
		child.Dependencies = append(child.Dependencies, record)
		children[depRows[idx].ServiceID] = child
	}

	out := make([]PortsvcServiceChildrenRecord, 0, len(children))
	for _, id := range ids {
		out = append(out, children[id])
	}
	return out, nil
}

func (s *Store) ReplacePortsvcServiceChildren(ctx context.Context, serviceID string) error {
	if _, err := s.db.NewDelete().Model((*ServiceRepositoriesModel)(nil)).Where("service_id = ?", serviceID).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*ServiceDependenciesModel)(nil)).Where("service_id = ?", serviceID).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Store) ListPortsvcPortAllocations(ctx context.Context, params PortsvcPortAllocationListFilter) ([]PortsvcPortAllocationRecord, error) {
	var rows []PortAllocationsModel
	query := s.db.NewSelect().Model(&rows)
	if !params.Admin {
		query = query.Where("allocation.owner_subject = ?", params.OwnerSubject)
	} else if params.OwnerSubject != "" {
		query = query.Where("allocation.owner_subject = ?", params.OwnerSubject)
	}
	if status := utilCompactString(params.Status); status != "" {
		query = query.Where("allocation.status = ?", status)
	}
	switch strings.ToLower(params.Sort) {
	case "status":
		query = query.Order("allocation.status ASC", "allocation.port_start ASC")
	case "updated_desc":
		query = query.Order("allocation.updated_at DESC", "allocation.port_start ASC")
	default:
		query = query.Order("allocation.port_start ASC")
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]PortsvcPortAllocationRecord, 0, len(rows))
	for idx := range rows {
		out = append(out, *utilPortsvcPortAllocationRecordFromModel(&rows[idx]))
	}
	return out, nil
}

func (s *Store) FetchPortsvcPortAllocationByID(ctx context.Context, id string) (*PortsvcPortAllocationRecord, error) {
	row := new(PortAllocationsModel)
	err := s.db.NewSelect().Model(row).Where("allocation.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcPortAllocationRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcPortAllocation(ctx context.Context, group *PortsvcPortAllocationRecord) (*PortsvcPortAllocationRecord, error) {
	row := &PortAllocationsModel{
		ID:           idgen.NewUUID7(),
		OwnerSubject: group.OwnerSubject,
		PortStart:    group.PortStart,
		PortEnd:      group.PortEnd,
		Status:       group.Status,
		Notes:        group.Notes,
		CreatedAt:    group.CreatedAt,
		UpdatedAt:    group.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return s.FetchPortsvcPortAllocationByID(ctx, row.ID)
}

func (s *Store) UpdatePortsvcPortAllocation(ctx context.Context, id string, group *PortsvcPortAllocationRecord) (*PortsvcPortAllocationRecord, error) {
	row := &PortAllocationsModel{
		ID:           id,
		OwnerSubject: group.OwnerSubject,
		PortStart:    group.PortStart,
		PortEnd:      group.PortEnd,
		Status:       group.Status,
		Notes:        group.Notes,
		UpdatedAt:    group.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("owner_subject", "port_start", "port_end", "status", "notes", "updated_at").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, ErrNotFound
	}
	return s.FetchPortsvcPortAllocationByID(ctx, id)
}

func (s *Store) DeletePortsvcPortAllocations(ctx context.Context, ids []string, ownerSubject string, admin bool) error {
	query := s.db.NewDelete().Model((*PortAllocationsModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("owner_subject = ?", ownerSubject)
	}
	_, err := query.Exec(ctx)
	return err
}

func (s *Store) CountPortsvcPortAllocationsByIDs(ctx context.Context, ids []string, ownerSubject string, admin bool) (int, error) {
	query := s.db.NewSelect().Model((*PortAllocationsModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("owner_subject = ?", ownerSubject)
	}
	return query.Count(ctx)
}

func (s *Store) CountPortsvcOverlappingPortAllocations(ctx context.Context, currentID string, group *PortsvcPortAllocationRecord) (int, error) {
	query := s.db.NewSelect().
		Model((*PortAllocationsModel)(nil)).
		Where("owner_subject = ?", group.OwnerSubject).
		Where("port_start = ?", group.PortStart)
	if currentID != "" {
		query = query.Where("id != ?", currentID)
	}
	return query.Count(ctx)
}

func (s *Store) ListPortsvcPortAllocationStartsByOwner(ctx context.Context, ownerSubject string, excludeID string) ([]int, error) {
	var rows []struct {
		PortStart int `bun:"port_start"`
	}
	query := s.db.NewSelect().Model((*PortAllocationsModel)(nil)).Column("port_start").Where("owner_subject = ?", ownerSubject)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Scan(ctx, &rows); err != nil {
		return nil, err
	}
	out := make([]int, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.PortStart)
	}
	return out, nil
}

func (s *Store) FetchPortsvcRepositoryByOwnerAndURL(ctx context.Context, ownerSubject string, url string) (*PortsvcRepositoryRecord, error) {
	row := new(RepositoriesModel)
	err := s.db.NewSelect().Model(row).Where("owner_subject = ?", ownerSubject).Where("url = ?", strings.TrimSpace(url)).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcRepositoryRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcRepository(ctx context.Context, repo *PortsvcRepositoryRecord) (*PortsvcRepositoryRecord, error) {
	row := &RepositoriesModel{
		ID:           idgen.NewUUID7(),
		OwnerSubject: repo.OwnerSubject,
		Name:         repo.Name,
		URL:          repo.URL,
		Kind:         repo.Kind,
		Notes:        repo.Notes,
		CreatedAt:    repo.CreatedAt,
		UpdatedAt:    repo.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return utilPortsvcRepositoryRecordFromModel(row), nil
}

func (s *Store) FetchPortsvcDependencyByIdentity(ctx context.Context, ownerSubject string, name string, itemType string, version string) (*PortsvcDependencyRecord, error) {
	row := new(DependenciesModel)
	err := s.db.NewSelect().Model(row).
		Where("owner_subject = ?", ownerSubject).
		Where("name = ?", strings.TrimSpace(name)).
		Where("type = ?", strings.TrimSpace(itemType)).
		Where("version = ?", strings.TrimSpace(version)).
		Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcDependencyRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcDependency(ctx context.Context, component *PortsvcDependencyRecord) (*PortsvcDependencyRecord, error) {
	row := &DependenciesModel{
		ID:           idgen.NewUUID7(),
		OwnerSubject: component.OwnerSubject,
		Name:         component.Name,
		Type:         component.Type,
		URL:          component.URL,
		Version:      component.Version,
		Notes:        component.Notes,
		CreatedAt:    component.CreatedAt,
		UpdatedAt:    component.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return utilPortsvcDependencyRecordFromModel(row), nil
}

func (s *Store) AddPortsvcServiceRepositories(ctx context.Context, links []PortsvcServiceRepositoryRecord) error {
	if len(links) == 0 {
		return nil
	}
	rows := make([]ServiceRepositoriesModel, 0, len(links))
	for _, link := range links {
		rows = append(rows, ServiceRepositoriesModel{
			ID:           idgen.NewUUID7(),
			ServiceID:    link.ServiceID,
			RepositoryID: link.RepositoryID,
			Role:         link.Role,
			CreatedAt:    link.CreatedAt,
			UpdatedAt:    link.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
}

func (s *Store) AddPortsvcServiceDependencies(ctx context.Context, links []PortsvcServiceDependencyRecord) error {
	if len(links) == 0 {
		return nil
	}
	rows := make([]ServiceDependenciesModel, 0, len(links))
	for _, link := range links {
		rows = append(rows, ServiceDependenciesModel{
			ID:           idgen.NewUUID7(),
			ServiceID:    link.ServiceID,
			DependencyID: link.DependencyID,
			Role:         link.Role,
			Notes:        link.Notes,
			CreatedAt:    link.CreatedAt,
			UpdatedAt:    link.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
}
