package query

import (
	"strings"
	"testing"
	"time"
)

func TestParseWhoisData(t *testing.T) {
	sampleWhoisData := `Domain Name: EXAMPLE.COM
   Registry Domain ID: 2336799_DOMAIN_COM-VRSN
   Registrar WHOIS Server: whois.iana.org
   Registrar URL: http://res-dom.iana.org
   Updated Date: 2025-08-14T07:01:39Z
   Creation Date: 1995-08-14T04:00:00Z
   Registry Expiry Date: 2026-08-13T04:00:00Z
   Registrar: RESERVED-Internet Assigned Numbers Authority
   Registrar IANA ID: 376
   Domain Status: clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
   Domain Status: clientTransferProhibited https://icann.org/epp#clientTransferProhibited
   Domain Status: clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
   Name Server: A.IANA-SERVERS.NET
   Name Server: B.IANA-SERVERS.NET
   DNSSEC: signedDelegation

% This is a comment line
# This is another comment line

>>> Last update of whois database: 2025-09-06T14:10:36Z <<<`

	result := parseWhoisData(sampleWhoisData)

	// Check that parsed_fields exists
	fields, ok := result["parsed_fields"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected parsed_fields to be a map[string]interface{}")
	}

	// Check specific fields
	tests := []struct {
		key      string
		expected string
	}{
		{"domain_name", "EXAMPLE.COM"},
		{"creation_date", "1995-08-14T04:00:00Z"},
		{"registry_expiry_date", "2026-08-13T04:00:00Z"},
		{"registrar", "RESERVED-Internet Assigned Numbers Authority"},
		{"registrar_iana_id", "376"},
		{"dnssec", "signedDelegation"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			value, exists := fields[tt.key]
			if !exists {
				t.Errorf("Expected field %q to exist", tt.key)
				return
			}
			if value != tt.expected {
				t.Errorf("Expected %q = %q, got %q", tt.key, tt.expected, value)
			}
		})
	}

	// Check that raw_response is preserved
	rawResponse, ok := result["raw_response"].(string)
	if !ok {
		t.Fatal("Expected raw_response to be a string")
	}
	if rawResponse != sampleWhoisData {
		t.Error("raw_response should preserve original data")
	}

	// Check that comments are ignored
	if strings.Contains(rawResponse, "%") || strings.Contains(rawResponse, "#") {
		// Comments should be in raw but not in parsed fields
		for key := range fields {
			if strings.Contains(key, "%") || strings.Contains(key, "#") {
				t.Errorf("Comments should not be parsed as fields, found key: %q", key)
			}
		}
	}
}

func TestPerformWhoisQuery_Structure(t *testing.T) {
	// Test the structure of the returned QueryResult without making actual network calls
	// This tests the function structure, not the actual WHOIS lookup
	
	testQuery := "example.com"
	
	// We can't easily mock the whois.Whois call, so we'll test the structure
	// by calling the function and checking error handling
	result := PerformWhoisQuery(testQuery)
	
	// Basic structure checks
	if result.Query != testQuery {
		t.Errorf("Expected Query = %q, got %q", testQuery, result.Query)
	}
	
	if result.Protocol != "WHOIS" {
		t.Errorf("Expected Protocol = WHOIS, got %q", result.Protocol)
	}
	
	if result.Type != string(QueryTypeDomain) {
		t.Errorf("Expected Type = %q, got %q", string(QueryTypeDomain), result.Type)
	}
	
	// Check timestamp is recent (within last minute)
	now := time.Now()
	if result.Timestamp.After(now) || result.Timestamp.Before(now.Add(-time.Minute)) {
		t.Errorf("Timestamp seems incorrect: %v", result.Timestamp)
	}
	
	// Result should either be successful with data, or failed with error
	if result.Success {
		if result.Data == nil {
			t.Error("Expected Data to be set when Success is true")
		}
		if result.RawData == "" {
			t.Error("Expected RawData to be set when Success is true")
		}
	} else {
		if result.Error == "" {
			t.Error("Expected Error to be set when Success is false")
		}
	}
}