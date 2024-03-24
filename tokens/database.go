package tokens

type Repository string

var (
	AnalyticsRepositoryToken      Repository = "analytics"
	ConfigRepositoryToken         Repository = "config"
	ContentRepositoryToken        Repository = "content"
	UploadRepositoryToken         Repository = "upload"
	CredentialHashRepositoryToken Repository = "credential-hash"
	RecoveryKeyRepositoryToken    Repository = "recovery-key"
	UserRepositoryToken           Repository = "users"
)
