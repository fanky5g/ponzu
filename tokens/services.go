package tokens

type Service string

var (
	UploadServiceToken        Service = "UploadService"
	ConfigServiceToken        Service = "ConfigService"
	AnalyticsServiceToken     Service = "AnalyticsService"
	AuthServiceToken          Service = "AuthService"
	ContentServiceToken       Service = "ContentService"
	ContentSearchServiceToken Service = "ContentSearchService"
	UploadSearchServiceToken  Service = "UploadSearchService"
	TLSServiceToken           Service = "TLSService"
	UserServiceToken          Service = "UserService"
)
