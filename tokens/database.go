package tokens

type RepositoryToken string

var (
	AnalyticsRequestsRepositoryToken RepositoryToken = "analytics_http_request_metadata"
	AnalyticsMetricsRepositoryToken  RepositoryToken = "analytics_metrics"
	ConfigRepositoryToken            RepositoryToken = "config"
	SlugRepositoryToken              RepositoryToken = "slugs"
	UploadRepositoryToken            RepositoryToken = "uploads"
	CredentialHashRepositoryToken    RepositoryToken = "credential_hashes"
	RecoveryKeyRepositoryToken       RepositoryToken = "recovery_keys"
	UserRepositoryToken              RepositoryToken = "users"
)
