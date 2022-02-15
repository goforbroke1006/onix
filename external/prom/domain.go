package prom

import "time"

type ResultType string

const (
	ResultTypeMatrix = ResultType("matrix")
	ResultTypeVector = ResultType("vector")
	ResultTypeScalar = ResultType("scalar")
	ResultTypeString = ResultType("string")
)

type ResponseStatus string

const (
	ResponseStatusSuccess = ResponseStatus("success")
	ResponseStatusError   = ResponseStatus("error")
)

type QueryResponse struct {
	Status ResponseStatus `json:"status"`
	Data   struct {
		ResultType ResultType   `json:"resultType"`
		Result     []ResultItem `json:"result"`
	} `json:"data"`
	Warnings []string `json:"warnings"`
}

type ResultItem struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type QueryRangeResponse struct {
	Status ResponseStatus `json:"status"`
	Data   struct {
		ResultType ResultType `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

type ApiClient interface {

	// Query wraps call to https://prometheus.io/docs/prometheus/latest/querying/api/#expression-queries
	Query(query string, timestamp time.Time, timeout time.Duration) (*QueryResponse, error)

	// QueryRange wraps call to https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries
	QueryRange(query string, start, end time.Time, step, timeout time.Duration) (*QueryRangeResponse, error)
}
