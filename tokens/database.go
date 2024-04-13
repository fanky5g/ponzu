package tokens

type Repository string

var (
	AnalyticsRequestsRepositoryToken Repository = "analytics_requests"
	AnalyticsMetricsRepositoryToken  Repository = "analytics_metrics"
	ConfigRepositoryToken            Repository = "config"
	ContentRepositoryToken           Repository = "content"
	UploadRepositoryToken            Repository = "upload"
	CredentialHashRepositoryToken    Repository = "credential-hash"
	RecoveryKeyRepositoryToken       Repository = "recovery-key"
	UserRepositoryToken              Repository = "users"
)
