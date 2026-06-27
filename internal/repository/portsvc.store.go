package repository

import (
	"context"
	"strings"

	"github.com/uptrace/bun"
)

func (s *Store) ListPortsvcServices(ctx context.Context, params PortsvcServiceListFilter) ([]PortsvcServiceRecord, error) {
	var rows []ServicesModel
	query := s.db.NewSelect().Model(&rows).Relation("User").Relation("PortAllocation")
	if !params.Admin {
		query = query.Where("service.user_id = ?", params.UserID)
	} else if params.UserID > 0 {
		query = query.Where("service.user_id = ?", params.UserID)
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

func (s *Store) FetchPortsvcServiceByID(ctx context.Context, id int64) (*PortsvcServiceRecord, error) {
	row := new(ServicesModel)
	err := s.db.NewSelect().Model(row).Relation("User").Relation("PortAllocation").Where("service.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcServiceRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcService(ctx context.Context, service *PortsvcServiceRecord) (*PortsvcServiceRecord, error) {
	row := &ServicesModel{
		UserID:           service.UserID,
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

func (s *Store) UpdatePortsvcService(ctx context.Context, id int64, service *PortsvcServiceRecord) (*PortsvcServiceRecord, error) {
	row := &ServicesModel{
		ID:               id,
		UserID:           service.UserID,
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
		Column("user_id", "port_allocation_id", "name", "project_name", "dind_ip", "dind_container", "status", "owner", "tags", "notes", "updated_at").
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

func (s *Store) DeletePortsvcServices(ctx context.Context, ids []int64, userID int64, admin bool) error {
	query := s.db.NewDelete().Model((*ServicesModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("user_id = ?", userID)
	}
	_, err := query.Exec(ctx)
	return err
}

func (s *Store) CountPortsvcServicesByIDs(ctx context.Context, ids []int64, userID int64, admin bool) (int, error) {
	query := s.db.NewSelect().Model((*ServicesModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("user_id = ?", userID)
	}
	return query.Count(ctx)
}

func (s *Store) ListPortsvcServiceChildrenByServiceIDs(ctx context.Context, ids []int64) ([]PortsvcServiceChildrenRecord, error) {
	children := make(map[int64]PortsvcServiceChildrenRecord, len(ids))
	for _, id := range ids {
		children[id] = PortsvcServiceChildrenRecord{ServiceID: id}
	}

	var repoRows []struct {
		bun.BaseModel `bun:"table:service_repositories,alias:sr"`
		RepositoriesModel

		ServiceID int64 `bun:"service_id"`
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

		ServiceID int64 `bun:"service_id"`
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

func (s *Store) ReplacePortsvcServiceChildren(ctx context.Context, serviceID int64) error {
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
	query := s.db.NewSelect().Model(&rows).Relation("User")
	if !params.Admin {
		query = query.Where("allocation.user_id = ?", params.UserID)
	} else if params.UserID > 0 {
		query = query.Where("allocation.user_id = ?", params.UserID)
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

func (s *Store) FetchPortsvcPortAllocationByID(ctx context.Context, id int64) (*PortsvcPortAllocationRecord, error) {
	row := new(PortAllocationsModel)
	err := s.db.NewSelect().Model(row).Relation("User").Where("allocation.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcPortAllocationRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcPortAllocation(ctx context.Context, group *PortsvcPortAllocationRecord) (*PortsvcPortAllocationRecord, error) {
	row := &PortAllocationsModel{
		UserID:    group.UserID,
		PortStart: group.PortStart,
		PortEnd:   group.PortEnd,
		Status:    group.Status,
		Notes:     group.Notes,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return s.FetchPortsvcPortAllocationByID(ctx, row.ID)
}

func (s *Store) UpdatePortsvcPortAllocation(ctx context.Context, id int64, group *PortsvcPortAllocationRecord) (*PortsvcPortAllocationRecord, error) {
	row := &PortAllocationsModel{
		ID:        id,
		UserID:    group.UserID,
		PortStart: group.PortStart,
		PortEnd:   group.PortEnd,
		Status:    group.Status,
		Notes:     group.Notes,
		UpdatedAt: group.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("user_id", "port_start", "port_end", "status", "notes", "updated_at").
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

func (s *Store) DeletePortsvcPortAllocations(ctx context.Context, ids []int64, userID int64, admin bool) error {
	query := s.db.NewDelete().Model((*PortAllocationsModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("user_id = ?", userID)
	}
	_, err := query.Exec(ctx)
	return err
}

func (s *Store) CountPortsvcPortAllocationsByIDs(ctx context.Context, ids []int64, userID int64, admin bool) (int, error) {
	query := s.db.NewSelect().Model((*PortAllocationsModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("user_id = ?", userID)
	}
	return query.Count(ctx)
}

func (s *Store) CountPortsvcOverlappingPortAllocations(ctx context.Context, currentID int64, group *PortsvcPortAllocationRecord) (int, error) {
	query := s.db.NewSelect().
		Model((*PortAllocationsModel)(nil)).
		Where("user_id = ?", group.UserID).
		Where("port_start = ?", group.PortStart)
	if currentID > 0 {
		query = query.Where("id != ?", currentID)
	}
	return query.Count(ctx)
}

func (s *Store) ListPortsvcPortAllocationStartsByUser(ctx context.Context, userID int64, excludeID int64) ([]int, error) {
	var rows []struct {
		PortStart int `bun:"port_start"`
	}
	query := s.db.NewSelect().Model((*PortAllocationsModel)(nil)).Column("port_start").Where("user_id = ?", userID)
	if excludeID > 0 {
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

func (s *Store) FetchPortsvcRepositoryByUserAndURL(ctx context.Context, userID int64, url string) (*PortsvcRepositoryRecord, error) {
	row := new(RepositoriesModel)
	err := s.db.NewSelect().Model(row).Where("user_id = ?", userID).Where("url = ?", strings.TrimSpace(url)).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcRepositoryRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcRepository(ctx context.Context, repo *PortsvcRepositoryRecord) (*PortsvcRepositoryRecord, error) {
	row := &RepositoriesModel{
		UserID:    repo.UserID,
		Name:      repo.Name,
		URL:       repo.URL,
		Kind:      repo.Kind,
		Notes:     repo.Notes,
		CreatedAt: repo.CreatedAt,
		UpdatedAt: repo.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return utilPortsvcRepositoryRecordFromModel(row), nil
}

func (s *Store) FetchPortsvcDependencyByIdentity(ctx context.Context, userID int64, name string, itemType string, version string) (*PortsvcDependencyRecord, error) {
	row := new(DependenciesModel)
	err := s.db.NewSelect().Model(row).
		Where("user_id = ?", userID).
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
		UserID:    component.UserID,
		Name:      component.Name,
		Type:      component.Type,
		URL:       component.URL,
		Version:   component.Version,
		Notes:     component.Notes,
		CreatedAt: component.CreatedAt,
		UpdatedAt: component.UpdatedAt,
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
