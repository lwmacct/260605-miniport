package repository

import (
	"context"
	"strings"

	"github.com/lwmacct/260630-go-hsr-shared/pkg/idgen"
	"github.com/uptrace/bun"
)

func (s *Store) ListPortsvcHosts(ctx context.Context, params PortsvcHostListFilter) ([]PortsvcHostRecord, error) {
	var rows []HostsModel
	query := utilFilterPortsvcList(
		s.db.NewSelect().Model(&rows), "host", params.Status, params.Query,
		[]string{"host.name", "host.ip", "host.spec", "host.notes"},
	).Order("host.name ASC", "host.id ASC")
	return utilScanPortsvcList(ctx, query, &rows, utilPortsvcHostRecordFromModel)
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

func (s *Store) ListPortsvcDependencyAssets(ctx context.Context, params PortsvcDependencyAssetListFilter) ([]PortsvcDependencyAssetRecord, error) {
	var rows []DependenciesModel
	query := s.db.NewSelect().Model(&rows)
	if assetKind := utilCompactString(params.AssetKind); assetKind != "" {
		query = query.Where("asset.asset_kind = ?", assetKind)
	}
	if assetType := utilCompactString(params.AssetType); assetType != "" {
		query = query.Where("asset.asset_type = ?", assetType)
	}
	if provider := utilCompactString(params.Provider); provider != "" {
		query = query.Where("asset.provider = ?", provider)
	}
	if status := utilCompactString(params.Status); status != "" {
		query = query.Where("asset.status = ?", status)
	}
	if keyword := utilCompactString(params.Query); keyword != "" {
		pattern := utilSearchPattern(keyword)
		query = query.Where(
			utilJoinSearchClauses([]string{
				"asset.name",
				"asset.url",
				"asset.full_name",
				"asset.description",
				"asset.notes",
			}),
			utilJoinSearchArgs(pattern, 5)...,
		)
	}
	query = query.Order("asset.name ASC", "asset.id ASC")
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]PortsvcDependencyAssetRecord, 0, len(rows))
	for idx := range rows {
		out = append(out, *utilPortsvcDependencyAssetRecordFromModel(&rows[idx]))
	}
	return out, nil
}

func (s *Store) FetchPortsvcDependencyAssetByID(ctx context.Context, id string) (*PortsvcDependencyAssetRecord, error) {
	row := new(DependenciesModel)
	err := s.db.NewSelect().Model(row).Where("asset.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcDependencyAssetRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcDependencyAsset(ctx context.Context, asset *PortsvcDependencyAssetRecord) (*PortsvcDependencyAssetRecord, error) {
	row := &DependenciesModel{
		ID:              idgen.NewUUID7(),
		Name:            asset.Name,
		AssetKind:       asset.AssetKind,
		AssetType:       asset.AssetType,
		Provider:        asset.Provider,
		URL:             asset.URL,
		FullName:        asset.FullName,
		ExternalID:      asset.ExternalID,
		Visibility:      asset.Visibility,
		Controllability: asset.Controllability,
		Status:          asset.Status,
		Description:     asset.Description,
		Metadata:        asset.Metadata,
		LastSyncedAt:    asset.LastSyncedAt,
		Notes:           asset.Notes,
		CreatedAt:       asset.CreatedAt,
		UpdatedAt:       asset.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return utilPortsvcDependencyAssetRecordFromModel(row), nil
}

func (s *Store) UpdatePortsvcDependencyAsset(ctx context.Context, id string, asset *PortsvcDependencyAssetRecord) (*PortsvcDependencyAssetRecord, error) {
	row := &DependenciesModel{
		ID:              id,
		Name:            asset.Name,
		AssetKind:       asset.AssetKind,
		AssetType:       asset.AssetType,
		Provider:        asset.Provider,
		URL:             asset.URL,
		FullName:        asset.FullName,
		ExternalID:      asset.ExternalID,
		Visibility:      asset.Visibility,
		Controllability: asset.Controllability,
		Status:          asset.Status,
		Description:     asset.Description,
		Metadata:        asset.Metadata,
		LastSyncedAt:    asset.LastSyncedAt,
		Notes:           asset.Notes,
		UpdatedAt:       asset.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("name", "asset_kind", "asset_type", "provider", "url", "full_name", "external_id", "visibility", "controllability", "status", "description", "metadata", "last_synced_at", "notes", "updated_at").
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
	return s.FetchPortsvcDependencyAssetByID(ctx, id)
}

func (s *Store) DeletePortsvcDependencyAsset(ctx context.Context, id string) error {
	res, err := s.db.NewDelete().Model((*DependenciesModel)(nil)).Where("id = ?", id).Exec(ctx)
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

func (s *Store) ListPortsvcServiceGroups(ctx context.Context, params PortsvcServiceGroupListFilter) ([]PortsvcServiceGroupRecord, error) {
	var rows []ServiceGroupsModel
	query := utilFilterPortsvcList(
		s.db.NewSelect().Model(&rows), "service_group", params.Status, params.Query,
		[]string{"service_group.name", "service_group.kind", "service_group.description", "service_group.notes"},
	).Order("service_group.name ASC", "service_group.id ASC")
	return utilScanPortsvcList(ctx, query, &rows, utilPortsvcServiceGroupRecordFromModel)
}

func (s *Store) FetchPortsvcServiceGroupByID(ctx context.Context, id string) (*PortsvcServiceGroupRecord, error) {
	row := new(ServiceGroupsModel)
	err := s.db.NewSelect().Model(row).Where("service_group.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return utilPortsvcServiceGroupRecordFromModel(row), nil
}

func (s *Store) CreatePortsvcServiceGroup(ctx context.Context, group *PortsvcServiceGroupRecord) (*PortsvcServiceGroupRecord, error) {
	row := &ServiceGroupsModel{
		ID:          idgen.NewUUID7(),
		Name:        group.Name,
		Kind:        group.Kind,
		Status:      group.Status,
		Description: group.Description,
		Notes:       group.Notes,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return utilPortsvcServiceGroupRecordFromModel(row), nil
}

func (s *Store) UpdatePortsvcServiceGroup(ctx context.Context, id string, group *PortsvcServiceGroupRecord) (*PortsvcServiceGroupRecord, error) {
	row := &ServiceGroupsModel{
		ID:          id,
		Name:        group.Name,
		Kind:        group.Kind,
		Status:      group.Status,
		Description: group.Description,
		Notes:       group.Notes,
		UpdatedAt:   group.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("name", "kind", "status", "description", "notes", "updated_at").
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
	return s.FetchPortsvcServiceGroupByID(ctx, id)
}

func (s *Store) DeletePortsvcServiceGroup(ctx context.Context, id string) error {
	res, err := s.db.NewDelete().Model((*ServiceGroupsModel)(nil)).Where("id = ?", id).Exec(ctx)
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

func (s *Store) CountPortsvcServiceGroupsByIDs(ctx context.Context, ids []string) (int, error) {
	return s.db.NewSelect().Model((*ServiceGroupsModel)(nil)).Where("id IN (?)", bun.List(ids)).Count(ctx)
}

func (s *Store) ListPortsvcPortGroups(ctx context.Context, params PortsvcPortGroupListFilter) ([]PortsvcPortGroupRecord, error) {
	var rows []PortAllocationsModel
	query := s.db.NewSelect().Model(&rows).Relation("Host")
	if status := utilCompactString(params.Status); status != "" {
		query = query.Where("port_group.status = ?", status)
	}
	if keyword := utilCompactString(params.Query); keyword != "" {
		pattern := utilSearchPattern(keyword)
		query = query.Where(
			utilJoinSearchClauses([]string{
				"port_group.environment_name",
				"port_group.environment_owner",
				"port_group.runtime_name",
				"port_group.service_ip",
				"port_group.tags",
				"port_group.notes",
			}),
			utilJoinSearchArgs(pattern, 6)...,
		)
	}
	switch strings.ToLower(params.Sort) {
	case "environment":
		query = query.Order("port_group.environment_name ASC", "port_group.port_prefix ASC")
	case "status":
		query = query.Order("port_group.status ASC", "port_group.port_prefix ASC")
	case "updated_desc":
		query = query.Order("port_group.updated_at DESC", "port_group.port_prefix ASC")
	default:
		query = query.Order("port_group.port_prefix ASC")
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
		ID:               idgen.NewUUID7(),
		HostID:           group.HostID,
		PortPrefix:       group.PortPrefix,
		EnvironmentName:  group.EnvironmentName,
		EnvironmentOwner: group.EnvironmentOwner,
		RuntimeMode:      group.RuntimeMode,
		RuntimeName:      group.RuntimeName,
		ServiceIP:        group.ServiceIP,
		Status:           group.Status,
		Tags:             group.Tags,
		Notes:            group.Notes,
		CreatedAt:        group.CreatedAt,
		UpdatedAt:        group.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return s.FetchPortsvcPortGroupByID(ctx, row.ID)
}

func (s *Store) UpdatePortsvcPortGroup(ctx context.Context, id string, group *PortsvcPortGroupRecord) (*PortsvcPortGroupRecord, error) {
	row := &PortAllocationsModel{
		ID:               id,
		HostID:           group.HostID,
		PortPrefix:       group.PortPrefix,
		EnvironmentName:  group.EnvironmentName,
		EnvironmentOwner: group.EnvironmentOwner,
		RuntimeMode:      group.RuntimeMode,
		RuntimeName:      group.RuntimeName,
		ServiceIP:        group.ServiceIP,
		Status:           group.Status,
		Tags:             group.Tags,
		Notes:            group.Notes,
		UpdatedAt:        group.UpdatedAt,
	}
	res, err := s.db.NewUpdate().
		Model(row).
		Column("host_id", "port_prefix", "environment_name", "environment_owner", "runtime_mode", "runtime_name", "service_ip", "status", "tags", "notes", "updated_at").
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

func (s *Store) DeletePortsvcPortGroups(ctx context.Context, ids []string) error {
	_, err := s.db.NewDelete().Model((*PortAllocationsModel)(nil)).Where("id IN (?)", bun.List(ids)).Exec(ctx)
	return err
}

func (s *Store) CountPortsvcPortGroupsByIDs(ctx context.Context, ids []string) (int, error) {
	return s.db.NewSelect().Model((*PortAllocationsModel)(nil)).Where("id IN (?)", bun.List(ids)).Count(ctx)
}

func (s *Store) CountPortsvcOverlappingPortGroups(ctx context.Context, currentID string, group *PortsvcPortGroupRecord) (int, error) {
	query := s.db.NewSelect().
		Model((*PortAllocationsModel)(nil)).
		Where("port_prefix = ?", group.PortPrefix)
	if currentID != "" {
		query = query.Where("id != ?", currentID)
	}
	return query.Count(ctx)
}

func (s *Store) ListPortsvcPortGroupPrefixes(ctx context.Context, excludeID string) ([]int, error) {
	var rows []struct {
		PortPrefix int `bun:"port_prefix"`
	}
	query := s.db.NewSelect().Model((*PortAllocationsModel)(nil)).Column("port_prefix")
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Scan(ctx, &rows); err != nil {
		return nil, err
	}
	out := make([]int, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.PortPrefix)
	}
	return out, nil
}

func (s *Store) ListPortsvcServiceGroupChildrenByGroupIDs(ctx context.Context, ids []string) ([]PortsvcServiceGroupChildrenRecord, error) {
	children := make(map[string]PortsvcServiceGroupChildrenRecord, len(ids))
	for _, id := range ids {
		children[id] = PortsvcServiceGroupChildrenRecord{ServiceGroupID: id}
	}

	var rows []ServiceGroupsPortGroupsModel
	if err := s.db.NewSelect().
		Model(&rows).
		Relation("PortGroup").
		Relation("PortGroup.Host").
		Where("service_group_port_group.service_group_id IN (?)", bun.List(ids)).
		Order("port_group.port_prefix ASC").
		Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range rows {
		record := *utilPortsvcServiceGroupPortGroupRecordFromModel(&rows[idx])
		child := children[record.ServiceGroupID]
		child.PortGroups = append(child.PortGroups, record)
		children[record.ServiceGroupID] = child
	}

	out := make([]PortsvcServiceGroupChildrenRecord, 0, len(children))
	for _, id := range ids {
		out = append(out, children[id])
	}
	return out, nil
}

func (s *Store) ReplacePortsvcServiceGroupPortGroups(ctx context.Context, serviceGroupID string) error {
	_, err := s.db.NewDelete().Model((*ServiceGroupsPortGroupsModel)(nil)).Where("service_group_id = ?", serviceGroupID).Exec(ctx)
	return err
}

func (s *Store) AddPortsvcServiceGroupPortGroups(ctx context.Context, groups []PortsvcServiceGroupPortGroupRecord) error {
	if len(groups) == 0 {
		return nil
	}
	rows := make([]ServiceGroupsPortGroupsModel, 0, len(groups))
	for _, group := range groups {
		rows = append(rows, ServiceGroupsPortGroupsModel{
			ID:             idgen.NewUUID7(),
			ServiceGroupID: group.ServiceGroupID,
			PortGroupID:    group.PortGroupID,
			Role:           group.Role,
			Notes:          group.Notes,
			CreatedAt:      group.CreatedAt,
			UpdatedAt:      group.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
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

	var linkRows []PortGroupAssetLinksModel
	if err := s.db.NewSelect().
		Model(&linkRows).
		Relation("Asset").
		Where("asset_link.port_group_id IN (?)", bun.List(ids)).
		Order("asset_link.relation_type ASC", "asset.name ASC").
		Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range linkRows {
		record := *utilPortsvcPortGroupAssetLinkRecordFromModel(&linkRows[idx])
		child := children[record.PortGroupID]
		child.AssetLinks = append(child.AssetLinks, record)
		children[record.PortGroupID] = child
	}

	var repositoryLinkRows []PortGroupRepositoryLinksModel
	if err := s.db.NewSelect().
		Model(&repositoryLinkRows).
		Relation("Repository").
		Where("repository_link.port_group_id IN (?)", bun.List(ids)).
		Order("repository_link.relation_type ASC", "repository.full_name ASC").
		Scan(ctx); err != nil {
		return nil, err
	}
	for idx := range repositoryLinkRows {
		record := *utilPortsvcPortGroupRepositoryLinkRecordFromModel(&repositoryLinkRows[idx])
		child := children[record.PortGroupID]
		child.RepositoryLinks = append(child.RepositoryLinks, record)
		children[record.PortGroupID] = child
	}

	out := make([]PortsvcPortGroupChildrenRecord, 0, len(children))
	for _, id := range ids {
		out = append(out, children[id])
	}
	return out, nil
}

func (s *Store) ReplacePortsvcPortGroupChildren(ctx context.Context, portGroupID string) error {
	if _, err := s.db.NewDelete().Model((*PortGroupRepositoryLinksModel)(nil)).Where("port_group_id = ?", portGroupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*PortGroupAssetLinksModel)(nil)).Where("port_group_id = ?", portGroupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.NewDelete().Model((*ServicesModel)(nil)).Where("port_group_id = ?", portGroupID).Exec(ctx); err != nil {
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
		id := slot.ID
		if !IsUUID7(id) {
			id = idgen.NewUUID7()
		}
		rows = append(rows, ServicesModel{
			ID:            id,
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

func (s *Store) AddPortsvcPortGroupAssetLinks(ctx context.Context, links []PortsvcPortGroupAssetLinkRecord) error {
	if len(links) == 0 {
		return nil
	}
	rows := make([]PortGroupAssetLinksModel, 0, len(links))
	for _, link := range links {
		rows = append(rows, PortGroupAssetLinksModel{
			ID:           idgen.NewUUID7(),
			PortGroupID:  link.PortGroupID,
			PortSlotID:   link.PortSlotID,
			AssetID:      link.AssetID,
			RelationType: link.RelationType,
			Required:     link.Required,
			Notes:        link.Notes,
			CreatedAt:    link.CreatedAt,
			UpdatedAt:    link.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
}

func (s *Store) AddPortsvcPortGroupRepositoryLinks(ctx context.Context, links []PortsvcPortGroupRepositoryLinkRecord) error {
	if len(links) == 0 {
		return nil
	}
	rows := make([]PortGroupRepositoryLinksModel, 0, len(links))
	for _, link := range links {
		rows = append(rows, PortGroupRepositoryLinksModel{
			ID: idgen.NewUUID7(), PortGroupID: link.PortGroupID, PortSlotID: link.PortSlotID,
			RepositoryID: link.RepositoryID, RelationType: link.RelationType, Required: link.Required,
			Notes: link.Notes, CreatedAt: link.CreatedAt, UpdatedAt: link.UpdatedAt,
		})
	}
	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	return err
}

func (s *Store) ExistsPortsvcPortGroupRepositoryLink(ctx context.Context, portGroupID, repositoryID string) (bool, error) {
	count, err := s.db.NewSelect().Model((*PortGroupRepositoryLinksModel)(nil)).
		Where("port_group_id = ?", portGroupID).Where("repository_id = ?", repositoryID).Count(ctx)
	return count > 0, err
}
