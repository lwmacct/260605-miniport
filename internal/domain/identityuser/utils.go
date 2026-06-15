package identityuser

import "strings"

func defaultDisplayName(username string, displayName string) string {
	displayName = strings.TrimSpace(displayName)
	if displayName != "" {
		return displayName
	}
	return username
}
