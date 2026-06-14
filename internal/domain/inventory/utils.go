package inventory

import (
	"encoding/csv"
	"strconv"
	"strings"
)

func compactString(value string) string {
	return strings.TrimSpace(value)
}

func searchPattern(value string) string {
	return "%" + strings.ToLower(strings.TrimSpace(value)) + "%"
}

func joinSearchClauses(columns []string) string {
	clauses := make([]string, 0, len(columns))
	for _, column := range columns {
		clauses = append(clauses, "LOWER("+column+") LIKE ?")
	}
	return "(" + strings.Join(clauses, " OR ") + ")"
}

func joinSearchArgs(pattern string, count int) []any {
	args := make([]any, 0, count)
	for idx := 0; idx < count; idx++ {
		args = append(args, pattern)
	}
	return args
}

func csvBytes(records [][]string) ([]byte, error) {
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

func formatPortGroupComponents(items []Component) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		label := item.Name
		if version := compactString(item.Version); version != "" {
			label += "@" + version
		}
		parts = append(parts, label)
	}
	return strings.Join(parts, "; ")
}

func formatPortGroupRepositories(items []Repository) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, item.Kind+":"+item.URL)
	}
	return strings.Join(parts, "; ")
}

func formatPortGroupSlots(items []PortSlot) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		label := strconv.Itoa(item.Port) + "/" + item.Protocol + "/" + item.Status
		if name := compactString(item.Name); name != "" {
			label += "/" + name
		}
		parts = append(parts, label)
	}
	return strings.Join(parts, "; ")
}
