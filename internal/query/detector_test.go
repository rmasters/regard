package query

import "testing"

func TestDetectQueryType(t *testing.T) {
	tests := []struct {
		query    string
		expected QueryType
	}{
		// Domain tests
		{"example.com", QueryTypeDomain},
		{"sub.example.org", QueryTypeDomain},
		{"very.long.subdomain.example.co.uk", QueryTypeDomain},
		{"test-domain.net", QueryTypeDomain},
		{"xn--e1afmkfd.xn--p1ai", QueryTypeDomain}, // IDN domain

		// IPv4 tests
		{"8.8.8.8", QueryTypeIP},
		{"192.168.1.1", QueryTypeIP},
		{"10.0.0.1", QueryTypeIP},
		{"255.255.255.255", QueryTypeIP},
		{"0.0.0.0", QueryTypeIP},

		// IPv6 tests
		{"2001:4860:4860::8888", QueryTypeIP},
		{"::1", QueryTypeIP},
		{"fe80::1", QueryTypeIP},
		{"2001:db8::1", QueryTypeIP},

		// ASN tests
		{"AS15169", QueryTypeASN},
		{"as13335", QueryTypeASN},
		{"AS1", QueryTypeASN},
		{"AS999999", QueryTypeASN},

		// Edge cases - domains that look like IPs but have long segments
		{"192.168.1.reallylong", QueryTypeDomain},
		{"8.8.8.example", QueryTypeDomain},

		// Default to domain for ambiguous cases
		{"localhost", QueryTypeDomain},
		{"test", QueryTypeDomain},
		{"", QueryTypeDomain},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			result := DetectQueryType(tt.query)
			if result != tt.expected {
				t.Errorf("DetectQueryType(%q) = %v, want %v", tt.query, result, tt.expected)
			}
		})
	}
}
