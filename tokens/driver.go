package tokens

type Driver string

var (
	SearchClientInfrastructureToken Driver = "search-client"
	StorageClientInfrastructureToken       Driver = "storage-client"
	AssetStorageClientInfrastructureToken  Driver = "asset-storage-client"
	DatabaseInfrastructureToken            Driver = "database"
)
