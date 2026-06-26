package service

import (
	"context"
	"strings"
	"time"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

func utilInventoryPortGroupFromPayload(ctx context.Context, store *repository.Store, currentID int64, payload PortGroupPayload) (*InventoryPortGroup, error) {
	group := &InventoryPortGroup{
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
	if err := utilValidatePortGroup(ctx, store, currentID, group); err != nil {
		return nil, err
	}
	return group, nil
}

func utilValidatePortGroup(ctx context.Context, store *repository.Store, currentID int64, group *InventoryPortGroup) error {
	if group.HostID <= 0 {
		return utilBadInventoryRequest("hostId is required")
	}
	if group.ServiceName == "" {
		return utilBadInventoryRequest("serviceName is required")
	}
	if group.PortStart <= 0 || group.PortEnd <= 0 {
		return utilBadInventoryRequest("portStart and portEnd are required")
	}
	if group.PortEnd < group.PortStart {
		return utilBadInventoryRequest("portEnd must be greater than or equal to portStart")
	}
	if group.PortEnd-group.PortStart+1 != 10 {
		return utilBadInventoryRequest("port group must contain exactly 10 ports")
	}
	if _, err := store.FetchInventoryHostByID(ctx, group.HostID); err != nil {
		return err
	}
	count, err := store.CountInventoryOverlappingPortGroups(ctx, currentID, group)
	if err != nil {
		return err
	}
	if count > 0 {
		return utilBadInventoryRequest("port range overlaps an existing group on this host")
	}
	return nil
}

func (s *InventoryService) replacePortGroupChildren(ctx context.Context, groupID int64, payload PortGroupPayload, now time.Time) error {
	if err := s.store.ReplaceInventoryPortGroupChildren(ctx, groupID); err != nil {
		return err
	}
	slots, err := utilSlotsFromPayload(groupID, payload, now)
	if err != nil {
		return err
	}
	if err := s.store.AddInventoryPortSlots(ctx, slots); err != nil {
		return err
	}
	components := utilComponentsFromPayload(groupID, payload.Components, now)
	if err := s.store.AddInventoryComponents(ctx, components); err != nil {
		return err
	}
	repositories, err := utilRepositoriesFromPayload(groupID, payload.Repositories, now)
	if err != nil {
		return err
	}
	return s.store.AddInventoryRepositoryRefs(ctx, repositories)
}

func utilSlotsFromPayload(groupID int64, payload PortGroupPayload, now time.Time) ([]InventoryPortSlot, error) {
	byPort := map[int]PortSlotPayload{}
	for _, slot := range payload.Slots {
		if slot.Port < payload.PortStart || slot.Port > payload.PortEnd {
			return nil, utilBadInventoryRequest("slot port must be inside the port group range")
		}
		if _, exists := byPort[slot.Port]; exists {
			return nil, utilBadInventoryRequest("slot ports must be unique")
		}
		byPort[slot.Port] = slot
	}
	slots := make([]InventoryPortSlot, 0, payload.PortEnd-payload.PortStart+1)
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
		slots = append(slots, InventoryPortSlot{
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

func utilComponentsFromPayload(groupID int64, payload []ComponentPayload, now time.Time) []InventoryComponent {
	items := make([]InventoryComponent, 0, len(payload))
	for _, item := range payload {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		itemType := strings.TrimSpace(item.Type)
		if itemType == "" {
			itemType = defaultComponentType
		}
		items = append(items, InventoryComponent{
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
	return items
}

func utilRepositoriesFromPayload(groupID int64, payload []RepositoryPayload, now time.Time) ([]InventoryRepositoryRef, error) {
	items := make([]InventoryRepositoryRef, 0, len(payload))
	for _, item := range payload {
		name := strings.TrimSpace(item.Name)
		url := strings.TrimSpace(item.URL)
		if name == "" && url == "" {
			continue
		}
		if name == "" || url == "" {
			return nil, utilBadInventoryRequest("repository name and url are required together")
		}
		kind := strings.TrimSpace(item.Kind)
		if kind == "" {
			kind = defaultRepositoryKind
		}
		items = append(items, InventoryRepositoryRef{
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

func utilNormalizeBatchUpdateInput(input PortGroupBatchUpdateInput) (PortGroupBatchUpdateInput, error) {
	input.IDs = utilUniqueInt64s(input.IDs)
	if len(input.IDs) == 0 {
		return input, utilBadInventoryRequest("ids are required")
	}
	hasChanges := false
	if input.Status != nil {
		status := strings.TrimSpace(*input.Status)
		if status == "" {
			return input, utilBadInventoryRequest("status cannot be empty")
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
		return input, utilBadInventoryRequest("at least one field must be provided")
	}
	return input, nil
}

func utilNormalizeBatchDeleteInput(input PortGroupBatchDeleteInput) (PortGroupBatchDeleteInput, error) {
	input.IDs = utilUniqueInt64s(input.IDs)
	if len(input.IDs) == 0 {
		return input, utilBadInventoryRequest("ids are required")
	}
	return input, nil
}
