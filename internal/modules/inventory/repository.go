package inventory

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/uptrace/bun"
)

type repository struct {
	db bun.IDB
}

func newRepository(db *bun.DB) *repository {
	return &repository{db: db}
}

func (r *repository) runInTx(ctx context.Context, fn func(context.Context, *repository) error) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return fn(ctx, &repository{db: tx})
	})
}

func (r *repository) listHosts(ctx context.Context) ([]Host, error) {
	var hosts []Host
	if err := r.db.NewSelect().Model(&hosts).Order("ip ASC").Scan(ctx); err != nil {
		return nil, err
	}
	return hosts, nil
}

func (r *repository) createHost(ctx context.Context, host *Host) error {
	_, err := r.db.NewInsert().Model(host).Exec(ctx)
	return err
}

func (r *repository) updateHost(ctx context.Context, id int64, host *Host) error {
	res, err := r.db.NewUpdate().
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

func (r *repository) deleteHost(ctx context.Context, id int64) error {
	count, err := r.db.NewSelect().Model((*PortGroup)(nil)).Where("host_id = ?", id).Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return huma.Error400BadRequest("host still has port groups")
	}
	res, err := r.db.NewDelete().Model((*Host)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return huma.Error404NotFound("host not found")
	}
	return nil
}

func (r *repository) getHost(ctx context.Context, id int64) (*Host, error) {
	host := &Host{}
	if err := r.db.NewSelect().Model(host).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("host not found")
		}
		return nil, err
	}
	return host, nil
}

func (r *repository) listPortGroups(ctx context.Context) ([]PortGroup, error) {
	var groups []PortGroup
	if err := r.db.NewSelect().
		Model(&groups).
		Relation("Host").
		Order("host.ip ASC", "port_group.port_start ASC").
		Scan(ctx); err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *repository) getPortGroupWithHost(ctx context.Context, id int64) (*PortGroup, error) {
	group := &PortGroup{}
	if err := r.db.NewSelect().Model(group).Relation("Host").Where("port_group.id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("port group not found")
		}
		return nil, err
	}
	return group, nil
}

func (r *repository) createPortGroup(ctx context.Context, group *PortGroup) error {
	_, err := r.db.NewInsert().Model(group).Exec(ctx)
	return err
}

func (r *repository) updatePortGroup(ctx context.Context, id int64, group *PortGroup) error {
	res, err := r.db.NewUpdate().
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

func (r *repository) deletePortGroup(ctx context.Context, id int64) error {
	if _, err := r.db.NewDelete().Model((*Repository)(nil)).Where("port_group_id = ?", id).Exec(ctx); err != nil {
		return err
	}
	if _, err := r.db.NewDelete().Model((*Component)(nil)).Where("port_group_id = ?", id).Exec(ctx); err != nil {
		return err
	}
	if _, err := r.db.NewDelete().Model((*PortSlot)(nil)).Where("port_group_id = ?", id).Exec(ctx); err != nil {
		return err
	}
	res, err := r.db.NewDelete().Model((*PortGroup)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return huma.Error404NotFound("port group not found")
	}
	return nil
}

func (r *repository) getPortGroup(ctx context.Context, id int64) (*PortGroup, error) {
	group := &PortGroup{}
	if err := r.db.NewSelect().Model(group).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("port group not found")
		}
		return nil, err
	}
	return group, nil
}

func (r *repository) insertPortSlots(ctx context.Context, slots []PortSlot) error {
	if len(slots) == 0 {
		return nil
	}
	_, err := r.db.NewInsert().Model(&slots).Exec(ctx)
	return err
}

func (r *repository) insertComponents(ctx context.Context, components []Component) error {
	if len(components) == 0 {
		return nil
	}
	_, err := r.db.NewInsert().Model(&components).Exec(ctx)
	return err
}

func (r *repository) insertRepositories(ctx context.Context, repositories []Repository) error {
	if len(repositories) == 0 {
		return nil
	}
	_, err := r.db.NewInsert().Model(&repositories).Exec(ctx)
	return err
}

func (r *repository) replaceChildren(ctx context.Context, groupID int64) error {
	if _, err := r.db.NewDelete().Model((*Repository)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := r.db.NewDelete().Model((*Component)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := r.db.NewDelete().Model((*PortSlot)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *repository) countOverlappingPortGroups(ctx context.Context, currentID int64, group *PortGroup) (int, error) {
	countQuery := r.db.NewSelect().
		Model((*PortGroup)(nil)).
		Where("host_id = ?", group.HostID).
		Where("NOT (port_end < ? OR port_start > ?)", group.PortStart, group.PortEnd)
	if currentID > 0 {
		countQuery = countQuery.Where("id != ?", currentID)
	}
	return countQuery.Count(ctx)
}

func (r *repository) getPortGroupView(ctx context.Context, id int64) (*PortGroupView, error) {
	group, err := r.getPortGroupWithHost(ctx, id)
	if err != nil {
		return nil, err
	}
	views, err := r.buildPortGroupViews(ctx, []PortGroup{*group})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (r *repository) buildPortGroupViews(ctx context.Context, groups []PortGroup) ([]PortGroupView, error) {
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
	if err := r.db.NewSelect().Model(&slots).Where("port_group_id IN (?)", bun.In(ids)).Order("port ASC").Scan(ctx); err != nil {
		return nil, err
	}
	var components []Component
	if err := r.db.NewSelect().Model(&components).Where("port_group_id IN (?)", bun.In(ids)).Order("name ASC").Scan(ctx); err != nil {
		return nil, err
	}
	var repositories []Repository
	if err := r.db.NewSelect().Model(&repositories).Where("port_group_id IN (?)", bun.In(ids)).Order("name ASC").Scan(ctx); err != nil {
		return nil, err
	}

	viewByID := map[int64]*PortGroupView{}
	for idx := range views {
		viewByID[views[idx].ID] = &views[idx]
	}
	for _, slot := range slots {
		viewByID[slot.PortGroupID].Slots = append(viewByID[slot.PortGroupID].Slots, slot)
	}
	for _, component := range components {
		viewByID[component.PortGroupID].Components = append(viewByID[component.PortGroupID].Components, component)
	}
	for _, repository := range repositories {
		viewByID[repository.PortGroupID].Repositories = append(viewByID[repository.PortGroupID].Repositories, repository)
	}
	return views, nil
}
