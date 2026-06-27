package repository

import (
	"context"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

func (s *Store) ListInventoryPortGroups(ctx context.Context, params InventoryPortGroupListFilter) ([]InventoryPortGroupRecord, error) {
	var groups []InventoryPortGroupModel
	query := s.db.NewSelect().Model(&groups).Relation("User")
	if !params.Admin {
		query = query.Where("allocation.user_id = ?", params.UserID)
	} else if params.UserID > 0 {
		query = query.Where("allocation.user_id = ?", params.UserID)
	}
	if status := utilCompactString(params.Status); status != "" {
		query = query.Where("allocation.status = ?", status)
	}
	if dindIP := utilCompactString(params.DindIP); dindIP != "" {
		query = query.Where("allocation.dind_ip = ?", dindIP)
	}
	if projectName := utilCompactString(params.ProjectName); projectName != "" {
		pattern := utilSearchPattern(projectName)
		query = query.Where("EXISTS (SELECT 1 FROM allocation_projects project WHERE project.allocation_id = allocation.id AND LOWER(project.name) LIKE ?)", pattern)
	}
	if keyword := utilCompactString(params.Query); keyword != "" {
		pattern := utilSearchPattern(keyword)
		query = query.Where(
			utilJoinSearchClauses([]string{
				"allocation.name",
				"allocation.dind_ip",
				"allocation.dind_container",
				"allocation.owner",
				"allocation.tags",
				"allocation.notes",
			})+" OR EXISTS (SELECT 1 FROM users user_search WHERE user_search.id = allocation.user_id AND (LOWER(user_search.username) LIKE ? OR LOWER(user_search.display_name) LIKE ?)) OR EXISTS (SELECT 1 FROM allocation_projects project WHERE project.allocation_id = allocation.id AND LOWER(project.name) LIKE ?)",
			append(utilJoinSearchArgs(pattern, 6), pattern, pattern, pattern)...,
		)
	}
	switch strings.ToLower(params.Sort) {
	case "name":
		query = query.Order("allocation.name ASC", "allocation.port_start ASC")
	case "status":
		query = query.Order("allocation.status ASC", "allocation.port_start ASC")
	case "user":
		query = query.Order("allocation.user_id ASC", "allocation.port_start ASC")
	case "updated_desc":
		query = query.Order("allocation.updated_at DESC", "allocation.port_start ASC")
	default:
		query = query.Order("allocation.port_start ASC")
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
	err := s.db.NewSelect().Model(group).Relation("User").Where("allocation.id = ?", id).Scan(ctx)
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
		UserID:        group.UserID,
		PortStart:     group.PortStart,
		PortEnd:       group.PortEnd,
		Name:          group.Name,
		DindIP:        group.DindIP,
		DindContainer: group.DindContainer,
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
	return s.FetchInventoryPortGroupWithHostByID(ctx, row.ID)
}

func (s *Store) UpdateInventoryPortGroup(ctx context.Context, id int64, group *InventoryPortGroupRecord) (*InventoryPortGroupRecord, error) {
	row := &InventoryPortGroupModel{
		ID:            id,
		UserID:        group.UserID,
		PortStart:     group.PortStart,
		PortEnd:       group.PortEnd,
		Name:          group.Name,
		DindIP:        group.DindIP,
		DindContainer: group.DindContainer,
		Status:        group.Status,
		Owner:         group.Owner,
		Tags:          group.Tags,
		Notes:         group.Notes,
		UpdatedAt:     group.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("user_id", "port_start", "port_end", "name", "dind_ip", "dind_container", "status", "owner", "tags", "notes", "updated_at").
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
	return s.FetchInventoryPortGroupWithHostByID(ctx, id)
}

func (s *Store) CountInventoryOverlappingPortGroups(ctx context.Context, currentID int64, group *InventoryPortGroupRecord) (int, error) {
	query := s.db.NewSelect().
		Model((*InventoryPortGroupModel)(nil)).
		Where("user_id = ?", group.UserID).
		Where("port_start = ?", group.PortStart)
	if currentID > 0 {
		query = query.Where("id != ?", currentID)
	}
	return query.Count(ctx)
}

func (s *Store) CountInventoryPortGroupsByIDs(ctx context.Context, ids []int64, userID int64, admin bool) (int, error) {
	query := s.db.NewSelect().Model((*InventoryPortGroupModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("user_id = ?", userID)
	}
	return query.Count(ctx)
}

func (s *Store) ListInventoryPortGroupsByIDs(ctx context.Context, ids []int64, userID int64, admin bool) ([]InventoryPortGroupRecord, error) {
	var groups []InventoryPortGroupModel
	query := s.db.NewSelect().Model(&groups).Relation("User").Where("allocation.id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("allocation.user_id = ?", userID)
	}
	query = query.Order("allocation.port_start ASC")
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]InventoryPortGroupRecord, 0, len(groups))
	for idx := range groups {
		out = append(out, *utilInventoryPortGroupRecordFromModel(&groups[idx]))
	}
	return out, nil
}

func (s *Store) UpdateInventoryPortGroupsBatch(ctx context.Context, ids []int64, status *string, owner *string, tags *string, updatedAt time.Time, userID int64, admin bool) (*InventoryPortGroupRecord, error) {
	query := s.db.NewUpdate().
		Model((*InventoryPortGroupModel)(nil)).
		Set("updated_at = ?", updatedAt).
		Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("user_id = ?", userID)
	}
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

func (s *Store) DeleteInventoryPortGroups(ctx context.Context, ids []int64, userID int64, admin bool) error {
	query := s.db.NewDelete().Model((*InventoryPortGroupModel)(nil)).Where("id IN (?)", bun.List(ids))
	if !admin {
		query = query.Where("user_id = ?", userID)
	}
	_, err := query.Exec(ctx)
	return err
}

func (s *Store) ListInventoryPortStartsByUser(ctx context.Context, userID int64, excludeID int64) ([]int, error) {
	var rows []struct {
		PortStart int `bun:"port_start"`
	}
	query := s.db.NewSelect().Model((*InventoryPortGroupModel)(nil)).Column("port_start").Where("user_id = ?", userID)
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

func (s *Store) ListInventoryPortGroupChildrenByPortGroupIDs(ctx context.Context, ids []int64) ([]InventoryPortGroupChildrenRecord, error) {
	children := make(map[int64]InventoryPortGroupChildrenRecord, len(ids))
	for _, id := range ids {
		children[id] = InventoryPortGroupChildrenRecord{PortGroupID: id}
	}

	var slots []InventoryPortSlotModel
	if err := s.db.NewSelect().Model(&slots).Where("allocation_id IN (?)", bun.List(ids)).Order("port ASC").Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range slots {
		record := *utilInventoryPortSlotRecordFromModel(&slots[idx])
		child := children[record.PortGroupID]
		child.Slots = append(child.Slots, record)
		children[record.PortGroupID] = child
	}

	var projects []InventoryProjectModel
	if err := s.db.NewSelect().Model(&projects).Where("allocation_id IN (?)", bun.List(ids)).Order("id ASC").Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range projects {
		record := *utilInventoryProjectRecordFromModel(&projects[idx])
		child := children[record.PortGroupID]
		child.Projects = append(child.Projects, record)
		children[record.PortGroupID] = child
	}

	var components []InventoryComponentModel
	if err := s.db.NewSelect().Model(&components).Where("allocation_id IN (?)", bun.List(ids)).Order("id ASC").Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range components {
		record := *utilInventoryComponentRecordFromModel(&components[idx])
		child := children[record.PortGroupID]
		child.Components = append(child.Components, record)
		children[record.PortGroupID] = child
	}

	var repositories []InventoryRepositoryRefModel
	if err := s.db.NewSelect().Model(&repositories).Where("allocation_id IN (?)", bun.List(ids)).Order("id ASC").Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range repositories {
		record := *utilInventoryRepositoryRefRecordFromModel(&repositories[idx])
		child := children[record.PortGroupID]
		child.Repositories = append(child.Repositories, record)
		children[record.PortGroupID] = child
	}

	out := make([]InventoryPortGroupChildrenRecord, 0, len(children))
	for _, id := range ids {
		out = append(out, children[id])
	}
	return out, nil
}

func (s *Store) ReplaceInventoryPortGroupChildren(ctx context.Context, groupID int64) error {
	if _, err := s.db.NewDelete().Model((*InventoryRepositoryRefModel)(nil)).Where("allocation_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*InventoryProjectModel)(nil)).Where("allocation_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*InventoryComponentModel)(nil)).Where("allocation_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*InventoryPortSlotModel)(nil)).Where("allocation_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Store) AddInventoryPortSlots(ctx context.Context, slots []InventoryPortSlotRecord) error {
	if len(slots) == 0 {
		return nil
	}
	rows := make([]InventoryPortSlotModel, 0, len(slots))
	for _, slot := range slots {
		rows = append(rows, InventoryPortSlotModel{
			PortGroupID: slot.PortGroupID,
			Port:        slot.Port,
			Name:        slot.Name,
			Protocol:    slot.Protocol,
			Purpose:     slot.Purpose,
			Status:      slot.Status,
			Notes:       slot.Notes,
			CreatedAt:   slot.CreatedAt,
			UpdatedAt:   slot.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
}

func (s *Store) AddInventoryProjects(ctx context.Context, projects []InventoryProjectRecord) error {
	if len(projects) == 0 {
		return nil
	}
	rows := make([]InventoryProjectModel, 0, len(projects))
	for _, project := range projects {
		rows = append(rows, InventoryProjectModel{
			PortGroupID: project.PortGroupID,
			Name:        project.Name,
			Description: project.Description,
			Notes:       project.Notes,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
}

func (s *Store) AddInventoryComponents(ctx context.Context, components []InventoryComponentRecord) error {
	if len(components) == 0 {
		return nil
	}
	rows := make([]InventoryComponentModel, 0, len(components))
	for _, component := range components {
		rows = append(rows, InventoryComponentModel{
			PortGroupID: component.PortGroupID,
			Name:        component.Name,
			Type:        component.Type,
			URL:         component.URL,
			Version:     component.Version,
			Notes:       component.Notes,
			CreatedAt:   component.CreatedAt,
			UpdatedAt:   component.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
}

func (s *Store) AddInventoryRepositoryRefs(ctx context.Context, repositories []InventoryRepositoryRefRecord) error {
	if len(repositories) == 0 {
		return nil
	}
	rows := make([]InventoryRepositoryRefModel, 0, len(repositories))
	for _, repo := range repositories {
		rows = append(rows, InventoryRepositoryRefModel{
			PortGroupID: repo.PortGroupID,
			ProjectID:   repo.ProjectID,
			Name:        repo.Name,
			URL:         repo.URL,
			Kind:        repo.Kind,
			Notes:       repo.Notes,
			CreatedAt:   repo.CreatedAt,
			UpdatedAt:   repo.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
}
