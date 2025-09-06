package query

import "strings"

// DetectQueryType determines the type of query based on the input string
func DetectQueryType(query string) QueryType {
	// Simple heuristics to detect query type
	if strings.Contains(query, ".") {
		// Could be domain or IP
		if strings.Count(query, ".") == 3 {
			// Likely IPv4
			parts := strings.Split(query, ".")
			for _, part := range parts {
				if len(part) > 3 {
					return QueryTypeDomain
				}
			}
			return QueryTypeIP
		}
		return QueryTypeDomain
	}

	// Check for ASN pattern (AS followed by numbers)
	if strings.HasPrefix(strings.ToUpper(query), "AS") && len(query) > 2 {
		return QueryTypeASN
	}

	// Check for IPv6
	if strings.Contains(query, ":") {
		return QueryTypeIP
	}

	return QueryTypeDomain
}
