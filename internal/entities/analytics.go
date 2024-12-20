package entities

type AnalyticsMetric struct {
	MetricID string `json:"metric_id"`
	Date     string `json:"date"`
	Total    int    `json:"total"`
	Unique   int    `json:"unique"`
}

func (*AnalyticsMetric) GetRepositoryToken() string {
	return "analytics_metrics"
}

func (*AnalyticsMetric) EntityName() string {
	return "AnalyticsMetric"
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

func (*AnalyticsHTTPRequestMetadata) GetRepositoryToken() string {
	return "analytics_http_request_metadata"
}

func (*AnalyticsHTTPRequestMetadata) EntityName() string {
	return "AnalyticsHTTPRequestMetadata"
}
