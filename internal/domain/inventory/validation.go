package inventory

import (
	"context"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

const (
	defaultGroupStatus    = "planned"
	defaultSlotProtocol   = "tcp"
	defaultSlotStatus     = "empty"
	defaultComponentType  = "opensource"
	defaultRepositoryKind = "source"
)

func hostFromPayload(payload HostPayload) *Host {
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

func portGroupFromPayload(ctx context.Context, st *store, currentID int64, payload PortGroupPayload) (*PortGroup, error) {
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
	if err := validatePortGroup(ctx, st, currentID, group); err != nil {
		return nil, err
	}
	return group, nil
}

func validatePortGroup(ctx context.Context, st *store, currentID int64, group *PortGroup) error {
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
	if _, err := st.getHost(ctx, group.HostID); err != nil {
		return err
	}
	count, err := st.countOverlappingPortGroups(ctx, currentID, group)
	if err != nil {
		return err
	}
	if count > 0 {
		return huma.Error400BadRequest("port range overlaps an existing group on this host")
	}
	return nil
}

func replacePortGroupChildren(ctx context.Context, st *store, groupID int64, payload PortGroupPayload, now time.Time) error {
	if err := st.replaceChildren(ctx, groupID); err != nil {
		return err
	}

	slots, err := slotsFromPayload(groupID, payload, now)
	if err != nil {
		return err
	}
	if err := st.insertPortSlots(ctx, slots); err != nil {
		return err
	}

	components, err := componentsFromPayload(groupID, payload.Components, now)
	if err != nil {
		return err
	}
	if err := st.insertComponents(ctx, components); err != nil {
		return err
	}

	repositories, err := repositoriesFromPayload(groupID, payload.Repositories, now)
	if err != nil {
		return err
	}
	return st.insertRepositories(ctx, repositories)
}

func slotsFromPayload(groupID int64, payload PortGroupPayload, now time.Time) ([]PortSlot, error) {
	byPort := map[int]PortSlotPayload{}
	for _, slot := range payload.Slots {
		if slot.Port < payload.PortStart || slot.Port > payload.PortEnd {
			return nil, huma.Error400BadRequest("slot port must be inside the port group range")
		}
		if _, exists := byPort[slot.Port]; exists {
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

func componentsFromPayload(groupID int64, payload []ComponentPayload, now time.Time) ([]Component, error) {
	items := make([]Component, 0, len(payload))
	for _, item := range payload {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}

		itemType := strings.TrimSpace(item.Type)
		if itemType == "" {
			itemType = defaultComponentType
		}
		items = append(items, Component{
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
	return items, nil
}

func repositoriesFromPayload(groupID int64, payload []RepositoryPayload, now time.Time) ([]Repository, error) {
	items := make([]Repository, 0, len(payload))
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
		items = append(items, Repository{
			PortGroupID: groupID,
			Name:        name,
			URL:         url,
			Kind:        kind,
			Notes:       strings.TrimSpace(item.Notes),
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return items, nil
}

func normalizeBatchUpdateInput(input PortGroupBatchUpdateInput) (PortGroupBatchUpdateInput, error) {
	input.IDs = uniqueInt64s(input.IDs)
	if len(input.IDs) == 0 {
		return input, huma.Error400BadRequest("ids are required")
	}

	hasChanges := false
	if input.Status != nil {
		status := strings.TrimSpace(*input.Status)
		if status == "" {
			return input, huma.Error400BadRequest("status cannot be empty")
		}
		input.Status = &status
		hasChanges = true
	}
	if input.Owner != nil {
		owner := strings.TrimSpace(*input.Owner)
		input.Owner = &owner
		hasChanges = true
	}
	if input.Tags != nil {
		tags := strings.TrimSpace(*input.Tags)
		input.Tags = &tags
		hasChanges = true
	}
	if !hasChanges {
		return input, huma.Error400BadRequest("at least one field must be provided")
	}

	return input, nil
}

func normalizeBatchDeleteInput(input PortGroupBatchDeleteInput) (PortGroupBatchDeleteInput, error) {
	input.IDs = uniqueInt64s(input.IDs)
	if len(input.IDs) == 0 {
		return input, huma.Error400BadRequest("ids are required")
	}
	return input, nil
}
