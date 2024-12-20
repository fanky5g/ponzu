package models

import (
	"fmt"

	"github.com/fanky5g/ponzu/models"
)

const ponzuModelNameSpace = "ponzu"

func WrapPonzuModelNameSpace(name string) string {
	return fmt.Sprintf("%s_%s", ponzuModelNameSpace, name)
}

func GetPonzuModels() []models.ModelInterface {
	return []models.ModelInterface{
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
