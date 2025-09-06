package domain

import "testing"

func TestInterpretStatus(t *testing.T) {
	tests := []struct {
		name          string
		statusDetails []string
		rawData       string
		expected      string
	}{
		{
			name:          "Available domain - not found",
			statusDetails: []string{},
			rawData:       "Domain not found in registry",
			expected:      "available",
		},
		{
			name:          "Available domain - no match",
			statusDetails: []string{},
			rawData:       "No match for domain example.com",
			expected:      "available",
		},
		{
			name:          "Available domain - no data found",
			statusDetails: []string{},
			rawData:       "No data found for example.com",
			expected:      "available",
		},
		{
			name:          "Available domain - available for registration",
			statusDetails: []string{},
			rawData:       "Domain available for registration",
			expected:      "available",
		},
		{
			name:          "Active domain with active status",
			statusDetails: []string{"active"},
			rawData:       "Domain: example.com\nStatus: active",
			expected:      "active",
		},
		{
			name:          "Active domain with OK status",
			statusDetails: []string{"ok"},
			rawData:       "Domain: example.com\nStatus: ok",
			expected:      "active",
		},
		{
			name:          "Active domain with prohibited statuses",
			statusDetails: []string{"clientDeleteProhibited", "clientTransferProhibited"},
			rawData:       "Domain registered and protected",
			expected:      "active",
		},
		{
			name:          "Active domain mixed statuses",
			statusDetails: []string{"active", "clientDeleteProhibited"},
			rawData:       "Domain is active with restrictions",
			expected:      "active",
		},
		{
			name:          "Unknown status - empty details, substantial data",
			statusDetails: []string{},
			rawData: `Domain Name: GOOGLE.COM
Registry Domain ID: 2138514_DOMAIN_COM-VRSN
Registrar WHOIS Server: whois.markmonitor.com
Registrar URL: http://www.markmonitor.com
Updated Date: 2019-09-09T15:39:04Z
Creation Date: 1997-09-15T04:00:00Z
Registry Expiry Date: 2028-09-14T04:00:00Z
Registrar: MarkMonitor Inc.
Registrar IANA ID: 292
Domain Status: some_unusual_status
Name Server: NS1.GOOGLE.COM`,
			expected:      "unknown",
		},
		{
			name:          "Available - no entries found",
			statusDetails: []string{},
			rawData:       "no entries found",
			expected:      "available",
		},
		{
			name:          "Available - short response",
			statusDetails: []string{},
			rawData:       "Not found",
			expected:      "available",
		},
		{
			name:          "Unknown - no details, no data",
			statusDetails: []string{},
			rawData:       "",
			expected:      "unknown",
		},
		{
			name:          "Unknown - unrecognized status details",
			statusDetails: []string{"pending", "locked"},
			rawData:       "Domain has some unusual status",
			expected:      "unknown",
		},
		{
			name:          "Case insensitive matching",
			statusDetails: []string{"ACTIVE"},
			rawData:       "Domain has active status with some other text",
			expected:      "active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InterpretStatus(tt.statusDetails, tt.rawData)
			if result != tt.expected {
				t.Errorf("InterpretStatus(%v, %q) = %q, want %q", 
					tt.statusDetails, tt.rawData, result, tt.expected)
			}
		})
	}
}