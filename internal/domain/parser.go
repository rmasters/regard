package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"regard/internal/query"
)

// CreateSummary converts a QueryResult into a structured domain summary
func CreateSummary(result query.QueryResult) Summary {
	summary := Summary{
		Domain:    result.Query,
		Protocol:  result.Protocol,
		QueryType: result.Type,
	}

	if result.Protocol == "RDAP" {
		summary = parseRDAPSummary(result, summary)
	} else {
		summary = parseWhoisSummary(result, summary)
	}

	// Parse ASN information if this is an ASN query
	if result.Type == string(query.QueryTypeASN) {
		summary.ASN = parseASNInfo(result)
		// For ASNs, use the ASN status as the summary status
		if summary.ASN != nil && summary.ASN.Status != "" {
			summary.Status = summary.ASN.Status
		}
	}

	// Add post-expiration guidance if needed
	if summary.Timeline.Expiration != nil {
		summary.PostExpiration = GeneratePostExpirationGuidance(summary)
	}

	return summary
}

func parseRDAPSummary(result query.QueryResult, summary Summary) Summary {
	// Parse RDAP response - need to handle the fact that result.Data might be a struct
	var domainData map[string]interface{}

	// Convert the RDAP struct to a map for easier parsing
	if jsonBytes, err := json.Marshal(result.Data); err == nil {
		_ = json.Unmarshal(jsonBytes, &domainData)
	}

	if domainData != nil {
		// Domain status
		if statusArray, ok := domainData["Status"].([]interface{}); ok {
			for _, status := range statusArray {
				if statusStr, ok := status.(string); ok {
					summary.StatusDetails = append(summary.StatusDetails, statusStr)
				}
			}
			// For RDAP, we can check the raw data for better status detection
			rawData := ""
			if jsonBytes, err := json.Marshal(result.Data); err == nil {
				rawData = string(jsonBytes)
			}
			summary.Status = InterpretStatus(summary.StatusDetails, rawData)
		}

		// Nameservers
		if nsArray, ok := domainData["Nameservers"].([]interface{}); ok {
			for _, ns := range nsArray {
				if nsObj, ok := ns.(map[string]interface{}); ok {
					if name, ok := nsObj["LDHName"].(string); ok {
						summary.Nameservers = append(summary.Nameservers, name)
					}
				}
			}
		}

		// DNSSEC
		if secureDNS, ok := domainData["SecureDNS"].(map[string]interface{}); ok {
			if delegationSigned, ok := secureDNS["DelegationSigned"].(bool); ok {
				summary.DNSSEC.Enabled = delegationSigned
				if delegationSigned {
					summary.DNSSEC.Details = "Delegation signed"
				}
			}
		}

		// Events (timeline)
		if eventsArray, ok := domainData["Events"].([]interface{}); ok {
			for _, event := range eventsArray {
				if eventObj, ok := event.(map[string]interface{}); ok {
					action, _ := eventObj["Action"].(string)
					dateStr, _ := eventObj["Date"].(string)

					if date, err := time.Parse(time.RFC3339, dateStr); err == nil {
						timelineEvent := &TimelineEvent{
							Date:          date,
							HumanReadable: HumanReadableTime(date),
						}

						switch action {
						case "registration":
							summary.Timeline.Registration = timelineEvent
						case "last changed", "last update of RDAP database":
							summary.Timeline.LastUpdated = timelineEvent
						case "expiration":
							summary.Timeline.Expiration = timelineEvent
						}
					}
				}
			}
		}

		// Registrar info
		if entitiesArray, ok := domainData["Entities"].([]interface{}); ok {
			for _, entity := range entitiesArray {
				if entityObj, ok := entity.(map[string]interface{}); ok {
					if rolesArray, ok := entityObj["Roles"].([]interface{}); ok {
						for _, role := range rolesArray {
							if roleStr, ok := role.(string); ok && roleStr == "registrar" {
								if vcard, ok := entityObj["VCard"].(map[string]interface{}); ok {
									if props, ok := vcard["Properties"].([]interface{}); ok {
										for _, prop := range props {
											if propObj, ok := prop.(map[string]interface{}); ok {
												if name, _ := propObj["Name"].(string); name == "fn" {
													if value, ok := propObj["Value"].(string); ok {
														summary.Registrar.Name = value
													}
												}
											}
										}
									}
								}
								if handle, ok := entityObj["Handle"].(string); ok {
									summary.Registrar.ID = handle
								}
								break
							}
						}
					}
				}
			}
		}
	}

	return summary
}

func parseWhoisSummary(result query.QueryResult, summary Summary) Summary {
	// Parse WHOIS response
	if data, ok := result.Data.(map[string]interface{}); ok {
		if fields, ok := data["parsed_fields"].(map[string]interface{}); ok {
			// Domain status from various possible fields
			statusFields := []string{"domain_status", "status"}
			for _, field := range statusFields {
				if status, ok := fields[field].(string); ok {
					// Extract status codes from WHOIS format
					statusParts := strings.Fields(status)
					for _, part := range statusParts {
						if strings.Contains(part, "rohibited") || strings.Contains(part, "ctive") || strings.Contains(part, "ransfer") || strings.Contains(part, "client") {
							// Clean up status - remove URLs
							cleanStatus := strings.Split(part, " ")[0]
							if !strings.HasPrefix(cleanStatus, "http") {
								summary.StatusDetails = append(summary.StatusDetails, cleanStatus)
							}
						}
					}
				}
			}
			// Get raw data for better status detection
			rawData := ""
			if data, ok := result.Data.(map[string]interface{}); ok {
				if rawResponse, ok := data["raw_response"].(string); ok {
					rawData = rawResponse
				}
			}
			summary.Status = InterpretStatus(summary.StatusDetails, rawData)

			// Nameservers - collect all nameserver entries
			for key, value := range fields {
				if strings.Contains(strings.ToLower(key), "name_server") || strings.Contains(strings.ToLower(key), "nameserver") {
					if ns, ok := value.(string); ok {
						// Avoid duplicates
						found := false
						for _, existing := range summary.Nameservers {
							if existing == ns {
								found = true
								break
							}
						}
						if !found {
							summary.Nameservers = append(summary.Nameservers, ns)
						}
					}
				}
			}

			// DNSSEC
			if dnssec, ok := fields["dnssec"].(string); ok {
				summary.DNSSEC.Enabled = dnssec == "signedDelegation"
				summary.DNSSEC.Details = dnssec
			}

			// Timeline
			if created, ok := fields["creation_date"].(string); ok {
				if date, err := parseWhoisDate(created); err == nil {
					summary.Timeline.Registration = &TimelineEvent{
						Date:          date,
						HumanReadable: HumanReadableTime(date),
					}
				}
			}

			if updated, ok := fields["updated_date"].(string); ok {
				if date, err := parseWhoisDate(updated); err == nil {
					summary.Timeline.LastUpdated = &TimelineEvent{
						Date:          date,
						HumanReadable: HumanReadableTime(date),
					}
				}
			}

			if expiry, ok := fields["registry_expiry_date"].(string); ok {
				if date, err := parseWhoisDate(expiry); err == nil {
					summary.Timeline.Expiration = &TimelineEvent{
						Date:          date,
						HumanReadable: HumanReadableTime(date),
					}
				}
			}

			// Registrar
			if registrar, ok := fields["registrar"].(string); ok {
				summary.Registrar.Name = registrar
			}
			if registrarID, ok := fields["registrar_iana_id"].(string); ok {
				summary.Registrar.ID = registrarID
			}
		}
	}

	return summary
}

func parseWhoisDate(dateStr string) (time.Time, error) {
	// Try different date formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02",
		"2006-01-02 15:04:05",
		"02-Jan-2006",
		"2006/01/02",
	}

	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// HumanReadableTime converts a time to a human-readable relative format
func HumanReadableTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if t.After(now) {
		// Future date
		diff = t.Sub(now)
		days := int(diff.Hours() / 24)
		if days == 0 {
			return "today"
		} else if days == 1 {
			return "tomorrow"
		} else if days < 30 {
			return fmt.Sprintf("in %d days", days)
		} else if days < 365 {
			months := days / 30
			return fmt.Sprintf("in %d months", months)
		} else {
			years := days / 365
			return fmt.Sprintf("in %d years", years)
		}
	} else {
		// Past date
		days := int(diff.Hours() / 24)
		if days == 0 {
			return "today"
		} else if days == 1 {
			return "yesterday"
		} else if days < 30 {
			return fmt.Sprintf("%d days ago", days)
		} else if days < 365 {
			months := days / 30
			return fmt.Sprintf("%d months ago", months)
		} else {
			years := days / 365
			return fmt.Sprintf("%d years ago", years)
		}
	}
}

// parseASNInfo extracts ASN information from WHOIS data
func parseASNInfo(result query.QueryResult) *ASNInfo {
	if result.RawData == "" {
		return nil
	}

	asn := &ASNInfo{
		Number: result.Query,
	}

	lines := strings.Split(result.RawData, "\n")
	var currentSection string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle abuse contact extraction from comment lines
		if strings.HasPrefix(line, "%") {
			if strings.Contains(strings.ToLower(line), "abuse contact for") && strings.Contains(line, "@") {
				// Extract email from "% Abuse contact for 'AS9009' is 'abuse@m247.ro'"
				start := strings.Index(line, "'")
				if start != -1 {
					end := strings.LastIndex(line, "'")
					if end > start {
						email := line[start+1 : end]
						if strings.Contains(email, "@") {
							asn.AbuseContact = email
						}
					}
				}
			}
			continue
		}

		// Parse key-value pairs
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				if value == "" {
					continue
				}

				switch strings.ToLower(key) {
				case "aut-num", "asnumber":
					asn.Number = value
				case "as-name", "asname":
					asn.Name = value
				case "descr", "description":
					if asn.Description == "" {
						asn.Description = value
					}
				case "country":
					asn.Country = value
				case "org-name", "orgname":
					asn.Organization = value
				case "status":
					asn.Status = value
				case "abuse-mailbox", "orgabuseemail", "orgabuse-email":
					asn.AbuseContact = value
				case "organisation", "organization": //nolint:misspell // Both spellings are valid (RIPE uses 's', ARIN uses 'z')
					currentSection = "org"
				case "role":
					currentSection = "role"
				}

				// Handle organization section
				if currentSection == "org" && strings.ToLower(key) == "org-name" {
					asn.Organization = value
				}

				// Parse import/export statements to extract peers
				if strings.HasPrefix(strings.ToLower(key), "import") || strings.HasPrefix(strings.ToLower(key), "export") {
					if strings.Contains(value, "AS") {
						// Extract AS numbers from import/export lines
						words := strings.Fields(value)
						for _, word := range words {
							if strings.HasPrefix(word, "AS") && len(word) > 2 {
								// Clean up AS number (remove trailing punctuation)
								asNum := strings.TrimRight(word, ",;:")
								if len(asNum) > 2 {
									// Avoid duplicates
									found := false
									for _, existing := range asn.Peers {
										if existing == asNum {
											found = true
											break
										}
									}
									if !found && asNum != asn.Number {
										asn.Peers = append(asn.Peers, asNum)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// If no organization name was found, use the first description
	if asn.Organization == "" && asn.Description != "" {
		asn.Organization = asn.Description
	}

	// Set status to active if we have substantial data but no explicit status
	if asn.Status == "" && (asn.Name != "" || asn.Organization != "") {
		asn.Status = "active"
	}

	return asn
}
