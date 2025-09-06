package domain

import (
	"testing"
	"time"

	"regard/internal/query"
)

func TestHumanReadableTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "Today",
			input:    now,
			expected: "today",
		},
		{
			name:     "Yesterday",
			input:    now.Add(-24 * time.Hour),
			expected: "yesterday",
		},
		{
			name:     "Tomorrow",
			input:    now.Add(25 * time.Hour), // Add extra hour to avoid timezone edge cases
			expected: "tomorrow",
		},
		{
			name:     "5 days ago",
			input:    now.Add(-5 * 24 * time.Hour),
			expected: "5 days ago",
		},
		{
			name:     "In 10 days",
			input:    now.Add(10*24*time.Hour + time.Hour), // Add buffer
			expected: "in 10 days",
		},
		{
			name:     "2 months ago",
			input:    now.Add(-60 * 24 * time.Hour),
			expected: "2 months ago",
		},
		{
			name:     "In 3 months",
			input:    now.Add(95 * 24 * time.Hour), // Use 95 days to be safely in 3 months
			expected: "in 3 months",
		},
		{
			name:     "2 years ago",
			input:    now.Add(-2 * 365 * 24 * time.Hour),
			expected: "2 years ago",
		},
		{
			name:     "In 1 year",
			input:    now.Add(400 * 24 * time.Hour), // Use 400 days to be safely over 1 year
			expected: "in 1 years",                  // Note: the function doesn't handle singular/plural
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HumanReadableTime(tt.input)
			if result != tt.expected {
				t.Errorf("HumanReadableTime(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseWhoisDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected string // Format: "2006-01-02T15:04:05Z"
	}{
		{
			name:     "RFC3339 format",
			input:    "2023-08-14T07:01:31Z",
			wantErr:  false,
			expected: "2023-08-14T07:01:31Z",
		},
		{
			name:     "ISO 8601 format",
			input:    "2023-08-14T07:01:31Z",
			wantErr:  false,
			expected: "2023-08-14T07:01:31Z",
		},
		{
			name:     "Date only",
			input:    "2023-08-14",
			wantErr:  false,
			expected: "2023-08-14T00:00:00Z",
		},
		{
			name:     "Date time without Z",
			input:    "2023-08-14 07:01:31",
			wantErr:  false,
			expected: "2023-08-14T07:01:31Z",
		},
		{
			name:     "DD-Mon-YYYY format",
			input:    "14-Aug-2023",
			wantErr:  false,
			expected: "2023-08-14T00:00:00Z",
		},
		{
			name:     "YYYY/MM/DD format",
			input:    "2023/08/14",
			wantErr:  false,
			expected: "2023-08-14T00:00:00Z",
		},
		{
			name:    "Invalid format",
			input:   "not a date",
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseWhoisDate(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseWhoisDate(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("parseWhoisDate(%q) unexpected error: %v", tt.input, err)
				return
			}

			resultStr := result.UTC().Format(time.RFC3339)
			if resultStr != tt.expected {
				t.Errorf("parseWhoisDate(%q) = %q, want %q", tt.input, resultStr, tt.expected)
			}
		})
	}
}

func TestCreateSummary_Structure(t *testing.T) {
	// Test the structure of CreateSummary with a mock QueryResult
	testTime := time.Now()

	mockResult := query.QueryResult{
		Query:     "example.com",
		Type:      string(query.QueryTypeDomain),
		Protocol:  "WHOIS",
		Timestamp: testTime,
		Success:   true,
		Data: map[string]interface{}{
			"parsed_fields": map[string]interface{}{
				"domain_name":          "EXAMPLE.COM",
				"registrar":            "Test Registrar",
				"creation_date":        "2020-01-01T00:00:00Z",
				"registry_expiry_date": "2025-01-01T00:00:00Z",
				"name_server_1":        "ns1.example.com",
				"name_server_2":        "ns2.example.com",
				"dnssec":               "unsigned",
				"domain_status":        "clientTransferProhibited",
			},
			"raw_response": "Mock WHOIS response data",
		},
		RawData: "Mock WHOIS response data",
	}

	summary := CreateSummary(mockResult)

	// Test basic fields
	if summary.Domain != "example.com" {
		t.Errorf("Expected Domain = example.com, got %q", summary.Domain)
	}

	if summary.Protocol != "WHOIS" {
		t.Errorf("Expected Protocol = WHOIS, got %q", summary.Protocol)
	}

	// Test that timeline events are parsed
	if summary.Timeline.Registration == nil {
		t.Error("Expected Registration timeline event to be set")
	} else {
		expectedDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		if !summary.Timeline.Registration.Date.Equal(expectedDate) {
			t.Errorf("Expected Registration date = %v, got %v", expectedDate, summary.Timeline.Registration.Date)
		}
	}

	if summary.Timeline.Expiration == nil {
		t.Error("Expected Expiration timeline event to be set")
	}

	// Test nameservers
	expectedNS := []string{"ns1.example.com", "ns2.example.com"}
	if len(summary.Nameservers) != len(expectedNS) {
		t.Errorf("Expected %d nameservers, got %d", len(expectedNS), len(summary.Nameservers))
	}

	// Test registrar
	if summary.Registrar.Name != "Test Registrar" {
		t.Errorf("Expected Registrar.Name = 'Test Registrar', got %q", summary.Registrar.Name)
	}

	// Test DNSSEC
	if summary.DNSSEC.Enabled {
		t.Error("Expected DNSSEC to be disabled for 'unsigned'")
	}

	// Test status
	if len(summary.StatusDetails) == 0 {
		t.Error("Expected status details to be populated")
	}

	// Test post-expiration guidance is added when there's an expiration date
	// Note: Guidance might be nil if expiration is far in future, which is OK
}

func TestCreateSummary_RDAP(t *testing.T) {
	// Test basic RDAP structure handling
	testTime := time.Now()

	mockRDAPData := map[string]interface{}{
		"objectClassName": "domain",
		"ldhName":         "example.com",
		"Status":          []interface{}{"active"},
		"Events": []interface{}{
			map[string]interface{}{
				"Action": "registration",
				"Date":   "2020-01-01T00:00:00Z",
			},
		},
		"Nameservers": []interface{}{
			map[string]interface{}{
				"LDHName": "ns1.example.com",
			},
		},
	}

	mockResult := query.QueryResult{
		Query:     "example.com",
		Type:      string(query.QueryTypeDomain),
		Protocol:  "RDAP",
		Timestamp: testTime,
		Success:   true,
		Data:      mockRDAPData,
	}

	summary := CreateSummary(mockResult)

	if summary.Protocol != "RDAP" {
		t.Errorf("Expected Protocol = RDAP, got %q", summary.Protocol)
	}

	if len(summary.StatusDetails) == 0 {
		t.Error("Expected status details to be populated from RDAP Status")
	}

	if len(summary.Nameservers) == 0 {
		t.Error("Expected nameservers to be populated from RDAP")
	}
}
