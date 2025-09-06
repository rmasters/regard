package domain

import "time"

// Summary represents a structured summary of domain information
type Summary struct {
	Domain         string            `json:"domain"`
	Status         string            `json:"status"`
	StatusDetails  []string          `json:"status_details,omitempty"`
	Protocol       string            `json:"protocol"`
	QueryType      string            `json:"query_type,omitempty"`
	Timeline       Timeline          `json:"timeline"`
	Nameservers    []string          `json:"nameservers"`
	DNSSEC         DNSSECInfo        `json:"dnssec"`
	Registrar      RegistrarInfo     `json:"registrar"`
	PostExpiration *ExpirationInfo   `json:"post_expiration,omitempty"`
	ASN            *ASNInfo          `json:"asn,omitempty"`
}

// Timeline represents important dates in a domain's lifecycle
type Timeline struct {
	Registration *TimelineEvent `json:"registration,omitempty"`
	LastUpdated  *TimelineEvent `json:"last_updated,omitempty"`
	Expiration   *TimelineEvent `json:"expiration,omitempty"`
}

// TimelineEvent represents a single event with date and human-readable format
type TimelineEvent struct {
	Date          time.Time `json:"date"`
	HumanReadable string    `json:"human_readable"`
}

// DNSSECInfo represents DNSSEC status information
type DNSSECInfo struct {
	Enabled bool   `json:"enabled"`
	Details string `json:"details,omitempty"`
}

// RegistrarInfo represents registrar information
type RegistrarInfo struct {
	Name string `json:"name"`
	ID   string `json:"id,omitempty"`
}

// ExpirationInfo provides guidance for expired domains
type ExpirationInfo struct {
	DaysExpired     int        `json:"days_expired"`
	AvailableDate   *time.Time `json:"available_date,omitempty"`
	GuidanceMessage string     `json:"guidance_message"`
}

// ASNInfo represents Autonomous System Number information
type ASNInfo struct {
	Number       string   `json:"number"`
	Name         string   `json:"name"`
	Description  string   `json:"description,omitempty"`
	Country      string   `json:"country,omitempty"`
	Organization string   `json:"organization"`
	Status       string   `json:"status,omitempty"`
	AbuseContact string   `json:"abuse_contact,omitempty"`
	Peers        []string `json:"peers,omitempty"`
}