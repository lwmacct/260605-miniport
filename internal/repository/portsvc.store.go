package repository

import (
	"context"
	"strings"

	"github.com/lwmacct/260630-go-hsr-shared/pkg/idgen"
	"github.com/uptrace/bun"
)

func (s *Store) ListPortsvcHosts(ctx context.Context, params PortsvcHostListFilter) ([]PortsvcHostRecord, error) {
	var rows []HostsModel
	query := s.db.NewSelect().Model(&rows)
	if status := utilCompactString(params.Status); status != "" {
		query = query.Where("host.status = ?", status)
	}
	if keyword := utilCompactString(params.Query); keyword != "" {
		pattern := utilSearchPattern(keyword)
		query = query.Where(
			utilJoinSearchClauses([]string{"host.name", "host.ip", "host.spec", "host.notes"}),
			utilJoinSearchArgs(pattern, 4)...,
		)
	}
	query = query.Order("host.name ASC", "host.id ASC")
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]PortsvcHostRecord, 0, len(rows))
	for idx := range rows {
		out = append(out, *utilPortsvcHostRecordFromModel(&rows[idx]))
	}
	return out, nil
}

func (s *Store) FetchPortsvcHostByID(ctx context.Context, id string) (*PortsvcHostRecord, error) {
	row := new(HostsModel)
	err := s.db.NewSelect().Model(row).Where("host.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcHostRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcHost(ctx context.Context, host *PortsvcHostRecord) (*PortsvcHostRecord, error) {
	row := &HostsModel{
		ID:        idgen.NewUUID7(),
		Name:      host.Name,
		IP:        host.IP,
		Spec:      host.Spec,
		Status:    host.Status,
		Notes:     host.Notes,
		CreatedAt: host.CreatedAt,
		UpdatedAt: host.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return utilPortsvcHostRecordFromModel(row), nil
}

func (s *Store) UpdatePortsvcHost(ctx context.Context, id string, host *PortsvcHostRecord) (*PortsvcHostRecord, error) {
	row := &HostsModel{
		ID:        id,
		Name:      host.Name,
		IP:        host.IP,
		Spec:      host.Spec,
		Status:    host.Status,
		Notes:     host.Notes,
		UpdatedAt: host.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("name", "ip", "spec", "status", "notes", "updated_at").
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
	return s.FetchPortsvcHostByID(ctx, id)
}

func (s *Store) DeletePortsvcHost(ctx context.Context, id string) error {
	res, err := s.db.NewDelete().Model((*HostsModel)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) ListPortsvcPortGroups(ctx context.Context, params PortsvcPortGroupListFilter) ([]PortsvcPortGroupRecord, error) {
	var rows []PortAllocationsModel
	query := s.db.NewSelect().Model(&rows).Relation("Host")
	if !params.Admin {
		query = query.Where("port_group.owner_subject = ?", params.OwnerSubject)
	} else if params.OwnerSubject != "" {
		query = query.Where("port_group.owner_subject = ?", params.OwnerSubject)
	}
	if status := utilCompactString(params.Status); status != "" {
		query = query.Where("port_group.status = ?", status)
	}
	if keyword := utilCompactString(params.Query); keyword != "" {
		pattern := utilSearchPattern(keyword)
		query = query.Where(
			utilJoinSearchClauses([]string{
				"port_group.project_name",
				"port_group.project_owner",
				"port_group.runtime_name",
				"port_group.service_ip",
				"port_group.tags",
				"port_group.notes",
			}),
			utilJoinSearchArgs(pattern, 6)...,
		)
	}
	switch strings.ToLower(params.Sort) {
	case "project":
		query = query.Order("port_group.project_name ASC", "port_group.port_start ASC")
	case "status":
		query = query.Order("port_group.status ASC", "port_group.port_start ASC")
	case "updated_desc":
		query = query.Order("port_group.updated_at DESC", "port_group.port_start ASC")
	default:
		query = query.Order("port_group.port_start ASC")
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]PortsvcPortGroupRecord, 0, len(rows))
	for idx := range rows {
		out = append(out, *utilPortsvcPortGroupRecordFromModel(&rows[idx]))
	}
	return out, nil
}

func (s *Store) FetchPortsvcPortGroupByID(ctx context.Context, id string) (*PortsvcPortGroupRecord, error) {
	row := new(PortAllocationsModel)
	err := s.db.NewSelect().Model(row).Relation("Host").Where("port_group.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcPortGroupRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcPortGroup(ctx context.Context, group *PortsvcPortGroupRecord) (*PortsvcPortGroupRecord, error) {
	row := &PortAllocationsModel{
		ID:           idgen.NewUUID7(),
		OwnerSubject: group.OwnerSubject,
		HostID:       group.HostID,
		PortStart:    group.PortStart,
		PortEnd:      group.PortEnd,
		ProjectName:  group.ProjectName,
		ProjectOwner: group.ProjectOwner,
		RuntimeMode:  group.RuntimeMode,
		RuntimeName:  group.RuntimeName,
		ServiceIP:    group.ServiceIP,
		Status:       group.Status,
		Tags:         group.Tags,
		Notes:        group.Notes,
		CreatedAt:    group.CreatedAt,
		UpdatedAt:    group.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return s.FetchPortsvcPortGroupByID(ctx, row.ID)
}

func (s *Store) UpdatePortsvcPortGroup(ctx context.Context, id string, group *PortsvcPortGroupRecord) (*PortsvcPortGroupRecord, error) {
	row := &PortAllocationsModel{
		ID:           id,
		OwnerSubject: group.OwnerSubject,
		HostID:       group.HostID,
		PortStart:    group.PortStart,
		PortEnd:      group.PortEnd,
		ProjectName:  group.ProjectName,
		ProjectOwner: group.ProjectOwner,
		RuntimeMode:  group.RuntimeMode,
		RuntimeName:  group.RuntimeName,
		ServiceIP:    group.ServiceIP,
		Status:       group.Status,
		Tags:         group.Tags,
		Notes:        group.Notes,
		UpdatedAt:    group.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("owner_subject", "host_id", "port_start", "port_end", "project_name", "project_owner", "runtime_mode", "runtime_name", "service_ip", "status", "tags", "notes", "updated_at").
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
	return s.FetchPortsvcPortGroupByID(ctx, id)
}

func (s *Store) DeletePortsvcPortGroups(ctx context.Context, ids []string, ownerSubject string, admin bool) error {
	query := s.db.NewDelete().Model((*PortAllocationsModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("owner_subject = ?", ownerSubject)
	}
	_, err := query.Exec(ctx)
	return err
}

func (s *Store) CountPortsvcPortGroupsByIDs(ctx context.Context, ids []string, ownerSubject string, admin bool) (int, error) {
	query := s.db.NewSelect().Model((*PortAllocationsModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("owner_subject = ?", ownerSubject)
	}
	return query.Count(ctx)
}

func (s *Store) CountPortsvcOverlappingPortGroups(ctx context.Context, currentID string, group *PortsvcPortGroupRecord) (int, error) {
	query := s.db.NewSelect().
		Model((*PortAllocationsModel)(nil)).
		Where("owner_subject = ?", group.OwnerSubject).
		Where("port_start = ?", group.PortStart)
	if currentID != "" {
		query = query.Where("id != ?", currentID)
	}
	return query.Count(ctx)
}

func (s *Store) ListPortsvcPortGroupStartsByOwner(ctx context.Context, ownerSubject string, excludeID string) ([]int, error) {
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

func (s *Store) FetchPortsvcPortSlotByID(ctx context.Context, id string) (*PortsvcPortSlotRecord, error) {
	row := new(ServicesModel)
	err := s.db.NewSelect().Model(row).Where("slot.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcPortSlotRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcPortSlot(ctx context.Context, slot *PortsvcPortSlotRecord) (*PortsvcPortSlotRecord, error) {
	row := &ServicesModel{
		ID:            idgen.NewUUID7(),
		PortGroupID:   slot.PortGroupID,
		Port:          slot.Port,
		Name:          slot.Name,
		Kind:          slot.Kind,
		Protocol:      slot.Protocol,
		ContainerName: slot.ContainerName,
		Status:        slot.Status,
		Notes:         slot.Notes,
		CreatedAt:     slot.CreatedAt,
		UpdatedAt:     slot.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return utilPortsvcPortSlotRecordFromModel(row), nil
}

func (s *Store) UpdatePortsvcPortSlot(ctx context.Context, id string, slot *PortsvcPortSlotRecord) (*PortsvcPortSlotRecord, error) {
	row := &ServicesModel{
		ID:            id,
		PortGroupID:   slot.PortGroupID,
		Port:          slot.Port,
		Name:          slot.Name,
		Kind:          slot.Kind,
		Protocol:      slot.Protocol,
		ContainerName: slot.ContainerName,
		Status:        slot.Status,
		Notes:         slot.Notes,
		UpdatedAt:     slot.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("port_group_id", "port", "name", "kind", "protocol", "container_name", "status", "notes", "updated_at").
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
	return s.FetchPortsvcPortSlotByID(ctx, id)
}

func (s *Store) DeletePortsvcPortSlot(ctx context.Context, id string) error {
	res, err := s.db.NewDelete().Model((*ServicesModel)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) ListPortsvcPortGroupChildrenByGroupIDs(ctx context.Context, ids []string) ([]PortsvcPortGroupChildrenRecord, error) {
	children := make(map[string]PortsvcPortGroupChildrenRecord, len(ids))
	for _, id := range ids {
		children[id] = PortsvcPortGroupChildrenRecord{PortGroupID: id}
	}

	var slotRows []ServicesModel
	if err := s.db.NewSelect().
		Model(&slotRows).
		Where("slot.port_group_id IN (?)", bun.List(ids)).
		Order("slot.port ASC").
		Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range slotRows {
		record := *utilPortsvcPortSlotRecordFromModel(&slotRows[idx])
		child := children[record.PortGroupID]
		child.Slots = append(child.Slots, record)
		children[record.PortGroupID] = child
	}

	var repoRows []RepositoriesModel
	if err := s.db.NewSelect().
		Model(&repoRows).
		Where("repository_ref.port_group_id IN (?)", bun.List(ids)).
		Order("repository_ref.name ASC").
		Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range repoRows {
		record := *utilPortsvcRepositoryRecordFromModel(&repoRows[idx])
		child := children[record.PortGroupID]
		child.Repositories = append(child.Repositories, record)
		children[record.PortGroupID] = child
	}

	var depRows []DependenciesModel
	if err := s.db.NewSelect().
		Model(&depRows).
		Where("dependency.port_group_id IN (?)", bun.List(ids)).
		Order("dependency.name ASC").
		Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range depRows {
		record := *utilPortsvcDependencyRecordFromModel(&depRows[idx])
		child := children[record.PortGroupID]
		child.Dependencies = append(child.Dependencies, record)
		children[record.PortGroupID] = child
	}

	out := make([]PortsvcPortGroupChildrenRecord, 0, len(children))
	for _, id := range ids {
		out = append(out, children[id])
	}
	return out, nil
}

func (s *Store) ReplacePortsvcPortGroupChildren(ctx context.Context, portGroupID string) error {
	if _, err := s.db.NewDelete().Model((*ServicesModel)(nil)).Where("port_group_id = ?", portGroupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*RepositoriesModel)(nil)).Where("port_group_id = ?", portGroupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*DependenciesModel)(nil)).Where("port_group_id = ?", portGroupID).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Store) AddPortsvcPortSlots(ctx context.Context, slots []PortsvcPortSlotRecord) error {
	if len(slots) == 0 {
		return nil
	}
	rows := make([]ServicesModel, 0, len(slots))
	for _, slot := range slots {
		rows = append(rows, ServicesModel{
			ID:            idgen.NewUUID7(),
			PortGroupID:   slot.PortGroupID,
			Port:          slot.Port,
			Name:          slot.Name,
			Kind:          slot.Kind,
			Protocol:      slot.Protocol,
			ContainerName: slot.ContainerName,
			Status:        slot.Status,
			Notes:         slot.Notes,
			CreatedAt:     slot.CreatedAt,
			UpdatedAt:     slot.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
}

func (s *Store) AddPortsvcRepositories(ctx context.Context, repos []PortsvcRepositoryRecord) error {
	if len(repos) == 0 {
		return nil
	}
	rows := make([]RepositoriesModel, 0, len(repos))
	for _, repo := range repos {
		rows = append(rows, RepositoriesModel{
			ID:           idgen.NewUUID7(),
			OwnerSubject: repo.OwnerSubject,
			PortGroupID:  repo.PortGroupID,
			Name:         repo.Name,
			URL:          repo.URL,
			Kind:         repo.Kind,
			Notes:        repo.Notes,
			CreatedAt:    repo.CreatedAt,
			UpdatedAt:    repo.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
}

func (s *Store) AddPortsvcDependencies(ctx context.Context, dependencies []PortsvcDependencyRecord) error {
	if len(dependencies) == 0 {
		return nil
	}
	rows := make([]DependenciesModel, 0, len(dependencies))
	for _, dependency := range dependencies {
		rows = append(rows, DependenciesModel{
			ID:           idgen.NewUUID7(),
			OwnerSubject: dependency.OwnerSubject,
			PortGroupID:  dependency.PortGroupID,
			Name:         dependency.Name,
			Type:         dependency.Type,
			URL:          dependency.URL,
			Version:      dependency.Version,
			Notes:        dependency.Notes,
			CreatedAt:    dependency.CreatedAt,
			UpdatedAt:    dependency.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
}
