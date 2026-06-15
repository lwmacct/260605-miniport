package identityuser

import "strings"

func normalizeUsername(username string) string {
	username = strings.TrimSpace(strings.ToLower(username))
	if username == "" {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(username))
	lastDash := false
	for _, r := range username {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			lastDash = false
			continue
		}
		if r == '.' || r == '_' || r == '-' || r == ' ' {
			if builder.Len() == 0 || lastDash {
				continue
			}
			builder.WriteByte('-')
			lastDash = true
		}
	}

	value := strings.Trim(builder.String(), "-")
	if len(value) > 64 {
		value = strings.Trim(value[:64], "-")
	}
	return value
}
