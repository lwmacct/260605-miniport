package inventory

import (
	"context"
	"database/sql"
	"errors"

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

func (st *store) listHosts(ctx context.Context) ([]Host, error) {
	var hosts []Host
	if err := st.db.NewSelect().Model(&hosts).Order("ip ASC").Scan(ctx); err != nil {
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

func (st *store) listPortGroups(ctx context.Context) ([]PortGroup, error) {
	var groups []PortGroup
	if err := st.db.NewSelect().
		Model(&groups).
		Relation("Host").
		Order("host.ip ASC", "port_group.port_start ASC").
		Scan(ctx); err != nil {
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
	if err := st.db.NewSelect().Model(&slots).Where("port_group_id IN (?)", bun.In(ids)).Order("port ASC").Scan(ctx); err != nil {
		return nil, err
	}
	var components []Component
	if err := st.db.NewSelect().Model(&components).Where("port_group_id IN (?)", bun.In(ids)).Order("name ASC").Scan(ctx); err != nil {
		return nil, err
	}
	var repositories []Repository
	if err := st.db.NewSelect().Model(&repositories).Where("port_group_id IN (?)", bun.In(ids)).Order("name ASC").Scan(ctx); err != nil {
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
