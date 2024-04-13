package tokens

type Repository string

var (
	AnalyticsRequestsRepositoryToken Repository = "analytics_request"
	AnalyticsMetricsRepositoryToken  Repository = "analytics_metric"
	ConfigRepositoryToken            Repository = "config"
	SlugRepositoryToken              Repository = "slug"
	UploadRepositoryToken            Repository = "upload"
	CredentialHashRepositoryToken    Repository = "credential-hash"
	RecoveryKeyRepositoryToken       Repository = "recovery-key"
	UserRepositoryToken              Repository = "users"
)
