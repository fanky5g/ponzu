package tokens

type Driver string

var (
	ContentSearchClientInfrastructureToken Driver = "content-search-client"
	UploadSearchClientInfrastructureToken  Driver = "upload-search-client"
	StorageClientInfrastructureToken       Driver = "storage-client"
	AssetStorageClientInfrastructureToken  Driver = "asset-storage-client"
	DatabaseInfrastructureToken            Driver = "database"
)
