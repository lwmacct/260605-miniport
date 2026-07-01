package service

import (
	"context"
	"encoding/csv"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/lwmacct/260605-miniport/internal/repository"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/identity"
)

const (
	allocationPortMin      = 10000
	allocationPortMax      = 59999
	allocationPortMaxStart = 59990
	allocationGroupSize    = 10
	defaultHostStatus      = "active"
	defaultPortGroupStatus = "available"
	defaultRuntimeMode     = "dind"
	defaultSlotKind        = "app"
	defaultSlotProtocol    = "tcp"
	defaultSlotStatus      = "planned"
	defaultAssetKind       = "component"
	defaultAssetType       = "middleware"
	defaultAssetProvider   = "manual"
	defaultAssetVisibility = "unknown"
	defaultControllability = "unknown"
	defaultAssetStatus     = "active"
	defaultRelationType    = "runtime"
	runtimeModeDind        = "dind"
	runtimeModeHost        = "host"
)

type Host = repository.PortsvcHostRecord
type PortGroup = repository.PortsvcPortGroupRecord
type PortSlot = repository.PortsvcPortSlotRecord
type DependencyAsset = repository.PortsvcDependencyAssetRecord
type PortGroupAssetLink = repository.PortsvcPortGroupAssetLinkRecord

type PortsvcActor struct {
	OwnerSubject string
	OwnerName    string
	Admin        bool
}

type HostPayload struct {
	Name   string
	IP     string
	Spec   string
	Status string
	Notes  string
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
	OwnerSubject    string
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

type PortGroupAssetLinkPayload struct {
	ID           string
	PortSlotID   string
	AssetID      string
	RelationType string
	Required     bool
	Notes        string
}

type PortGroupPayload struct {
	OwnerSubject string
	HostID       string
	PortStart    int
	PortEnd      int
	ProjectName  string
	ProjectOwner string
	RuntimeMode  string
	RuntimeName  string
	ServiceIP    string
	Status       string
	Tags         string
	Notes        string
	Slots        []PortSlotPayload
	AssetLinks   []PortGroupAssetLinkPayload
}

type PortGroupView struct {
	PortGroup

	Slots      []PortSlot
	AssetLinks []PortGroupAssetLink
}

type HostListParams struct {
	Query  string
	Status string
}

type PortGroupListParams struct {
	Actor        PortsvcActor
	OwnerSubject string
	Query        string
	Sort         string
	Status       string
}

type DependencyAssetListParams struct {
	Actor        PortsvcActor
	OwnerSubject string
	Query        string
	AssetKind    string
	AssetType    string
	Provider     string
	Status       string
}

var (
	ErrPortsvcBadRequest = errors.New("portsvc bad request")
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

func utilNowUTC() time.Time {
	return time.Now().UTC()
}

func utilVisibleOwnerSubject(actor PortsvcActor, requestedOwnerSubject string) string {
	if actor.Admin {
		return strings.TrimSpace(requestedOwnerSubject)
	}
	return actor.OwnerSubject
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
	return host, nil
}

func utilNormalizeDependencyAsset(actor PortsvcActor, current *DependencyAsset, payload DependencyAssetPayload) (*DependencyAsset, error) {
	ownerSubject := strings.TrimSpace(actor.OwnerSubject)
	if actor.Admin {
		if subject := strings.TrimSpace(payload.OwnerSubject); subject != "" {
			ownerSubject = subject
		} else if current != nil {
			ownerSubject = current.OwnerSubject
		}
	}
	if ownerSubject == "" {
		return nil, utilBadPortsvcRequest("ownerSubject is required")
	}
	asset := &DependencyAsset{
		OwnerSubject:    ownerSubject,
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

func utilNormalizePortGroup(ctx context.Context, store *repository.Store, actor PortsvcActor, current *PortGroup, payload PortGroupPayload) (*PortGroup, error) {
	ownerSubject := strings.TrimSpace(actor.OwnerSubject)
	if actor.Admin {
		if subject := strings.TrimSpace(payload.OwnerSubject); subject != "" {
			ownerSubject = subject
		} else if current != nil {
			ownerSubject = current.OwnerSubject
		}
	}
	if ownerSubject == "" {
		return nil, utilBadPortsvcRequest("ownerSubject is required")
	}

	hostID := strings.TrimSpace(payload.HostID)
	if hostID != "" {
		if _, err := store.FetchPortsvcHostByID(ctx, hostID); err != nil {
			return nil, err
		}
	}

	portStart := payload.PortStart
	if portStart == 0 {
		allocated, err := utilNextAvailablePortStart(ctx, store, ownerSubject, utilCurrentPortGroupID(current))
		if err != nil {
			return nil, err
		}
		portStart = allocated
	}
	portEnd := portStart + allocationGroupSize - 1
	if payload.PortEnd > 0 && payload.PortEnd != portEnd {
		return nil, utilBadPortsvcRequest("portEnd must equal portStart + 9")
	}

	runtimeMode := strings.TrimSpace(payload.RuntimeMode)
	if runtimeMode == "" {
		runtimeMode = defaultRuntimeMode
	}
	group := &PortGroup{
		OwnerSubject: ownerSubject,
		HostID:       hostID,
		PortStart:    portStart,
		PortEnd:      portEnd,
		ProjectName:  strings.TrimSpace(payload.ProjectName),
		ProjectOwner: strings.TrimSpace(payload.ProjectOwner),
		RuntimeMode:  runtimeMode,
		RuntimeName:  strings.TrimSpace(payload.RuntimeName),
		ServiceIP:    strings.TrimSpace(payload.ServiceIP),
		Status:       strings.TrimSpace(payload.Status),
		Tags:         strings.TrimSpace(payload.Tags),
		Notes:        strings.TrimSpace(payload.Notes),
	}
	if group.Status == "" {
		group.Status = defaultPortGroupStatus
	}
	if err := utilValidatePortGroup(ctx, store, utilCurrentPortGroupID(current), group); err != nil {
		return nil, err
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
	if group.PortStart < allocationPortMin || group.PortStart > allocationPortMaxStart {
		return utilBadPortsvcRequest("portStart must be between 10000 and 59990")
	}
	if group.PortStart%allocationGroupSize != 0 {
		return utilBadPortsvcRequest("portStart must align to a 10-port group")
	}
	if group.PortEnd != group.PortStart+allocationGroupSize-1 {
		return utilBadPortsvcRequest("port group must contain exactly 10 ports")
	}
	if group.PortEnd > allocationPortMax {
		return utilBadPortsvcRequest("portEnd must be less than or equal to 59999")
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
		return utilBadPortsvcRequest("port group is already allocated for this user")
	}
	return nil
}

func utilNextAvailablePortStart(ctx context.Context, store *repository.Store, ownerSubject string, excludeID string) (int, error) {
	used, err := store.ListPortsvcPortGroupStartsByOwner(ctx, ownerSubject, excludeID)
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
	return 0, utilBadPortsvcRequest("no available port groups for this user")
}

func utilNormalizePortSlot(group PortGroup, payload PortSlotPayload) (*PortSlot, error) {
	slot := &PortSlot{
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
	if slot.Port < group.PortStart || slot.Port > group.PortEnd {
		return nil, utilBadPortsvcRequest("slot port must be inside the port group range")
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
		asset, err := store.FetchPortsvcDependencyAssetByID(ctx, assetID)
		if err != nil {
			return nil, err
		}
		if asset.OwnerSubject != group.OwnerSubject {
			return nil, utilBadPortsvcRequest("asset must belong to the port group owner")
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

func utilResolveActivePrincipal(ctx context.Context, directory identity.Directory, ownerSubject string) error {
	if directory == nil {
		return utilBadPortsvcRequest("identity directory is required")
	}
	principal, err := directory.Principal(ctx, ownerSubject)
	if err != nil {
		return err
	}
	if !principal.Active() {
		return utilBadPortsvcRequest("ownerSubject is invalid")
	}
	return nil
}

func utilPrincipalName(principal *identity.Principal, fallback string) string {
	if principal == nil {
		return fallback
	}
	if principal.DisplayName != "" {
		return principal.DisplayName
	}
	if principal.Username != "" {
		return principal.Username
	}
	if principal.Subject != "" {
		return principal.Subject
	}
	return fallback
}

func utilCSVBytes(records [][]string) ([]byte, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}
	return []byte(builder.String()), nil
}

func utilPortRange(group PortGroup) string {
	return strconv.Itoa(group.PortStart) + "-" + strconv.Itoa(group.PortEnd)
}

func utilWrapNotFound(err error, message string) error {
	if errors.Is(err, repository.ErrNotFound) {
		return utilPortsvcNotFound(message)
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
