package models

import (
	"fmt"

	"github.com/fanky5g/ponzu/database"
)

const ponzuModelNameSpace = "ponzu"

func WrapPonzuModelNameSpace(name string) string {
	return fmt.Sprintf("%s_%s", ponzuModelNameSpace, name)
}

func GetPonzuModels() []database.ModelInterface {
	return []database.ModelInterface{
		new(UserModel),
		new(UploadModel),
		new(SlugModel),
		new(RecoveryKeyModel),
		new(CredentialHashModel),
		new(ConfigModel),
		new(AnalyticsMetricModel),
		new(AnalyticsHTTPRequestMetadataModel),
	}
}
