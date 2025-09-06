package query

import (
	"encoding/json"
	"time"

	"github.com/openrdap/rdap"
)

// PerformRDAPQuery executes an RDAP query for the given input
func PerformRDAPQuery(query string) QueryResult {
	result := QueryResult{
		Query:     query,
		Type:      string(DetectQueryType(query)),
		Protocol:  "RDAP",
		Timestamp: time.Now(),
	}

	client := &rdap.Client{}

	var response interface{}
	var err error

	switch result.Type {
	case string(QueryTypeDomain):
		response, err = client.QueryDomain(query)
	case string(QueryTypeIP):
		response, err = client.QueryIP(query)
	case string(QueryTypeASN):
		response, err = client.QueryAutnum(query)
	default:
		response, err = client.QueryDomain(query)
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result
	}

	result.Success = true
	result.Data = response

	// Also store raw JSON for fallback
	if rawBytes, err := json.Marshal(response); err == nil {
		result.RawData = string(rawBytes)
	}

	return result
}