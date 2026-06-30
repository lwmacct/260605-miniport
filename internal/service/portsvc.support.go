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
	allocationPortMin       = 10000
	allocationPortMax       = 59999
	allocationPortMaxStart  = 59990
	allocationGroupSize     = 10
	defaultAllocationStatus = "available"
	defaultServiceStatus    = "planned"
	defaultDependencyType   = "opensource"
	defaultRepositoryKind   = "source"
)

type PortsvcServiceRecord = repository.PortsvcServiceRecord
type PortAllocation = repository.PortsvcPortAllocationRecord
type Dependency = repository.PortsvcDependencyRecord
type RepositoryRef = repository.PortsvcRepositoryRecord

type PortsvcActor struct {
	OwnerSubject string
	OwnerName    string
	Admin        bool
}

type RepositoryPayload struct {
	ID    string
	Name  string
	URL   string
	Kind  string
	Role  string
	Notes string
}

type DependencyPayload struct {
	ID      string
	Name    string
	Type    string
	URL     string
	Version string
	Role    string
	Notes   string
}

type ServicePayload struct {
	OwnerSubject     string
	PortAllocationID string
	Name             string
	ProjectName      string
	DindIP           string
	DindContainer    string
	Status           string
	Owner            string
	Tags             string
	Notes            string
	Repositories     []RepositoryPayload
	Dependencies     []DependencyPayload
}

type PortAllocationPayload struct {
	OwnerSubject string
	PortStart    int
	PortEnd      int
	Status       string
	Notes        string
}

type ServiceView struct {
	PortsvcServiceRecord

	Repositories []RepositoryRef
	Dependencies []Dependency
}

type ServiceListParams struct {
	Actor        PortsvcActor
	OwnerSubject string
	Query        string
	Sort         string
	Status       string
	ProjectName  string
}

type PortAllocationListParams struct {
	Actor        PortsvcActor
	OwnerSubject string
	Sort         string
	Status       string
}

type ServiceBatchDeleteInput struct {
	Actor PortsvcActor
	IDs   []string
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

func utilNormalizeService(ctx context.Context, store *repository.Store, actor PortsvcActor, current *PortsvcServiceRecord, payload ServicePayload) (*PortsvcServiceRecord, error) {
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
	portID := payload.PortAllocationID
	if portID != "" {
		port, err := store.FetchPortsvcPortAllocationByID(ctx, portID)
		if err != nil {
			return nil, err
		}
		if port.OwnerSubject != ownerSubject {
			return nil, utilBadPortsvcRequest("port allocation must belong to the service owner")
		}
	}
	service := &PortsvcServiceRecord{
		OwnerSubject:     ownerSubject,
		PortAllocationID: portID,
		Name:             strings.TrimSpace(payload.Name),
		ProjectName:      strings.TrimSpace(payload.ProjectName),
		DindIP:           strings.TrimSpace(payload.DindIP),
		DindContainer:    strings.TrimSpace(payload.DindContainer),
		Status:           strings.TrimSpace(payload.Status),
		Owner:            strings.TrimSpace(payload.Owner),
		Tags:             strings.TrimSpace(payload.Tags),
		Notes:            strings.TrimSpace(payload.Notes),
	}
	if service.Name == "" {
		return nil, utilBadPortsvcRequest("name is required")
	}
	if service.Status == "" {
		service.Status = defaultServiceStatus
	}
	return service, nil
}

func utilNormalizePortAllocation(ctx context.Context, store *repository.Store, actor PortsvcActor, current *PortAllocation, payload PortAllocationPayload) (*PortAllocation, error) {
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

	portStart := payload.PortStart
	if portStart == 0 {
		allocated, err := utilNextAvailablePortStart(ctx, store, ownerSubject, utilCurrentPortID(current))
		if err != nil {
			return nil, err
		}
		portStart = allocated
	}
	portEnd := portStart + allocationGroupSize - 1
	if payload.PortEnd > 0 && payload.PortEnd != portEnd {
		return nil, utilBadPortsvcRequest("portEnd must equal portStart + 9")
	}
	group := &PortAllocation{
		OwnerSubject: ownerSubject,
		PortStart:    portStart,
		PortEnd:      portEnd,
		Status:       strings.TrimSpace(payload.Status),
		Notes:        strings.TrimSpace(payload.Notes),
	}
	if group.Status == "" {
		group.Status = defaultAllocationStatus
	}
	if err := utilValidatePortAllocation(ctx, store, utilCurrentPortID(current), group); err != nil {
		return nil, err
	}
	return group, nil
}

func utilCurrentPortID(group *PortAllocation) string {
	if group == nil {
		return ""
	}
	return group.ID
}

func utilValidatePortAllocation(ctx context.Context, store *repository.Store, currentID string, group *PortAllocation) error {
	if group.PortStart < allocationPortMin || group.PortStart > allocationPortMaxStart {
		return utilBadPortsvcRequest("portStart must be between 10000 and 59990")
	}
	if group.PortStart%allocationGroupSize != 0 {
		return utilBadPortsvcRequest("portStart must align to a 10-port group")
	}
	if group.PortEnd != group.PortStart+allocationGroupSize-1 {
		return utilBadPortsvcRequest("port allocation must contain exactly 10 ports")
	}
	if group.PortEnd > allocationPortMax {
		return utilBadPortsvcRequest("portEnd must be less than or equal to 59999")
	}
	count, err := store.CountPortsvcOverlappingPortAllocations(ctx, currentID, group)
	if err != nil {
		return err
	}
	if count > 0 {
		return utilBadPortsvcRequest("port allocation is already allocated for this user")
	}
	return nil
}

func utilNextAvailablePortStart(ctx context.Context, store *repository.Store, ownerSubject string, excludeID string) (int, error) {
	used, err := store.ListPortsvcPortAllocationStartsByOwner(ctx, ownerSubject, excludeID)
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
	return 0, utilBadPortsvcRequest("no available port allocations for this user")
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

func utilNormalizeIDs(ids []string) ([]string, error) {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			return nil, utilBadPortsvcRequest("ids must be non-empty")
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	if len(out) == 0 {
		return nil, utilBadPortsvcRequest("ids are required")
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

func utilServiceRepositories(items []RepositoryRef) string {
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

func utilServiceDependencies(items []Dependency) string {
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

func utilPortRange(group *PortAllocation) string {
	if group == nil {
		return ""
	}
	return strconv.Itoa(group.PortStart) + "-" + strconv.Itoa(group.PortEnd)
}
