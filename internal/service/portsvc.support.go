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
	allocationPortPrefixMin   = 1000
	allocationPortPrefixMax   = 5999
	allocationGroupSize       = 10
	defaultHostStatus         = "active"
	hostStatusStopped         = "stopped"
	defaultPortGroupStatus    = "available"
	defaultRuntimeMode        = "dind"
	defaultSlotKind           = "app"
	defaultSlotProtocol       = "tcp"
	defaultSlotStatus         = "planned"
	defaultAssetKind          = "component"
	defaultAssetType          = "middleware"
	defaultAssetProvider      = "manual"
	defaultAssetVisibility    = "unknown"
	defaultControllability    = "unknown"
	defaultAssetStatus        = "active"
	defaultRelationType       = "runtime"
	defaultServiceGroupKind   = "service"
	defaultServiceGroupStatus = "active"
	runtimeModeDind           = "dind"
	runtimeModeHost           = "host"
)

type Host = repository.PortsvcHostRecord
type PortGroup = repository.PortsvcPortGroupRecord
type PortSlot = repository.PortsvcPortSlotRecord
type DependencyAsset = repository.PortsvcDependencyAssetRecord
type PortGroupAssetLink = repository.PortsvcPortGroupAssetLinkRecord
type PortGroupRepositoryLink = repository.PortsvcPortGroupRepositoryLinkRecord
type ServiceGroup = repository.PortsvcServiceGroupRecord
type ServiceGroupPortGroup = repository.PortsvcServiceGroupPortGroupRecord

type HostPayload struct {
	Name   string
	IP     string
	Spec   string
	Status string
	Notes  string
}

type HostUpdateInput struct {
	ID      string
	Payload HostPayload
}

type PortSlotPayload struct {
	ID            string
	Port          int
	Name          string
	Kind          string
	Protocol      string
	ContainerName string
	Status        string
	Notes         string
}

type DependencyAssetPayload struct {
	Name            string
	AssetKind       string
	AssetType       string
	Provider        string
	URL             string
	FullName        string
	ExternalID      string
	Visibility      string
	Controllability string
	Status          string
	Description     string
	Metadata        string
	Notes           string
}

type DependencyAssetUpdateInput struct {
	ID      string
	Payload DependencyAssetPayload
}

type PortGroupAssetLinkPayload struct {
	ID           string
	PortSlotID   string
	AssetID      string
	RelationType string
	Required     bool
	Notes        string
}

type PortGroupRepositoryLinkPayload struct {
	ID           string
	PortSlotID   string
	RepositoryID string
	RelationType string
	Required     bool
	Notes        string
}

type PortGroupPayload struct {
	HostID           string
	PortPrefix       int
	EnvironmentName  string
	EnvironmentOwner string
	RuntimeMode      string
	RuntimeName      string
	ServiceIP        string
	Status           string
	Tags             string
	Notes            string
	Slots            []PortSlotPayload
	AssetLinks       []PortGroupAssetLinkPayload
	RepositoryLinks  []PortGroupRepositoryLinkPayload
}

type PortGroupUpdateInput struct {
	ID      string
	Payload PortGroupPayload
}

type ServiceGroupPortGroupPayload struct {
	ID          string
	PortGroupID string
	Role        string
	Notes       string
}

type ServiceGroupPayload struct {
	Name        string
	Kind        string
	Status      string
	Description string
	Notes       string
	PortGroups  []ServiceGroupPortGroupPayload
}

type ServiceGroupUpdateInput struct {
	ID      string
	Payload ServiceGroupPayload
}

type PortGroupView struct {
	PortGroup

	Slots           []PortSlot
	AssetLinks      []PortGroupAssetLink
	RepositoryLinks []PortGroupRepositoryLink
}

type ServiceGroupView struct {
	ServiceGroup

	PortGroups []ServiceGroupPortGroup
}

type HostListParams struct {
	Query  string
	Status string
}

type PortGroupListParams struct {
	Query  string
	Sort   string
	Status string
}

type DependencyAssetListParams struct {
	Query     string
	AssetKind string
	AssetType string
	Provider  string
	Status    string
}

type ServiceGroupListParams struct {
	Query  string
	Status string
}

var (
	ErrPortsvcBadRequest = errors.New("portsvc bad request")
	ErrPortsvcConflict   = errors.New("portsvc conflict")
	ErrPortsvcNotFound   = errors.New("portsvc not found")
)

type PortsvcError struct {
	Kind    error
	Message string
}

func (e PortsvcError) Error() string { return e.Message }
func (e PortsvcError) Unwrap() error { return e.Kind }

func utilBadPortsvcRequest(message string) error {
	return PortsvcError{Kind: ErrPortsvcBadRequest, Message: message}
}
func utilPortsvcNotFound(message string) error {
	return PortsvcError{Kind: ErrPortsvcNotFound, Message: message}
}

func utilPortsvcConflict(message string) error {
	return PortsvcError{Kind: ErrPortsvcConflict, Message: message}
}

func utilNowUTC() time.Time {
	return time.Now().UTC()
}

func utilNormalizeHost(payload HostPayload) (*Host, error) {
	host := &Host{
		Name:   strings.TrimSpace(payload.Name),
		IP:     strings.TrimSpace(payload.IP),
		Spec:   strings.TrimSpace(payload.Spec),
		Status: strings.TrimSpace(payload.Status),
		Notes:  strings.TrimSpace(payload.Notes),
	}
	if host.Name == "" {
		return nil, utilBadPortsvcRequest("host name is required")
	}
	if host.Status == "" {
		host.Status = defaultHostStatus
	}
	if host.Status != defaultHostStatus && host.Status != hostStatusStopped {
		return nil, utilBadPortsvcRequest("host status must be active or stopped")
	}
	return host, nil
}

func utilNormalizeDependencyAsset(payload DependencyAssetPayload) (*DependencyAsset, error) {
	asset := &DependencyAsset{
		Name:            strings.TrimSpace(payload.Name),
		AssetKind:       strings.TrimSpace(payload.AssetKind),
		AssetType:       strings.TrimSpace(payload.AssetType),
		Provider:        strings.TrimSpace(payload.Provider),
		URL:             strings.TrimSpace(payload.URL),
		FullName:        strings.TrimSpace(payload.FullName),
		ExternalID:      strings.TrimSpace(payload.ExternalID),
		Visibility:      strings.TrimSpace(payload.Visibility),
		Controllability: strings.TrimSpace(payload.Controllability),
		Status:          strings.TrimSpace(payload.Status),
		Description:     strings.TrimSpace(payload.Description),
		Metadata:        strings.TrimSpace(payload.Metadata),
		Notes:           strings.TrimSpace(payload.Notes),
	}
	if asset.Name == "" {
		return nil, utilBadPortsvcRequest("asset name is required")
	}
	if asset.AssetKind == "" {
		asset.AssetKind = defaultAssetKind
	}
	if asset.AssetType == "" {
		asset.AssetType = defaultAssetType
	}
	if asset.Provider == "" {
		asset.Provider = defaultAssetProvider
	}
	if asset.Visibility == "" {
		asset.Visibility = defaultAssetVisibility
	}
	if asset.Controllability == "" {
		asset.Controllability = defaultControllability
	}
	if asset.Status == "" {
		asset.Status = defaultAssetStatus
	}
	return asset, nil
}

func utilNormalizePortGroup(ctx context.Context, store *repository.Store, current *PortGroup, payload PortGroupPayload) (*PortGroup, error) {
	hostID := strings.TrimSpace(payload.HostID)
	if hostID != "" {
		if _, err := store.FetchPortsvcHostByID(ctx, hostID); err != nil {
			return nil, err
		}
	}

	portPrefix := payload.PortPrefix
	if portPrefix == 0 {
		allocated, err := utilNextAvailablePortPrefix(ctx, store, utilCurrentPortGroupID(current))
		if err != nil {
			return nil, err
		}
		portPrefix = allocated
	}

	runtimeMode := strings.TrimSpace(payload.RuntimeMode)
	if runtimeMode == "" {
		runtimeMode = defaultRuntimeMode
	}
	group := &PortGroup{
		HostID:           hostID,
		PortPrefix:       portPrefix,
		EnvironmentName:  strings.TrimSpace(payload.EnvironmentName),
		EnvironmentOwner: strings.TrimSpace(payload.EnvironmentOwner),
		RuntimeMode:      runtimeMode,
		RuntimeName:      strings.TrimSpace(payload.RuntimeName),
		ServiceIP:        strings.TrimSpace(payload.ServiceIP),
		Status:           strings.TrimSpace(payload.Status),
		Tags:             strings.TrimSpace(payload.Tags),
		Notes:            strings.TrimSpace(payload.Notes),
	}
	if group.Status == "" {
		group.Status = defaultPortGroupStatus
	}
	if err := utilValidatePortGroup(ctx, store, utilCurrentPortGroupID(current), group); err != nil {
		return nil, err
	}
	return group, nil
}

func utilNormalizeServiceGroup(payload ServiceGroupPayload) (*ServiceGroup, error) {
	group := &ServiceGroup{
		Name:        strings.TrimSpace(payload.Name),
		Kind:        strings.TrimSpace(payload.Kind),
		Status:      strings.TrimSpace(payload.Status),
		Description: strings.TrimSpace(payload.Description),
		Notes:       strings.TrimSpace(payload.Notes),
	}
	if group.Name == "" {
		return nil, utilBadPortsvcRequest("service group name is required")
	}
	if group.Kind == "" {
		group.Kind = defaultServiceGroupKind
	}
	if group.Status == "" {
		group.Status = defaultServiceGroupStatus
	}
	return group, nil
}

func utilCurrentPortGroupID(group *PortGroup) string {
	if group == nil {
		return ""
	}
	return group.ID
}

func utilValidatePortGroup(ctx context.Context, store *repository.Store, currentID string, group *PortGroup) error {
	if group.PortPrefix < allocationPortPrefixMin || group.PortPrefix > allocationPortPrefixMax {
		return utilBadPortsvcRequest("portPrefix must be between 1000 and 5999")
	}
	if group.RuntimeMode != runtimeModeDind && group.RuntimeMode != runtimeModeHost {
		return utilBadPortsvcRequest("runtimeMode must be dind or host")
	}
	if group.RuntimeMode == runtimeModeHost && group.HostID == "" {
		return utilBadPortsvcRequest("hostId is required when runtimeMode is host")
	}
	count, err := store.CountPortsvcOverlappingPortGroups(ctx, currentID, group)
	if err != nil {
		return err
	}
	if count > 0 {
		return utilPortsvcConflict("port group is already allocated")
	}
	return nil
}

func utilNextAvailablePortPrefix(ctx context.Context, store *repository.Store, excludeID string) (int, error) {
	used, err := store.ListPortsvcPortGroupPrefixes(ctx, excludeID)
	if err != nil {
		return 0, err
	}
	usedSet := map[int]struct{}{}
	for _, portPrefix := range used {
		usedSet[portPrefix] = struct{}{}
	}
	for portPrefix := allocationPortPrefixMin; portPrefix <= allocationPortPrefixMax; portPrefix++ {
		if _, exists := usedSet[portPrefix]; !exists {
			return portPrefix, nil
		}
	}
	return 0, utilBadPortsvcRequest("no available port groups")
}

func utilNormalizePortSlot(group PortGroup, payload PortSlotPayload) (*PortSlot, error) {
	slot := &PortSlot{
		ID:            strings.TrimSpace(payload.ID),
		PortGroupID:   group.ID,
		Port:          payload.Port,
		Name:          strings.TrimSpace(payload.Name),
		Kind:          strings.TrimSpace(payload.Kind),
		Protocol:      strings.TrimSpace(payload.Protocol),
		ContainerName: strings.TrimSpace(payload.ContainerName),
		Status:        strings.TrimSpace(payload.Status),
		Notes:         strings.TrimSpace(payload.Notes),
	}
	if slot.Name == "" {
		return nil, utilBadPortsvcRequest("slot name is required")
	}
	portMin, portMax := utilPortGroupBounds(group.PortPrefix)
	if slot.Port < portMin || slot.Port > portMax {
		return nil, utilBadPortsvcRequest("slot port must be inside the port group")
	}
	if slot.Kind == "" {
		slot.Kind = defaultSlotKind
	}
	if slot.Protocol == "" {
		slot.Protocol = defaultSlotProtocol
	}
	if slot.Status == "" {
		slot.Status = defaultSlotStatus
	}
	return slot, nil
}

func utilNormalizePortSlots(group PortGroup, payloads []PortSlotPayload) ([]PortSlot, error) {
	out := make([]PortSlot, 0, len(payloads))
	ports := map[int]struct{}{}
	ids := map[string]struct{}{}
	for _, payload := range payloads {
		if strings.TrimSpace(payload.Name) == "" && payload.Port == 0 {
			continue
		}
		slot, err := utilNormalizePortSlot(group, payload)
		if err != nil {
			return nil, err
		}
		if _, exists := ports[slot.Port]; exists {
			return nil, utilBadPortsvcRequest("slot ports must be unique inside a port group")
		}
		if slot.ID != "" {
			if _, exists := ids[slot.ID]; exists {
				return nil, utilBadPortsvcRequest("slot ids must be unique inside a port group")
			}
			ids[slot.ID] = struct{}{}
		}
		ports[slot.Port] = struct{}{}
		out = append(out, *slot)
	}
	return out, nil
}

func utilNormalizeAssetLinks(ctx context.Context, store *repository.Store, group PortGroup, payloads []PortGroupAssetLinkPayload) ([]PortGroupAssetLink, error) {
	out := make([]PortGroupAssetLink, 0, len(payloads))
	seen := map[string]struct{}{}
	for _, payload := range payloads {
		assetID := strings.TrimSpace(payload.AssetID)
		if assetID == "" {
			continue
		}
		if _, err := store.FetchPortsvcDependencyAssetByID(ctx, assetID); err != nil {
			return nil, err
		}
		portSlotID := strings.TrimSpace(payload.PortSlotID)
		if portSlotID != "" {
			slot, err := store.FetchPortsvcPortSlotByID(ctx, portSlotID)
			if err != nil {
				return nil, err
			}
			if slot.PortGroupID != group.ID {
				return nil, utilBadPortsvcRequest("portSlotId must belong to the port group")
			}
		}
		relationType := strings.TrimSpace(payload.RelationType)
		if relationType == "" {
			relationType = defaultRelationType
		}
		key := portSlotID + "\x00" + assetID + "\x00" + relationType
		if _, exists := seen[key]; exists {
			return nil, utilBadPortsvcRequest("asset links must be unique inside a port group")
		}
		seen[key] = struct{}{}
		out = append(out, PortGroupAssetLink{
			PortGroupID:  group.ID,
			PortSlotID:   portSlotID,
			AssetID:      assetID,
			RelationType: relationType,
			Required:     payload.Required,
			Notes:        strings.TrimSpace(payload.Notes),
		})
	}
	return out, nil
}

func utilNormalizeRepositoryLinks(ctx context.Context, store *repository.Store, group PortGroup, payloads []PortGroupRepositoryLinkPayload) ([]PortGroupRepositoryLink, error) {
	out := make([]PortGroupRepositoryLink, 0, len(payloads))
	seen := map[string]struct{}{}
	for _, payload := range payloads {
		repositoryID := strings.TrimSpace(payload.RepositoryID)
		if repositoryID == "" {
			continue
		}
		if _, err := store.FetchGithubRepositoryByID(ctx, repositoryID); err != nil {
			return nil, err
		}
		portSlotID := strings.TrimSpace(payload.PortSlotID)
		if portSlotID != "" {
			slot, err := store.FetchPortsvcPortSlotByID(ctx, portSlotID)
			if err != nil {
				return nil, err
			}
			if slot.PortGroupID != group.ID {
				return nil, utilBadPortsvcRequest("portSlotId must belong to the port group")
			}
		}
		relationType := strings.TrimSpace(payload.RelationType)
		if relationType == "" {
			relationType = "source"
		}
		key := portSlotID + "\x00" + repositoryID + "\x00" + relationType
		if _, exists := seen[key]; exists {
			return nil, utilBadPortsvcRequest("repository links must be unique inside a port group")
		}
		seen[key] = struct{}{}
		out = append(out, PortGroupRepositoryLink{
			PortGroupID: group.ID, PortSlotID: portSlotID, RepositoryID: repositoryID,
			RelationType: relationType, Required: payload.Required, Notes: strings.TrimSpace(payload.Notes),
		})
	}
	return out, nil
}

func utilNormalizeServiceGroupPortGroups(ctx context.Context, store *repository.Store, group ServiceGroup, payloads []ServiceGroupPortGroupPayload) ([]ServiceGroupPortGroup, error) {
	out := make([]ServiceGroupPortGroup, 0, len(payloads))
	seen := map[string]struct{}{}
	for _, payload := range payloads {
		portGroupID := strings.TrimSpace(payload.PortGroupID)
		if portGroupID == "" {
			continue
		}
		if _, err := store.FetchPortsvcPortGroupByID(ctx, portGroupID); err != nil {
			return nil, err
		}
		if _, exists := seen[portGroupID]; exists {
			return nil, utilBadPortsvcRequest("port groups must be unique inside a service group")
		}
		seen[portGroupID] = struct{}{}
		out = append(out, ServiceGroupPortGroup{
			ServiceGroupID: group.ID,
			PortGroupID:    portGroupID,
			Role:           strings.TrimSpace(payload.Role),
			Notes:          strings.TrimSpace(payload.Notes),
		})
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

func utilPortGroupLabel(group PortGroup) string {
	return strconv.Itoa(group.PortPrefix)
}

func utilPortGroupBounds(portPrefix int) (int, int) {
	portMin := portPrefix * allocationGroupSize
	return portMin, portMin + allocationGroupSize - 1
}

func utilWrapNotFound(err error, message string) error {
	if errors.Is(err, repository.ErrNotFound) {
		return utilPortsvcNotFound(message)
	}
	return err
}

func utilWrapConflict(err error, message string) error {
	if repository.IsUniqueViolation(err) {
		return utilPortsvcConflict(message)
	}
	return err
}

func utilPortSlots(slots []PortSlot) string {
	items := make([]string, 0, len(slots))
	for _, slot := range slots {
		items = append(items, slot.Name+":"+strconv.Itoa(slot.Port))
	}
	return strings.Join(items, "; ")
}

func utilAssetLinks(links []PortGroupAssetLink) string {
	items := make([]string, 0, len(links))
	for _, link := range links {
		label := link.AssetID
		if link.Asset != nil && link.Asset.Name != "" {
			label = link.Asset.Name
		}
		items = append(items, link.RelationType+":"+label)
	}
	return strings.Join(items, "; ")
}

func utilRepositoryLinks(links []PortGroupRepositoryLink) string {
	items := make([]string, 0, len(links))
	for _, link := range links {
		label := link.RepositoryID
		if link.Repository != nil && link.Repository.FullName != "" {
			label = link.Repository.FullName
		}
		items = append(items, link.RelationType+":"+label)
	}
	return strings.Join(items, "; ")
}
