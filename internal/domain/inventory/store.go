package inventory

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/uptrace/bun"
)

type store struct {
	db bun.IDB
}

func newStore(db *bun.DB) *store {
	return &store{db: db}
}

func (st *store) runInTx(ctx context.Context, fn func(context.Context, *store) error) error {
	return st.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return fn(ctx, &store{db: tx})
	})
}

func (st *store) listHosts(ctx context.Context, params HostListParams) ([]Host, error) {
	var hosts []Host
	query := st.db.NewSelect().Model(&hosts)

	if environment := compactString(params.Environment); environment != "" {
		query = query.Where("environment = ?", environment)
	}
	if keyword := compactString(params.Query); keyword != "" {
		pattern := searchPattern(keyword)
		query = query.Where(
			joinSearchClauses([]string{"ip", "name", "network", "environment", "notes"}),
			joinSearchArgs(pattern, 5)...,
		)
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
	return hosts, nil
}

func (st *store) createHost(ctx context.Context, host *Host) error {
	_, err := st.db.NewInsert().Model(host).Exec(ctx)
	return err
}

func (st *store) updateHost(ctx context.Context, id int64, host *Host) error {
	res, err := st.db.NewUpdate().
		Model(host).
		Column("ip", "name", "network", "environment", "notes", "updated_at").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return huma.Error404NotFound("host not found")
	}
	return nil
}

func (st *store) deleteHost(ctx context.Context, id int64) error {
	count, err := st.db.NewSelect().Model((*PortGroup)(nil)).Where("host_id = ?", id).Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return huma.Error400BadRequest("host still has port groups")
	}

	res, err := st.db.NewDelete().Model((*Host)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return huma.Error404NotFound("host not found")
	}
	return nil
}

func (st *store) getHost(ctx context.Context, id int64) (*Host, error) {
	host := &Host{}
	if err := st.db.NewSelect().Model(host).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("host not found")
		}
		return nil, err
	}
	return host, nil
}

func (st *store) listPortGroups(ctx context.Context, params PortGroupListParams) ([]PortGroup, error) {
	var groups []PortGroup
	query := st.db.NewSelect().
		Model(&groups).
		Relation("Host")

	if params.HostID > 0 {
		query = query.Where("port_group.host_id = ?", params.HostID)
	}
	if status := compactString(params.Status); status != "" {
		query = query.Where("port_group.status = ?", status)
	}
	if keyword := compactString(params.Query); keyword != "" {
		pattern := searchPattern(keyword)
		query = query.Where(
			joinSearchClauses([]string{
				"port_group.service_name",
				"port_group.container_name",
				"port_group.dind_host",
				"port_group.owner",
				"port_group.tags",
				"port_group.notes",
				"host.ip",
				"host.name",
			}),
			joinSearchArgs(pattern, 8)...,
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
	return groups, nil
}

func (st *store) getPortGroupWithHost(ctx context.Context, id int64) (*PortGroup, error) {
	group := &PortGroup{}
	if err := st.db.NewSelect().Model(group).Relation("Host").Where("port_group.id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("port group not found")
		}
		return nil, err
	}
	return group, nil
}

func (st *store) createPortGroup(ctx context.Context, group *PortGroup) error {
	_, err := st.db.NewInsert().Model(group).Exec(ctx)
	return err
}

func (st *store) countPortGroupsByIDs(ctx context.Context, ids []int64) (int, error) {
	return st.db.NewSelect().
		Model((*PortGroup)(nil)).
		Where("id IN (?)", bun.List(ids)).
		Count(ctx)
}

func (st *store) updatePortGroup(ctx context.Context, id int64, group *PortGroup) error {
	res, err := st.db.NewUpdate().
		Model(group).
		Column("host_id", "port_start", "port_end", "service_name", "container_name", "dind_host", "status", "owner", "tags", "notes", "updated_at").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return huma.Error404NotFound("port group not found")
	}
	return nil
}

func (st *store) batchUpdatePortGroups(ctx context.Context, input PortGroupBatchUpdateInput, updatedAt time.Time) error {
	query := st.db.NewUpdate().
		Model((*PortGroup)(nil)).
		Set("updated_at = ?", updatedAt).
		Where("id IN (?)", bun.List(input.IDs))

	if input.Status != nil {
		query = query.Set("status = ?", *input.Status)
	}
	if input.Owner != nil {
		query = query.Set("owner = ?", *input.Owner)
	}
	if input.Tags != nil {
		query = query.Set("tags = ?", *input.Tags)
	}

	_, err := query.Exec(ctx)
	return err
}

func (st *store) deletePortGroup(ctx context.Context, id int64) error {
	if _, err := st.db.NewDelete().Model((*Repository)(nil)).Where("port_group_id = ?", id).Exec(ctx); err != nil {
		return err
	}
	if _, err := st.db.NewDelete().Model((*Component)(nil)).Where("port_group_id = ?", id).Exec(ctx); err != nil {
		return err
	}
	if _, err := st.db.NewDelete().Model((*PortSlot)(nil)).Where("port_group_id = ?", id).Exec(ctx); err != nil {
		return err
	}

	res, err := st.db.NewDelete().Model((*PortGroup)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return huma.Error404NotFound("port group not found")
	}
	return nil
}

func (st *store) deletePortGroups(ctx context.Context, ids []int64) error {
	if _, err := st.db.NewDelete().Model((*Repository)(nil)).Where("port_group_id IN (?)", bun.List(ids)).Exec(ctx); err != nil {
		return err
	}
	if _, err := st.db.NewDelete().Model((*Component)(nil)).Where("port_group_id IN (?)", bun.List(ids)).Exec(ctx); err != nil {
		return err
	}
	if _, err := st.db.NewDelete().Model((*PortSlot)(nil)).Where("port_group_id IN (?)", bun.List(ids)).Exec(ctx); err != nil {
		return err
	}
	if _, err := st.db.NewDelete().Model((*PortGroup)(nil)).Where("id IN (?)", bun.List(ids)).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (st *store) getPortGroup(ctx context.Context, id int64) (*PortGroup, error) {
	group := &PortGroup{}
	if err := st.db.NewSelect().Model(group).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("port group not found")
		}
		return nil, err
	}
	return group, nil
}

func (st *store) countOverlappingPortGroups(ctx context.Context, currentID int64, group *PortGroup) (int, error) {
	query := st.db.NewSelect().
		Model((*PortGroup)(nil)).
		Where("host_id = ?", group.HostID).
		Where("NOT (port_end < ? OR port_start > ?)", group.PortStart, group.PortEnd)
	if currentID > 0 {
		query = query.Where("id != ?", currentID)
	}
	return query.Count(ctx)
}

func (st *store) replaceChildren(ctx context.Context, groupID int64) error {
	if _, err := st.db.NewDelete().Model((*Repository)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := st.db.NewDelete().Model((*Component)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := st.db.NewDelete().Model((*PortSlot)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (st *store) insertPortSlots(ctx context.Context, slots []PortSlot) error {
	if len(slots) == 0 {
		return nil
	}
	_, err := st.db.NewInsert().Model(&slots).Exec(ctx)
	return err
}

func (st *store) insertComponents(ctx context.Context, components []Component) error {
	if len(components) == 0 {
		return nil
	}
	_, err := st.db.NewInsert().Model(&components).Exec(ctx)
	return err
}

func (st *store) insertRepositories(ctx context.Context, repositories []Repository) error {
	if len(repositories) == 0 {
		return nil
	}
	_, err := st.db.NewInsert().Model(&repositories).Exec(ctx)
	return err
}

func (st *store) getPortGroupView(ctx context.Context, id int64) (*PortGroupView, error) {
	group, err := st.getPortGroupWithHost(ctx, id)
	if err != nil {
		return nil, err
	}
	views, err := st.buildPortGroupViews(ctx, []PortGroup{*group})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (st *store) listPortGroupsByIDs(ctx context.Context, ids []int64) ([]PortGroup, error) {
	var groups []PortGroup
	if err := st.db.NewSelect().
		Model(&groups).
		Relation("Host").
		Where("port_group.id IN (?)", bun.List(ids)).
		Scan(ctx); err != nil {
		return nil, err
	}
	return groups, nil
}

func (st *store) buildPortGroupViews(ctx context.Context, groups []PortGroup) ([]PortGroupView, error) {
	views := make([]PortGroupView, len(groups))
	if len(groups) == 0 {
		return views, nil
	}

	ids := make([]int64, len(groups))
	for idx, group := range groups {
		ids[idx] = group.ID
		views[idx] = PortGroupView{PortGroup: group}
	}

	var slots []PortSlot
	if err := st.db.NewSelect().Model(&slots).Where("port_group_id IN (?)", bun.List(ids)).Order("port ASC").Scan(ctx); err != nil {
		return nil, err
	}
	var components []Component
	if err := st.db.NewSelect().Model(&components).Where("port_group_id IN (?)", bun.List(ids)).Order("name ASC").Scan(ctx); err != nil {
		return nil, err
	}
	var repositories []Repository
	if err := st.db.NewSelect().Model(&repositories).Where("port_group_id IN (?)", bun.List(ids)).Order("name ASC").Scan(ctx); err != nil {
		return nil, err
	}

	viewByID := make(map[int64]*PortGroupView, len(views))
	for idx := range views {
		viewByID[views[idx].ID] = &views[idx]
	}
	for _, slot := range slots {
		viewByID[slot.PortGroupID].Slots = append(viewByID[slot.PortGroupID].Slots, slot)
	}
	for _, component := range components {
		viewByID[component.PortGroupID].Components = append(viewByID[component.PortGroupID].Components, component)
	}
	for _, repo := range repositories {
		viewByID[repo.PortGroupID].Repositories = append(viewByID[repo.PortGroupID].Repositories, repo)
	}
	return views, nil
}
