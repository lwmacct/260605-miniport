package service

import (
	"encoding/csv"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

const (
	defaultGroupStatus    = "planned"
	defaultSlotProtocol   = "tcp"
	defaultSlotStatus     = "empty"
	defaultComponentType  = "opensource"
	defaultRepositoryKind = "source"
)

type InventoryHost = repository.InventoryHostRecord
type InventoryPortGroup = repository.InventoryPortGroupRecord
type InventoryPortSlot = repository.InventoryPortSlotRecord
type InventoryComponent = repository.InventoryComponentRecord
type InventoryRepositoryRef = repository.InventoryRepositoryRefRecord

type HostPayload struct {
	IP          string
	Name        string
	Network     string
	Environment string
	Notes       string
}

type PortSlotPayload struct {
	Port     int
	Name     string
	Protocol string
	Purpose  string
	Status   string
	Notes    string
}

type ComponentPayload struct {
	Name    string
	Type    string
	URL     string
	Version string
	Notes   string
}

type RepositoryPayload struct {
	Name  string
	URL   string
	Kind  string
	Notes string
}

type PortGroupPayload struct {
	HostID        int64
	PortStart     int
	PortEnd       int
	ServiceName   string
	ContainerName string
	DindHost      string
	Status        string
	Owner         string
	Tags          string
	Notes         string
	Slots         []PortSlotPayload
	Components    []ComponentPayload
	Repositories  []RepositoryPayload
}

type HostListParams struct {
	Environment string
	Query       string
	Sort        string
}

type PortGroupView struct {
	InventoryPortGroup

	Slots        []InventoryPortSlot
	Components   []InventoryComponent
	Repositories []InventoryRepositoryRef
}

type PortGroupListParams struct {
	HostID int64
	Query  string
	Sort   string
	Status string
}

type PortGroupBatchUpdateInput struct {
	IDs    []int64
	Owner  *string
	Status *string
	Tags   *string
}

type PortGroupBatchDeleteInput struct {
	IDs []int64
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

func utilInventoryHostFromPayload(payload HostPayload) *InventoryHost {
	return &InventoryHost{
		IP:          strings.TrimSpace(payload.IP),
		Name:        strings.TrimSpace(payload.Name),
		Network:     strings.TrimSpace(payload.Network),
		Environment: strings.TrimSpace(payload.Environment),
		Notes:       strings.TrimSpace(payload.Notes),
	}
}

func validateHost(host *InventoryHost) error {
	if host.IP == "" {
		return utilBadInventoryRequest("host ip is required")
	}
	return nil
}

func utilUniqueInt64s(items []int64) []int64 {
	seen := make(map[int64]struct{}, len(items))
	out := make([]int64, 0, len(items))
	for _, item := range items {
		if item <= 0 {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func utilCSVBytes(records [][]string) ([]byte, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return []byte(builder.String()), nil
}

func utilPortGroupComponents(items []InventoryComponent) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		label := item.Name
		if version := strings.TrimSpace(item.Version); version != "" {
			label += "@" + version
		}
		parts = append(parts, label)
	}
	return strings.Join(parts, "; ")
}

func utilPortGroupRepositories(items []InventoryRepositoryRef) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, item.Kind+":"+item.URL)
	}
	return strings.Join(parts, "; ")
}

func utilPortGroupSlots(items []InventoryPortSlot) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		label := strconv.Itoa(item.Port) + "/" + item.Protocol + "/" + item.Status
		if name := strings.TrimSpace(item.Name); name != "" {
			label += "/" + name
		}
		parts = append(parts, label)
	}
	return strings.Join(parts, "; ")
}

func utilNowUTC() time.Time { return time.Now().UTC() }
