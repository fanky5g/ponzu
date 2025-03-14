package analytics

const (
	MetricsRepositoryToken  = "analytics_metrics"
	RequestsRepositoryToken = "analytics_http_request_metadata"
)

type Metric struct {
	MetricID string `json:"metric_id"`
	Date     string `json:"date"`
	Total    int    `json:"total"`
	Unique   int    `json:"unique"`
}

func (*Metric) GetRepositoryToken() string {
	return MetricsRepositoryToken
}

func (*Metric) EntityName() string {
	return "AnalyticsMetric"
}

type HTTPRequestMetadata struct {
	RequestID  string `json:"request_id"`
	URL        string `json:"url"`
	Method     string `json:"http_method"`
	Origin     string `json:"origin"`
	Proto      string `json:"http_protocol"`
	RemoteAddr string `json:"ip_address"`
	Timestamp  int64  `json:"timestamp"`
	External   bool   `json:"external_content"`
}

func (*HTTPRequestMetadata) GetRepositoryToken() string {
	return RequestsRepositoryToken
}

func (*HTTPRequestMetadata) EntityName() string {
	return "AnalyticsHTTPRequestMetadata"
}
