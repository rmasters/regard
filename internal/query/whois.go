package query

import (
	"strings"
	"time"

	"github.com/likexian/whois"
)

// PerformWhoisQuery executes a WHOIS query for the given input
func PerformWhoisQuery(query string) QueryResult {
	result := QueryResult{
		Query:     query,
		Type:      string(DetectQueryType(query)),
		Protocol:  "WHOIS",
		Timestamp: time.Now(),
	}

	response, err := whois.Whois(query)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result
	}

	result.Success = true
	result.RawData = response

	// Parse WHOIS data into a structured format for JSON output
	result.Data = parseWhoisData(response)

	return result
}

func parseWhoisData(whoisData string) map[string]interface{} {
	data := make(map[string]interface{})
	lines := strings.Split(whoisData, "\n")

	sectionData := make(map[string]interface{})

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "%") || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				if value != "" {
					// Convert key to lowercase for consistency
					key = strings.ToLower(strings.ReplaceAll(key, " ", "_"))
					sectionData[key] = value
				}
			}
		}
	}

	if len(sectionData) > 0 {
		data["parsed_fields"] = sectionData
	}

	data["raw_response"] = whoisData
	return data
}