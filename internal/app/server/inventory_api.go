package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/uptrace/bun"
)

const (
	defaultGroupStatus    = "planned"
	defaultSlotProtocol   = "tcp"
	defaultSlotStatus     = "empty"
	defaultComponentType  = "opensource"
	defaultRepositoryKind = "source"
)

type hostListOutput struct {
	Body []Host `json:"body"`
}

type hostOutput struct {
	Body Host `json:"body"`
}

type hostInput struct {
	ID int64 `path:"id" example:"1"`
}

type hostBodyInput struct {
	Body hostPayload
}

type hostUpdateInput struct {
	ID   int64 `path:"id" example:"1"`
	Body hostPayload
}

type portGroupListOutput struct {
	Body []portGroupView `json:"body"`
}

type portGroupOutput struct {
	Body portGroupView `json:"body"`
}

type portGroupInput struct {
	ID int64 `path:"id" example:"1"`
}

type portGroupBodyInput struct {
	Body portGroupPayload
}

type portGroupUpdateInput struct {
	ID   int64 `path:"id" example:"1"`
	Body portGroupPayload
}

type deleteOutput struct {
	Body struct {
		Deleted bool `json:"deleted" example:"true"`
	}
}

func migrateInventory(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		return fmt.Errorf("enable sqlite foreign keys: %w", err)
	}

	models := []any{
		(*Host)(nil),
		(*PortGroup)(nil),
		(*PortSlot)(nil),
		(*Component)(nil),
		(*Repository)(nil),
	}
	for _, model := range models {
		if _, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			return fmt.Errorf("create inventory table: %w", err)
		}
	}

	statements := []string{
		`CREATE INDEX IF NOT EXISTS idx_port_groups_host ON port_groups(host_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_port_slots_group_port ON port_slots(port_group_id, port)`,
		`CREATE INDEX IF NOT EXISTS idx_components_group ON components(port_group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_repositories_group ON repositories(port_group_id)`,
	}
	for _, stmt := range statements {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("create inventory index: %w", err)
		}
	}
	return nil
}

func registerInventoryRoutes(api huma.API, db *bun.DB) {
	huma.Register(api, huma.Operation{
		OperationID: "list-hosts",
		Method:      http.MethodGet,
		Path:        "/hosts",
		Summary:     "List IP hosts",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, _ *struct{}) (*hostListOutput, error) {
		var hosts []Host
		if err := db.NewSelect().Model(&hosts).Order("ip ASC").Scan(ctx); err != nil {
			return nil, err
		}
		return &hostListOutput{Body: hosts}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-host",
		Method:      http.MethodPost,
		Path:        "/hosts",
		Summary:     "Create IP host",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *hostBodyInput) (*hostOutput, error) {
		host, err := createHost(ctx, db, input.Body)
		if err != nil {
			return nil, err
		}
		return &hostOutput{Body: *host}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-host",
		Method:      http.MethodPut,
		Path:        "/hosts/{id}",
		Summary:     "Update IP host",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *hostUpdateInput) (*hostOutput, error) {
		host, err := updateHost(ctx, db, input.ID, input.Body)
		if err != nil {
			return nil, err
		}
		return &hostOutput{Body: *host}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-host",
		Method:      http.MethodDelete,
		Path:        "/hosts/{id}",
		Summary:     "Delete IP host",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *hostInput) (*deleteOutput, error) {
		if err := deleteHost(ctx, db, input.ID); err != nil {
			return nil, err
		}
		out := &deleteOutput{}
		out.Body.Deleted = true
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-port-groups",
		Method:      http.MethodGet,
		Path:        "/port-groups",
		Summary:     "List port groups",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, _ *struct{}) (*portGroupListOutput, error) {
		groups, err := listPortGroups(ctx, db)
		if err != nil {
			return nil, err
		}
		return &portGroupListOutput{Body: groups}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-port-group",
		Method:      http.MethodGet,
		Path:        "/port-groups/{id}",
		Summary:     "Get port group",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *portGroupInput) (*portGroupOutput, error) {
		group, err := getPortGroupView(ctx, db, input.ID)
		if err != nil {
			return nil, err
		}
		return &portGroupOutput{Body: *group}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-port-group",
		Method:      http.MethodPost,
		Path:        "/port-groups",
		Summary:     "Create port group",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *portGroupBodyInput) (*portGroupOutput, error) {
		group, err := createPortGroup(ctx, db, input.Body)
		if err != nil {
			return nil, err
		}
		return &portGroupOutput{Body: *group}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-port-group",
		Method:      http.MethodPut,
		Path:        "/port-groups/{id}",
		Summary:     "Update port group",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *portGroupUpdateInput) (*portGroupOutput, error) {
		group, err := updatePortGroup(ctx, db, input.ID, input.Body)
		if err != nil {
			return nil, err
		}
		return &portGroupOutput{Body: *group}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-port-group",
		Method:      http.MethodDelete,
		Path:        "/port-groups/{id}",
		Summary:     "Delete port group",
		Tags:        []string{"inventory"},
	}, func(ctx context.Context, input *portGroupInput) (*deleteOutput, error) {
		if err := deletePortGroup(ctx, db, input.ID); err != nil {
			return nil, err
		}
		out := &deleteOutput{}
		out.Body.Deleted = true
		return out, nil
	})
}

func createHost(ctx context.Context, db *bun.DB, payload hostPayload) (*Host, error) {
	host := hostFromPayload(payload)
	if err := validateHost(host); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	host.CreatedAt = now
	host.UpdatedAt = now
	if _, err := db.NewInsert().Model(host).Exec(ctx); err != nil {
		return nil, err
	}
	return host, nil
}

func updateHost(ctx context.Context, db *bun.DB, id int64, payload hostPayload) (*Host, error) {
	host := hostFromPayload(payload)
	host.ID = id
	if err := validateHost(host); err != nil {
		return nil, err
	}
	host.UpdatedAt = time.Now().UTC()
	res, err := db.NewUpdate().
		Model(host).
		Column("ip", "name", "network", "environment", "notes", "updated_at").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return nil, huma.Error404NotFound("host not found")
	}
	return getHost(ctx, db, id)
}

func deleteHost(ctx context.Context, db *bun.DB, id int64) error {
	count, err := db.NewSelect().Model((*PortGroup)(nil)).Where("host_id = ?", id).Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return huma.Error400BadRequest("host still has port groups")
	}
	res, err := db.NewDelete().Model((*Host)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return huma.Error404NotFound("host not found")
	}
	return nil
}

func getHost(ctx context.Context, db bun.IDB, id int64) (*Host, error) {
	host := &Host{}
	if err := db.NewSelect().Model(host).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("host not found")
		}
		return nil, err
	}
	return host, nil
}

func hostFromPayload(payload hostPayload) *Host {
	return &Host{
		IP:          strings.TrimSpace(payload.IP),
		Name:        strings.TrimSpace(payload.Name),
		Network:     strings.TrimSpace(payload.Network),
		Environment: strings.TrimSpace(payload.Environment),
		Notes:       strings.TrimSpace(payload.Notes),
	}
}

func validateHost(host *Host) error {
	if host.IP == "" {
		return huma.Error400BadRequest("host ip is required")
	}
	return nil
}

func listPortGroups(ctx context.Context, db *bun.DB) ([]portGroupView, error) {
	var groups []PortGroup
	if err := db.NewSelect().
		Model(&groups).
		Relation("Host").
		Order("host.ip ASC", "port_group.port_start ASC").
		Scan(ctx); err != nil {
		return nil, err
	}
	return buildPortGroupViews(ctx, db, groups)
}

func getPortGroupView(ctx context.Context, db bun.IDB, id int64) (*portGroupView, error) {
	group := &PortGroup{}
	if err := db.NewSelect().Model(group).Relation("Host").Where("port_group.id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("port group not found")
		}
		return nil, err
	}
	views, err := buildPortGroupViews(ctx, db, []PortGroup{*group})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func createPortGroup(ctx context.Context, db *bun.DB, payload portGroupPayload) (*portGroupView, error) {
	var out *portGroupView
	err := db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		group, err := portGroupFromPayload(ctx, tx, 0, payload)
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		group.CreatedAt = now
		group.UpdatedAt = now
		if _, err := tx.NewInsert().Model(group).Exec(ctx); err != nil {
			return err
		}
		if err := replacePortGroupChildren(ctx, tx, group.ID, payload, now); err != nil {
			return err
		}
		view, err := getPortGroupView(ctx, tx, group.ID)
		if err != nil {
			return err
		}
		out = view
		return nil
	})
	return out, err
}

func updatePortGroup(ctx context.Context, db *bun.DB, id int64, payload portGroupPayload) (*portGroupView, error) {
	var out *portGroupView
	err := db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := getPortGroup(ctx, tx, id); err != nil {
			return err
		}
		group, err := portGroupFromPayload(ctx, tx, id, payload)
		if err != nil {
			return err
		}
		group.ID = id
		now := time.Now().UTC()
		group.UpdatedAt = now
		res, err := tx.NewUpdate().
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
		if err := replacePortGroupChildren(ctx, tx, id, payload, now); err != nil {
			return err
		}
		view, err := getPortGroupView(ctx, tx, id)
		if err != nil {
			return err
		}
		out = view
		return nil
	})
	return out, err
}

func deletePortGroup(ctx context.Context, db *bun.DB, id int64) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model((*Repository)(nil)).Where("port_group_id = ?", id).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*Component)(nil)).Where("port_group_id = ?", id).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*PortSlot)(nil)).Where("port_group_id = ?", id).Exec(ctx); err != nil {
			return err
		}
		res, err := tx.NewDelete().Model((*PortGroup)(nil)).Where("id = ?", id).Exec(ctx)
		if err != nil {
			return err
		}
		if rows, _ := res.RowsAffected(); rows == 0 {
			return huma.Error404NotFound("port group not found")
		}
		return nil
	})
}

func getPortGroup(ctx context.Context, db bun.IDB, id int64) (*PortGroup, error) {
	group := &PortGroup{}
	if err := db.NewSelect().Model(group).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("port group not found")
		}
		return nil, err
	}
	return group, nil
}

func portGroupFromPayload(ctx context.Context, db bun.IDB, currentID int64, payload portGroupPayload) (*PortGroup, error) {
	group := &PortGroup{
		HostID:        payload.HostID,
		PortStart:     payload.PortStart,
		PortEnd:       payload.PortEnd,
		ServiceName:   strings.TrimSpace(payload.ServiceName),
		ContainerName: strings.TrimSpace(payload.ContainerName),
		DindHost:      strings.TrimSpace(payload.DindHost),
		Status:        strings.TrimSpace(payload.Status),
		Owner:         strings.TrimSpace(payload.Owner),
		Tags:          strings.TrimSpace(payload.Tags),
		Notes:         strings.TrimSpace(payload.Notes),
	}
	if group.Status == "" {
		group.Status = defaultGroupStatus
	}
	if err := validatePortGroup(ctx, db, currentID, group); err != nil {
		return nil, err
	}
	return group, nil
}

func validatePortGroup(ctx context.Context, db bun.IDB, currentID int64, group *PortGroup) error {
	if group.HostID <= 0 {
		return huma.Error400BadRequest("hostId is required")
	}
	if group.ServiceName == "" {
		return huma.Error400BadRequest("serviceName is required")
	}
	if group.PortStart <= 0 || group.PortEnd <= 0 {
		return huma.Error400BadRequest("portStart and portEnd are required")
	}
	if group.PortEnd < group.PortStart {
		return huma.Error400BadRequest("portEnd must be greater than or equal to portStart")
	}
	if group.PortEnd-group.PortStart+1 != 10 {
		return huma.Error400BadRequest("port group must contain exactly 10 ports")
	}
	if _, err := getHost(ctx, db, group.HostID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return huma.Error404NotFound("host not found")
	}
	countQuery := db.NewSelect().
		Model((*PortGroup)(nil)).
		Where("host_id = ?", group.HostID).
		Where("NOT (port_end < ? OR port_start > ?)", group.PortStart, group.PortEnd)
	if currentID > 0 {
		countQuery = countQuery.Where("id != ?", currentID)
	}
	count, err := countQuery.Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return huma.Error400BadRequest("port range overlaps an existing group on this host")
	}
	return nil
}

func replacePortGroupChildren(ctx context.Context, db bun.IDB, groupID int64, payload portGroupPayload, now time.Time) error {
	if _, err := db.NewDelete().Model((*Repository)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := db.NewDelete().Model((*Component)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}
	if _, err := db.NewDelete().Model((*PortSlot)(nil)).Where("port_group_id = ?", groupID).Exec(ctx); err != nil {
		return err
	}

	slots, err := slotsFromPayload(groupID, payload, now)
	if err != nil {
		return err
	}
	if len(slots) > 0 {
		if _, err := db.NewInsert().Model(&slots).Exec(ctx); err != nil {
			return err
		}
	}

	components, err := componentsFromPayload(groupID, payload.Components, now)
	if err != nil {
		return err
	}
	if len(components) > 0 {
		if _, err := db.NewInsert().Model(&components).Exec(ctx); err != nil {
			return err
		}
	}

	repositories, err := repositoriesFromPayload(groupID, payload.Repositories, now)
	if err != nil {
		return err
	}
	if len(repositories) > 0 {
		if _, err := db.NewInsert().Model(&repositories).Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

func slotsFromPayload(groupID int64, payload portGroupPayload, now time.Time) ([]PortSlot, error) {
	byPort := map[int]portSlotPayload{}
	for _, slot := range payload.Slots {
		if slot.Port < payload.PortStart || slot.Port > payload.PortEnd {
			return nil, huma.Error400BadRequest("slot port must be inside the port group range")
		}
		if _, ok := byPort[slot.Port]; ok {
			return nil, huma.Error400BadRequest("slot ports must be unique")
		}
		byPort[slot.Port] = slot
	}

	slots := make([]PortSlot, 0, payload.PortEnd-payload.PortStart+1)
	for port := payload.PortStart; port <= payload.PortEnd; port++ {
		source := byPort[port]
		protocol := strings.TrimSpace(source.Protocol)
		if protocol == "" {
			protocol = defaultSlotProtocol
		}
		status := strings.TrimSpace(source.Status)
		if status == "" {
			status = defaultSlotStatus
		}
		slots = append(slots, PortSlot{
			PortGroupID: groupID,
			Port:        port,
			Name:        strings.TrimSpace(source.Name),
			Protocol:    protocol,
			Purpose:     strings.TrimSpace(source.Purpose),
			Status:      status,
			Notes:       strings.TrimSpace(source.Notes),
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return slots, nil
}

func componentsFromPayload(groupID int64, payload []componentPayload, now time.Time) ([]Component, error) {
	components := make([]Component, 0, len(payload))
	for _, item := range payload {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		itemType := strings.TrimSpace(item.Type)
		if itemType == "" {
			itemType = defaultComponentType
		}
		components = append(components, Component{
			PortGroupID: groupID,
			Name:        name,
			Type:        itemType,
			URL:         strings.TrimSpace(item.URL),
			Version:     strings.TrimSpace(item.Version),
			Notes:       strings.TrimSpace(item.Notes),
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return components, nil
}

func repositoriesFromPayload(groupID int64, payload []repositoryPayload, now time.Time) ([]Repository, error) {
	repositories := make([]Repository, 0, len(payload))
	for _, item := range payload {
		name := strings.TrimSpace(item.Name)
		url := strings.TrimSpace(item.URL)
		if name == "" && url == "" {
			continue
		}
		if name == "" || url == "" {
			return nil, huma.Error400BadRequest("repository name and url are required together")
		}
		kind := strings.TrimSpace(item.Kind)
		if kind == "" {
			kind = defaultRepositoryKind
		}
		repositories = append(repositories, Repository{
			PortGroupID: groupID,
			Name:        name,
			URL:         url,
			Kind:        kind,
			Notes:       strings.TrimSpace(item.Notes),
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return repositories, nil
}

func buildPortGroupViews(ctx context.Context, db bun.IDB, groups []PortGroup) ([]portGroupView, error) {
	views := make([]portGroupView, len(groups))
	if len(groups) == 0 {
		return views, nil
	}

	ids := make([]int64, len(groups))
	for idx, group := range groups {
		ids[idx] = group.ID
		views[idx] = portGroupView{PortGroup: group}
	}

	var slots []PortSlot
	if err := db.NewSelect().Model(&slots).Where("port_group_id IN (?)", bun.In(ids)).Order("port ASC").Scan(ctx); err != nil {
		return nil, err
	}
	var components []Component
	if err := db.NewSelect().Model(&components).Where("port_group_id IN (?)", bun.In(ids)).Order("name ASC").Scan(ctx); err != nil {
		return nil, err
	}
	var repositories []Repository
	if err := db.NewSelect().Model(&repositories).Where("port_group_id IN (?)", bun.In(ids)).Order("name ASC").Scan(ctx); err != nil {
		return nil, err
	}

	viewByID := map[int64]*portGroupView{}
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
