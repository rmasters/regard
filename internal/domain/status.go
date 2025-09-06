package domain

import "strings"

// InterpretStatus determines the overall domain status from status details and raw data
func InterpretStatus(statusDetails []string, rawData string) string {
	// Check for explicit "not found" indicators first
	lowerRaw := strings.ToLower(rawData)
	if strings.Contains(lowerRaw, "not found") ||
		strings.Contains(lowerRaw, "no match") ||
		strings.Contains(lowerRaw, "no data found") ||
		strings.Contains(lowerRaw, "domain not found") ||
		strings.Contains(lowerRaw, "available for registration") {
		return "available"
	}

	if len(statusDetails) == 0 {
		// If no status details but we have raw data, it might be available
		if rawData != "" && (strings.Contains(lowerRaw, "no entries found") || len(strings.TrimSpace(rawData)) < 100) {
			return "available"
		}
		return "unknown"
	}

	hasActive := false
	hasProhibited := false

	for _, status := range statusDetails {
		lowerStatus := strings.ToLower(status)
		if strings.Contains(lowerStatus, "active") || strings.Contains(lowerStatus, "ok") {
			hasActive = true
		}
		if strings.Contains(lowerStatus, "prohibited") {
			hasProhibited = true
		}
	}

	if hasActive || hasProhibited {
		return "active"
	}

	return "unknown"
}