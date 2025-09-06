package output

import (
	"testing"
	"time"

	"regard/internal/domain"
)

func TestStripAnsiCodes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No ANSI codes",
			input:    "plain text",
			expected: "plain text",
		},
		{
			name:     "Bold text",
			input:    "\033[1mBold\033[0m",
			expected: "Bold",
		},
		{
			name:     "Colored text",
			input:    "\033[32mgreen\033[0m",
			expected: "green",
		},
		{
			name:     "Multiple colors",
			input:    "\033[31mred\033[0m and \033[34mblue\033[0m",
			expected: "red and blue",
		},
		{
			name:     "Complex formatting",
			input:    "\033[1;32mBold Green\033[0m",
			expected: "Bold Green",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripAnsiCodes(tt.input)
			if result != tt.expected {
				t.Errorf("stripAnsiCodes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetTerminalWidth(t *testing.T) {
	// This test mainly ensures the function doesn't panic and returns a reasonable value
	width := getTerminalWidth()
	
	if width < 20 || width > 1000 {
		t.Errorf("getTerminalWidth() = %d, expected reasonable value between 20-1000", width)
	}
}

func TestOutputSummary_Structure(t *testing.T) {
	// Test that OutputSummary doesn't panic with various domain summaries
	now := time.Now()
	
	tests := []struct {
		name    string
		summary domain.Summary
	}{
		{
			name: "Active domain",
			summary: domain.Summary{
				Domain:   "example.com",
				Status:   "active", 
				Protocol: "RDAP",
				Timeline: domain.Timeline{
					Registration: &domain.TimelineEvent{
						Date:          now.Add(-365 * 24 * time.Hour),
						HumanReadable: "1 year ago",
					},
					Expiration: &domain.TimelineEvent{
						Date:          now.Add(365 * 24 * time.Hour),
						HumanReadable: "in 1 year",
					},
				},
				Nameservers: []string{"ns1.example.com", "ns2.example.com"},
				DNSSEC: domain.DNSSECInfo{
					Enabled: true,
					Details: "signed",
				},
				Registrar: domain.RegistrarInfo{
					Name: "Test Registrar",
					ID:   "123",
				},
				StatusDetails: []string{"active", "clientTransferProhibited"},
			},
		},
		{
			name: "Available domain", 
			summary: domain.Summary{
				Domain:   "available.example",
				Status:   "available",
				Protocol: "WHOIS",
			},
		},
		{
			name: "Expired domain with guidance",
			summary: domain.Summary{
				Domain:   "expired.example", 
				Status:   "expired",
				Protocol: "WHOIS",
				Timeline: domain.Timeline{
					Expiration: &domain.TimelineEvent{
						Date:          now.Add(-30 * 24 * time.Hour),
						HumanReadable: "30 days ago",
					},
				},
				PostExpiration: &domain.ExpirationInfo{
					DaysExpired:     30,
					GuidanceMessage: "Domain is in renewal grace period",
				},
			},
		},
		{
			name: "Minimal domain info",
			summary: domain.Summary{
				Domain:   "minimal.example",
				Status:   "unknown",
				Protocol: "WHOIS",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output by using a test that just ensures no panic
			// In a real scenario, you might want to capture stdout/stderr
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("OutputSummary panicked: %v", r)
				}
			}()
			
			// Test both with and without color
			OutputSummary(tt.summary, true)
			OutputSummary(tt.summary, false)
		})
	}
}

func TestOutputSummary_AvailableDomain(t *testing.T) {
	// Test specific behavior for available domains
	summary := domain.Summary{
		Domain:   "available.test",
		Status:   "available",
		Protocol: "RDAP",
	}
	
	// We can't easily capture stdout, but we can test that it doesn't panic
	// and that the function completes successfully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("OutputSummary panicked for available domain: %v", r)
		}
	}()
	
	OutputSummary(summary, false) // No color for predictable output
}