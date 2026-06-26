package repository

import (
	"context"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

func (s *Store) ListInventoryHosts(ctx context.Context, params InventoryHostListFilter) ([]InventoryHostRecord, error) {
	var hosts []InventoryHostModel
	query := s.db.NewSelect().Model(&hosts)
	if environment := utilCompactString(params.Environment); environment != "" {
		query = query.Where("environment = ?", environment)
	}
	if keyword := utilCompactString(params.Query); keyword != "" {
		pattern := utilSearchPattern(keyword)
		query = query.Where(utilJoinSearchClauses([]string{"ip", "name", "network", "environment", "notes"}), utilJoinSearchArgs(pattern, 5)...)
	}
	switch strings.ToLower(params.Sort) {
	case "environment":
		query = query.Order("environment ASC", "ip ASC")
	case "updated_desc":
		query = query.Order("updated_at DESC", "ip ASC")
	default:
		query = query.Order("ip ASC")
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]InventoryHostRecord, 0, len(hosts))
	for idx := range hosts {
		out = append(out, *utilInventoryHostRecordFromModel(&hosts[idx]))
	}
	return out, nil
}

func (s *Store) CreateInventoryHost(ctx context.Context, host *InventoryHostRecord) (*InventoryHostRecord, error) {
	row := &InventoryHostModel{
		ID:          host.ID,
		IP:          host.IP,
		Name:        host.Name,
		Network:     host.Network,
		Environment: host.Environment,
		Notes:       host.Notes,
		CreatedAt:   host.CreatedAt,
		UpdatedAt:   host.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return utilInventoryHostRecordFromModel(row), nil
}

func (s *Store) UpdateInventoryHost(ctx context.Context, id int64, host *InventoryHostRecord) (*InventoryHostRecord, error) {
	row := &InventoryHostModel{
		ID:          host.ID,
		IP:          host.IP,
		Name:        host.Name,
		Network:     host.Network,
		Environment: host.Environment,
		Notes:       host.Notes,
		CreatedAt:   host.CreatedAt,
		UpdatedAt:   host.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("ip", "name", "network", "environment", "notes", "updated_at").
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
	out, err := s.FetchInventoryHostByID(ctx, id)
	return out, err
}

func (s *Store) DeleteInventoryHost(ctx context.Context, id int64) (bool, error) {
	res, err := s.db.NewDelete().Model((*InventoryHostModel)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (s *Store) FetchInventoryHostByID(ctx context.Context, id int64) (*InventoryHostRecord, error) {
	host := new(InventoryHostModel)
	err := s.db.NewSelect().Model(host).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilInventoryHostRecordFromModel(host), nil
}

func (s *Store) CountInventoryPortGroupsByHostID(ctx context.Context, hostID int64) (int, error) {
	return s.db.NewSelect().Model((*InventoryPortGroupModel)(nil)).Where("host_id = ?", hostID).Count(ctx)
}

func (s *Store) ListInventoryPortGroups(ctx context.Context, params InventoryPortGroupListFilter) ([]InventoryPortGroupRecord, error) {
	var groups []InventoryPortGroupModel
	query := s.db.NewSelect().Model(&groups).Relation("InventoryHostModel")
	if params.HostID > 0 {
		query = query.Where("port_group.host_id = ?", params.HostID)
	}
	if status := utilCompactString(params.Status); status != "" {
		query = query.Where("port_group.status = ?", status)
	}
	if keyword := utilCompactString(params.Query); keyword != "" {
		pattern := utilSearchPattern(keyword)
		query = query.Where(
			utilJoinSearchClauses([]string{
				"port_group.service_name",
				"port_group.container_name",
				"port_group.dind_host",
				"port_group.owner",
				"port_group.tags",
				"port_group.notes",
				"host.ip",
				"host.name",
			}),
			utilJoinSearchArgs(pattern, 8)...,
		)
	}
	switch strings.ToLower(params.Sort) {
	case "service":
		query = query.Order("port_group.service_name ASC", "host.ip ASC", "port_group.port_start ASC")
	case "status":
		query = query.Order("port_group.status ASC", "host.ip ASC", "port_group.port_start ASC")
	case "updated_desc":
		query = query.Order("port_group.updated_at DESC", "host.ip ASC", "port_group.port_start ASC")
	default:
		query = query.Order("host.ip ASC", "port_group.port_start ASC")
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]InventoryPortGroupRecord, 0, len(groups))
	for idx := range groups {
		out = append(out, *utilInventoryPortGroupRecordFromModel(&groups[idx]))
	}
	return out, nil
}

func (s *Store) FetchInventoryPortGroupWithHostByID(ctx context.Context, id int64) (*InventoryPortGroupRecord, error) {
	group := new(InventoryPortGroupModel)
	err := s.db.NewSelect().Model(group).Relation("InventoryHostModel").Where("port_group.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilInventoryPortGroupRecordFromModel(group), nil
}

func (s *Store) FetchInventoryPortGroupByID(ctx context.Context, id int64) (*InventoryPortGroupRecord, error) {
	group := new(InventoryPortGroupModel)
	err := s.db.NewSelect().Model(group).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilInventoryPortGroupRecordFromModel(group), nil
}

func (s *Store) CreateInventoryPortGroup(ctx context.Context, group *InventoryPortGroupRecord) (*InventoryPortGroupRecord, error) {
	row := &InventoryPortGroupModel{
		ID:            group.ID,
		HostID:        group.HostID,
		PortStart:     group.PortStart,
		PortEnd:       group.PortEnd,
		ServiceName:   group.ServiceName,
		ContainerName: group.ContainerName,
		DindHost:      group.DindHost,
		Status:        group.Status,
		Owner:         group.Owner,
		Tags:          group.Tags,
		Notes:         group.Notes,
		CreatedAt:     group.CreatedAt,
		UpdatedAt:     group.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return s.FetchInventoryPortGroupByID(ctx, row.ID)
}

func (s *Store) UpdateInventoryPortGroup(ctx context.Context, id int64, group *InventoryPortGroupRecord) (*InventoryPortGroupRecord, error) {
	row := &InventoryPortGroupModel{
		ID:            group.ID,
		HostID:        group.HostID,
		PortStart:     group.PortStart,
		PortEnd:       group.PortEnd,
		ServiceName:   group.ServiceName,
		ContainerName: group.ContainerName,
		DindHost:      group.DindHost,
		Status:        group.Status,
		Owner:         group.Owner,
		Tags:          group.Tags,
		Notes:         group.Notes,
		CreatedAt:     group.CreatedAt,
		UpdatedAt:     group.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("host_id", "port_start", "port_end", "service_name", "container_name", "dind_host", "status", "owner", "tags", "notes", "updated_at").
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
	out, err := s.FetchInventoryPortGroupByID(ctx, id)
	return out, err
}

func (s *Store) CountInventoryOverlappingPortGroups(ctx context.Context, currentID int64, group *InventoryPortGroupRecord) (int, error) {
	query := s.db.NewSelect().
		Model((*InventoryPortGroupModel)(nil)).
		Where("host_id = ?", group.HostID).
		Where("NOT (port_end < ? OR port_start > ?)", group.PortStart, group.PortEnd)
	if currentID > 0 {
		query = query.Where("id != ?", currentID)
	}
	return query.Count(ctx)
}

func (s *Store) CountInventoryPortGroupsByIDs(ctx context.Context, ids []int64) (int, error) {
	return s.db.NewSelect().Model((*InventoryPortGroupModel)(nil)).Where("id IN (?)", bun.List(ids)).Count(ctx)
}

func (s *Store) UpdateInventoryPortGroupsBatch(ctx context.Context, ids []int64, status *string, owner *string, tags *string, updatedAt time.Time) (*InventoryPortGroupRecord, error) {
	query := s.db.NewUpdate().
		Model((*InventoryPortGroupModel)(nil)).
		Set("updated_at = ?", updatedAt).
		Where("id IN (?)", bun.List(ids))
	if status != nil {
		query = query.Set("status = ?", *status)
	}
	if owner != nil {
		query = query.Set("owner = ?", *owner)
	}
	if tags != nil {
		query = query.Set("tags = ?", *tags)
	}
	if _, err := query.Exec(ctx); err != nil {
		return nil, err
	}
	return s.FetchInventoryPortGroupByID(ctx, ids[0])
}

func (s *Store) DeleteInventoryPortGroups(ctx context.Context, ids []int64) error {
	if _, err := s.db.NewDelete().Model((*InventoryRepositoryRefModel)(nil)).Where("port_group_id IN (?)", bun.List(ids)).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*InventoryComponentModel)(nil)).Where("port_group_id IN (?)", bun.List(ids)).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*InventoryPortSlotModel)(nil)).Where("port_group_id IN (?)", bun.List(ids)).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*InventoryPortGroupModel)(nil)).Where("id IN (?)", bun.List(ids)).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Store) ReplaceInventoryPortGroupChildren(ctx context.Context, groupID int64) error {
	if _, err := s.db.NewDelete().Model((*InventoryRepositoryRefModel)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*InventoryComponentModel)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*InventoryPortSlotModel)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Store) AddInventoryPortSlots(ctx context.Context, slots []InventoryPortSlotRecord) error {
	if len(slots) == 0 {
		return nil
	}
	_, err := s.db.NewInsert().Model(&slots).Exec(ctx)
	return err
}

func (s *Store) AddInventoryComponents(ctx context.Context, components []InventoryComponentRecord) error {
	if len(components) == 0 {
		return nil
	}
	_, err := s.db.NewInsert().Model(&components).Exec(ctx)
	return err
}

func (s *Store) AddInventoryRepositoryRefs(ctx context.Context, repositories []InventoryRepositoryRefRecord) error {
	if len(repositories) == 0 {
		return nil
	}
	_, err := s.db.NewInsert().Model(&repositories).Exec(ctx)
	return err
}

func (s *Store) ListInventoryPortGroupsByIDs(ctx context.Context, ids []int64) ([]InventoryPortGroupRecord, error) {
	var groups []InventoryPortGroupModel
	if err := s.db.NewSelect().
		Model(&groups).
		Relation("InventoryHostModel").
		Where("port_group.id IN (?)", bun.List(ids)).
		Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]InventoryPortGroupRecord, 0, len(groups))
	for idx := range groups {
		out = append(out, *utilInventoryPortGroupRecordFromModel(&groups[idx]))
	}
	return out, nil
}

func (s *Store) FetchInventoryPortGroupChildrenByPortGroupIDs(ctx context.Context, ids []int64) (*InventoryPortGroupChildrenRecord, error) {
	var slots []InventoryPortSlotModel
	if err := s.db.NewSelect().Model(&slots).Where("port_group_id IN (?)", bun.List(ids)).Order("port ASC").Scan(ctx); err != nil {
		return nil, err
	}
	var components []InventoryComponentModel
	if err := s.db.NewSelect().Model(&components).Where("port_group_id IN (?)", bun.List(ids)).Order("name ASC").Scan(ctx); err != nil {
		return nil, err
	}
	var repositories []InventoryRepositoryRefModel
	if err := s.db.NewSelect().Model(&repositories).Where("port_group_id IN (?)", bun.List(ids)).Order("name ASC").Scan(ctx); err != nil {
		return nil, err
	}
	out := &InventoryPortGroupChildrenRecord{
		Slots:        make([]InventoryPortSlotRecord, 0, len(slots)),
		Components:   make([]InventoryComponentRecord, 0, len(components)),
		Repositories: make([]InventoryRepositoryRefRecord, 0, len(repositories)),
	}
	for idx := range slots {
		out.Slots = append(out.Slots, *utilInventoryPortSlotRecordFromModel(&slots[idx]))
	}
	for idx := range components {
		out.Components = append(out.Components, *utilInventoryComponentRecordFromModel(&components[idx]))
	}
	for idx := range repositories {
		out.Repositories = append(out.Repositories, *utilInventoryRepositoryRefRecordFromModel(&repositories[idx]))
	}
	return out, nil
}
