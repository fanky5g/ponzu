package entities

import "github.com/fanky5g/ponzu/tokens"

type AnalyticsMetric struct {
	MetricID string `json:"metric_id"`
	Date     string `json:"date"`
	Total    int    `json:"total"`
	Unique   int    `json:"unique"`
}

func (*AnalyticsMetric) GetRepositoryToken() tokens.RepositoryToken {
	return tokens.AnalyticsMetricsRepositoryToken
}

type AnalyticsHTTPRequestMetadata struct {
	RequestID  string `json:"request_id"`
	URL        string `json:"url"`
	Method     string `json:"http_method"`
	Origin     string `json:"origin"`
	Proto      string `json:"http_protocol"`
	RemoteAddr string `json:"ip_address"`
	Timestamp  int64  `json:"timestamp"`
	External   bool   `json:"external_content"`
}

func (*AnalyticsHTTPRequestMetadata) GetRepositoryToken() tokens.RepositoryToken {
	return tokens.AnalyticsRequestsRepositoryToken
}
