package tokens

type Service string

var (
	StorageServiceToken                Service = "StorageService"
	ConfigServiceToken                 Service = "ConfigService"
	AnalyticsServiceToken              Service = "AnalyticsService"
	AuthServiceToken                   Service = "AuthService"
	ContentServiceToken                Service = "ContentService"
	ContentSearchServiceToken          Service = "ContentSearchService"
	UploadSearchServiceToken           Service = "UploadSearchService"
	TLSServiceToken                    Service = "TLSService"
	UserServiceToken                   Service = "UserService"
)
