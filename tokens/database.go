package tokens

type Repository string

var (
	AnalyticsRequestsRepositoryToken Repository = "analytics_http_request_metadata"
	AnalyticsMetricsRepositoryToken  Repository = "analytics_metrics"
	ConfigRepositoryToken            Repository = "config"
	SlugRepositoryToken              Repository = "slugs"
	UploadRepositoryToken            Repository = "uploads"
	CredentialHashRepositoryToken    Repository = "credential_hashes"
	RecoveryKeyRepositoryToken       Repository = "recovery_keys"
	UserRepositoryToken              Repository = "users"
)
