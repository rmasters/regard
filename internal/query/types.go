package query

import "time"

// QueryResult represents the result of a domain/IP/ASN query
type QueryResult struct {
	Query     string      `json:"query"`
	Type      string      `json:"type"`
	Protocol  string      `json:"protocol"`
	Timestamp time.Time   `json:"timestamp"`
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	RawData   string      `json:"raw_data,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// QueryType represents the type of query being performed
type QueryType string

const (
	QueryTypeDomain QueryType = "domain"
	QueryTypeIP     QueryType = "ip"
	QueryTypeASN    QueryType = "asn"
)