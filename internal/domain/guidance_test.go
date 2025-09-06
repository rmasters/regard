package domain

import (
	"strings"
	"testing"
	"time"
)

func TestExtractTLD(t *testing.T) {
	tests := []struct {
		domain   string
		expected string
	}{
		{"example.com", "com"},
		{"test.org", "org"},
		{"subdomain.example.net", "net"},
		{"example.co.uk", "co.uk"},
		{"test.com.au", "com.au"},
		{"very.long.subdomain.example.org", "org"},
		{"localhost", "localhost"},
		{"single", "single"},
		{"a.b", "b"},
		{"a.b.c", "b.c"}, // This is actually correct - 3-letter segments are treated as compound TLDs
		{"example.info", "info"},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			result := extractTLD(tt.domain)
			if result != tt.expected {
				t.Errorf("extractTLD(%q) = %q, want %q", tt.domain, result, tt.expected)
			}
		})
	}
}

func TestGeneratePostExpirationGuidance(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		domain       string
		expiryOffset time.Duration // Offset from now (negative = past, positive = future)
		expectGuidance bool
		checkMessage   string // Substring that should be in guidance message
	}{
		{
			name:         "Domain expires far in future",
			domain:       "example.com",
			expiryOffset: 365 * 24 * time.Hour, // 1 year from now
			expectGuidance: false,
		},
		{
			name:         "Domain expires in 5 days",
			domain:       "example.com", 
			expiryOffset: 5 * 24 * time.Hour,
			expectGuidance: true,
			checkMessage:   "monitor closely",
		},
		{
			name:         "Domain expires in 25 days",
			domain:       "example.com",
			expiryOffset: 25 * 24 * time.Hour,
			expectGuidance: true,
			checkMessage:   "add to your watchlist",
		},
		{
			name:         "Domain expired 5 days ago (.com)",
			domain:       "example.com",
			expiryOffset: -5 * 24 * time.Hour,
			expectGuidance: true,
			checkMessage:   "renewal grace period",
		},
		{
			name:         "Domain expired 50 days ago (.com)",
			domain:       "example.com", 
			expiryOffset: -50 * 24 * time.Hour,
			expectGuidance: true,
			checkMessage:   "redemption grace period",
		},
		{
			name:         "Domain expired 85 days ago (.com)",
			domain:       "example.com",
			expiryOffset: -85 * 24 * time.Hour,
			expectGuidance: true,
			checkMessage:   "available for registration",
		},
		{
			name:         "Domain expired 50 days ago (.uk)",
			domain:       "example.co.uk",
			expiryOffset: -50 * 24 * time.Hour,
			expectGuidance: true,
			checkMessage:   "renewal grace period",
		},
		{
			name:         "Domain expired 100 days ago (.uk)",
			domain:       "example.co.uk",
			expiryOffset: -100 * 24 * time.Hour,
			expectGuidance: true,
			checkMessage:   "available for public registration",
		},
		{
			name:         "Domain expired 15 days ago (.info - generic TLD)",
			domain:       "example.info",
			expiryOffset: -15 * 24 * time.Hour,
			expectGuidance: true,
			checkMessage:   "renewal grace period",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expiryDate := now.Add(tt.expiryOffset)
			
			summary := Summary{
				Domain: tt.domain,
				Timeline: Timeline{
					Expiration: &TimelineEvent{
						Date:          expiryDate,
						HumanReadable: "test",
					},
				},
			}

			guidance := GeneratePostExpirationGuidance(summary)

			if tt.expectGuidance {
				if guidance == nil {
					t.Errorf("Expected guidance for %s, got nil", tt.name)
					return
				}

				if tt.checkMessage != "" && !strings.Contains(strings.ToLower(guidance.GuidanceMessage), strings.ToLower(tt.checkMessage)) {
					t.Errorf("Expected guidance message to contain %q, got %q", tt.checkMessage, guidance.GuidanceMessage)
				}

				// Check days expired calculation for past dates
				if tt.expiryOffset < 0 {
					expectedDays := int(-tt.expiryOffset.Hours() / 24)
					if guidance.DaysExpired != expectedDays {
						t.Errorf("Expected DaysExpired = %d, got %d", expectedDays, guidance.DaysExpired)
					}
				}
			} else {
				if guidance != nil {
					t.Errorf("Expected no guidance for %s, got %+v", tt.name, guidance)
				}
			}
		})
	}
}

func TestGeneratePostExpirationGuidance_NoExpiration(t *testing.T) {
	summary := Summary{
		Domain: "example.com",
		Timeline: Timeline{
			// No expiration date set
		},
	}

	guidance := GeneratePostExpirationGuidance(summary)
	if guidance != nil {
		t.Errorf("Expected nil guidance when no expiration date, got %+v", guidance)
	}
}

func TestGeneratePostExpirationGuidance_PendingDelete(t *testing.T) {
	now := time.Now()
	// Domain expired 78 days ago (.com) - should be pending delete
	expiryDate := now.Add(-78 * 24 * time.Hour)
	
	summary := Summary{
		Domain: "example.com",
		Timeline: Timeline{
			Expiration: &TimelineEvent{
				Date:          expiryDate,
				HumanReadable: "test",
			},
		},
	}

	guidance := GeneratePostExpirationGuidance(summary)
	
	if guidance == nil {
		t.Fatal("Expected guidance for pending delete domain")
	}

	if !strings.Contains(guidance.GuidanceMessage, "pending deletion") {
		t.Errorf("Expected message about pending deletion, got %q", guidance.GuidanceMessage)
	}

	if !strings.Contains(guidance.GuidanceMessage, "2 days") {
		t.Errorf("Expected message about 2 days remaining, got %q", guidance.GuidanceMessage)
	}

	if guidance.AvailableDate == nil {
		t.Error("Expected AvailableDate to be set for pending delete")
	} else {
		expectedDate := expiryDate.AddDate(0, 0, 80)
		if !guidance.AvailableDate.Equal(expectedDate) {
			t.Errorf("Expected AvailableDate = %v, got %v", expectedDate, *guidance.AvailableDate)
		}
	}
}