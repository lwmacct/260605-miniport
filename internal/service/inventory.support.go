package service

import (
	"context"
	"encoding/csv"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

const (
	allocationPortMin       = 10000
	allocationPortMax       = 59999
	allocationPortMaxStart  = 59990
	allocationGroupSize     = 10
	defaultAllocationStatus = "planned"
	defaultSlotProtocol     = "tcp"
	defaultSlotStatus       = "empty"
	defaultComponentType    = "opensource"
	defaultRepositoryKind   = "source"
)

type InventoryPortGroup = repository.InventoryPortGroupRecord
type InventoryPortSlot = repository.InventoryPortSlotRecord
type InventoryProject = repository.InventoryProjectRecord
type InventoryComponent = repository.InventoryComponentRecord
type InventoryRepositoryRef = repository.InventoryRepositoryRefRecord

type InventoryActor struct {
	UserID   int64
	Username string
	Admin    bool
}

type PortSlotPayload struct {
	Port     int
	Name     string
	Protocol string
	Purpose  string
	Status   string
	Notes    string
}

type ProjectPayload struct {
	Name        string
	Description string
	Notes       string
}

type ComponentPayload struct {
	Name    string
	Type    string
	URL     string
	Version string
	Notes   string
}

type RepositoryPayload struct {
	ProjectID int64
	Name      string
	URL       string
	Kind      string
	Notes     string
}

type PortGroupPayload struct {
	UserID        int64
	PortStart     int
	PortEnd       int
	Name          string
	DindIP        string
	DindContainer string
	Status        string
	Owner         string
	Tags          string
	Notes         string
	Slots         []PortSlotPayload
	Projects      []ProjectPayload
	Components    []ComponentPayload
	Repositories  []RepositoryPayload
}

type PortGroupView struct {
	InventoryPortGroup

	Slots        []InventoryPortSlot
	Projects     []InventoryProject
	Components   []InventoryComponent
	Repositories []InventoryRepositoryRef
}

type PortGroupListParams struct {
	Actor       InventoryActor
	UserID      int64
	Query       string
	Sort        string
	Status      string
	ProjectName string
	DindIP      string
}

type PortGroupBatchUpdateInput struct {
	Actor  InventoryActor
	IDs    []int64
	Owner  *string
	Status *string
	Tags   *string
}

type PortGroupBatchDeleteInput struct {
	Actor InventoryActor
	IDs   []int64
}

var (
	ErrInventoryBadRequest = errors.New("inventory bad request")
	ErrInventoryNotFound   = errors.New("inventory not found")
)

type InventoryError struct {
	Kind    error
	Message string
}

func (e InventoryError) Error() string { return e.Message }
func (e InventoryError) Unwrap() error { return e.Kind }

func utilBadInventoryRequest(message string) error {
	return InventoryError{Kind: ErrInventoryBadRequest, Message: message}
}
func utilInventoryNotFound(message string) error {
	return InventoryError{Kind: ErrInventoryNotFound, Message: message}
}

func utilNowUTC() time.Time {
	return time.Now().UTC()
}

func utilVisibleUserID(actor InventoryActor, requestedUserID int64) int64 {
	if actor.Admin {
		return requestedUserID
	}
	return actor.UserID
}

func utilInventoryPortGroupFromPayload(ctx context.Context, store *repository.Store, actor InventoryActor, current *InventoryPortGroup, payload PortGroupPayload) (*InventoryPortGroup, error) {
	userID := actor.UserID
	if actor.Admin {
		if payload.UserID > 0 {
			userID = payload.UserID
		} else if current != nil {
			userID = current.UserID
		}
	}
	if userID <= 0 {
		return nil, utilBadInventoryRequest("userId is required")
	}
	if _, err := store.FetchUserByID(ctx, userID); err != nil {
		return nil, err
	}

	portStart := payload.PortStart
	if portStart == 0 {
		allocated, err := utilNextAvailablePortStart(ctx, store, userID, utilCurrentID(current))
		if err != nil {
			return nil, err
		}
		portStart = allocated
	}
	portEnd := portStart + allocationGroupSize - 1
	if payload.PortEnd > 0 && payload.PortEnd != portEnd {
		return nil, utilBadInventoryRequest("portEnd must equal portStart + 9")
	}

	group := &InventoryPortGroup{
		UserID:        userID,
		PortStart:     portStart,
		PortEnd:       portEnd,
		Name:          strings.TrimSpace(payload.Name),
		DindIP:        strings.TrimSpace(payload.DindIP),
		DindContainer: strings.TrimSpace(payload.DindContainer),
		Status:        strings.TrimSpace(payload.Status),
		Owner:         strings.TrimSpace(payload.Owner),
		Tags:          strings.TrimSpace(payload.Tags),
		Notes:         strings.TrimSpace(payload.Notes),
	}
	if group.Status == "" {
		group.Status = defaultAllocationStatus
	}
	if err := utilValidatePortGroup(ctx, store, utilCurrentID(current), group); err != nil {
		return nil, err
	}
	return group, nil
}

func utilCurrentID(group *InventoryPortGroup) int64 {
	if group == nil {
		return 0
	}
	return group.ID
}

func utilValidatePortGroup(ctx context.Context, store *repository.Store, currentID int64, group *InventoryPortGroup) error {
	if group.Name == "" {
		return utilBadInventoryRequest("name is required")
	}
	if group.PortStart < allocationPortMin || group.PortStart > allocationPortMaxStart {
		return utilBadInventoryRequest("portStart must be between 10000 and 59990")
	}
	if group.PortStart%allocationGroupSize != 0 {
		return utilBadInventoryRequest("portStart must align to a 10-port group")
	}
	if group.PortEnd != group.PortStart+allocationGroupSize-1 {
		return utilBadInventoryRequest("port group must contain exactly 10 ports")
	}
	if group.PortEnd > allocationPortMax {
		return utilBadInventoryRequest("portEnd must be less than or equal to 59999")
	}
	count, err := store.CountInventoryOverlappingPortGroups(ctx, currentID, group)
	if err != nil {
		return err
	}
	if count > 0 {
		return utilBadInventoryRequest("port group is already allocated for this user")
	}
	return nil
}

func utilNextAvailablePortStart(ctx context.Context, store *repository.Store, userID int64, excludeID int64) (int, error) {
	used, err := store.ListInventoryPortStartsByUser(ctx, userID, excludeID)
	if err != nil {
		return 0, err
	}
	usedSet := map[int]struct{}{}
	for _, portStart := range used {
		usedSet[portStart] = struct{}{}
	}
	for portStart := allocationPortMin; portStart <= allocationPortMaxStart; portStart += allocationGroupSize {
		if _, exists := usedSet[portStart]; !exists {
			return portStart, nil
		}
	}
	return 0, utilBadInventoryRequest("no available port groups for this user")
}

func utilReplacePortGroupChildren(ctx context.Context, store *repository.Store, groupID int64, payload PortGroupPayload, now time.Time, group InventoryPortGroup) error {
	if err := store.ReplaceInventoryPortGroupChildren(ctx, groupID); err != nil {
		return err
	}
	slots, err := utilSlotsFromPayload(groupID, payload, now, group)
	if err != nil {
		return err
	}
	if addSlotErr := store.AddInventoryPortSlots(ctx, slots); addSlotErr != nil {
		return addSlotErr
	}
	if addProjectErr := store.AddInventoryProjects(ctx, utilProjectsFromPayload(groupID, payload.Projects, now)); addProjectErr != nil {
		return addProjectErr
	}
	if addComponentErr := store.AddInventoryComponents(ctx, utilComponentsFromPayload(groupID, payload.Components, now)); addComponentErr != nil {
		return addComponentErr
	}
	repositories := utilRepositoriesFromPayload(groupID, payload.Repositories, now)
	return store.AddInventoryRepositoryRefs(ctx, repositories)
}

func utilSlotsFromPayload(groupID int64, payload PortGroupPayload, now time.Time, group InventoryPortGroup) ([]InventoryPortSlot, error) {
	byPort := map[int]PortSlotPayload{}
	for _, slot := range payload.Slots {
		if slot.Port < group.PortStart || slot.Port > group.PortEnd {
			return nil, utilBadInventoryRequest("slot port must be inside the port group range")
		}
		if _, exists := byPort[slot.Port]; exists {
			return nil, utilBadInventoryRequest("slot ports must be unique")
		}
		byPort[slot.Port] = slot
	}
	slots := make([]InventoryPortSlot, 0, allocationGroupSize)
	for port := group.PortStart; port <= group.PortEnd; port++ {
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

func utilProjectsFromPayload(groupID int64, payload []ProjectPayload, now time.Time) []InventoryProject {
	items := make([]InventoryProject, 0, len(payload))
	for _, item := range payload {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		items = append(items, InventoryProject{
			PortGroupID: groupID,
			Name:        name,
			Description: strings.TrimSpace(item.Description),
			Notes:       strings.TrimSpace(item.Notes),
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return items
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

func utilRepositoriesFromPayload(groupID int64, payload []RepositoryPayload, now time.Time) []InventoryRepositoryRef {
	items := make([]InventoryRepositoryRef, 0, len(payload))
	for _, item := range payload {
		name := strings.TrimSpace(item.Name)
		url := strings.TrimSpace(item.URL)
		if name == "" && url == "" {
			continue
		}
		kind := strings.TrimSpace(item.Kind)
		if kind == "" {
			kind = defaultRepositoryKind
		}
		items = append(items, InventoryRepositoryRef{
			PortGroupID: groupID,
			ProjectID:   item.ProjectID,
			Name:        name,
			URL:         url,
			Kind:        kind,
			Notes:       strings.TrimSpace(item.Notes),
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return items
}

func utilNormalizeBatchUpdateInput(input PortGroupBatchUpdateInput) (PortGroupBatchUpdateInput, error) {
	ids, err := utilNormalizeIDs(input.IDs)
	if err != nil {
		return input, err
	}
	input.IDs = ids
	if input.Owner != nil {
		value := strings.TrimSpace(*input.Owner)
		input.Owner = &value
	}
	if input.Status != nil {
		value := strings.TrimSpace(*input.Status)
		input.Status = &value
	}
	if input.Tags != nil {
		value := strings.TrimSpace(*input.Tags)
		input.Tags = &value
	}
	if input.Owner == nil && input.Status == nil && input.Tags == nil {
		return input, utilBadInventoryRequest("no changes provided")
	}
	return input, nil
}

func utilNormalizeBatchDeleteInput(input PortGroupBatchDeleteInput) (PortGroupBatchDeleteInput, error) {
	ids, err := utilNormalizeIDs(input.IDs)
	if err != nil {
		return input, err
	}
	input.IDs = ids
	return input, nil
}

func utilNormalizeIDs(ids []int64) ([]int64, error) {
	seen := map[int64]struct{}{}
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			return nil, utilBadInventoryRequest("ids must be positive")
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	if len(out) == 0 {
		return nil, utilBadInventoryRequest("ids are required")
	}
	return out, nil
}

func utilCSVBytes(records [][]string) ([]byte, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}
	return []byte(builder.String()), nil
}

func utilPortGroupProjects(items []InventoryProject) string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, item.Name)
	}
	return strings.Join(values, "; ")
}

func utilPortGroupComponents(items []InventoryComponent) string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		value := item.Name
		if item.Version != "" {
			value += " " + item.Version
		}
		values = append(values, value)
	}
	return strings.Join(values, "; ")
}

func utilPortGroupRepositories(items []InventoryRepositoryRef) string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		value := item.Name
		if item.URL != "" {
			value += " " + item.URL
		}
		values = append(values, value)
	}
	return strings.Join(values, "; ")
}

func utilPortGroupSlots(items []InventoryPortSlot) string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		value := strconv.Itoa(item.Port)
		if item.Name != "" {
			value += " " + item.Name
		}
		if item.Purpose != "" {
			value += " " + item.Purpose
		}
		values = append(values, value)
	}
	return strings.Join(values, "; ")
}
