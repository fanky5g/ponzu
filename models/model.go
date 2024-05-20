package models

import (
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/google/uuid"
	"time"
)

const ponzuModelNameSpace = "ponzu"

type Model struct {
	ID        uuid.UUID         `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt time.Time         `json:"deleted_at"`
	Document  DocumentInterface `json:"document"`
}

type ModelInterface interface {
	Name() string
	ToDocument(entity interface{}) DocumentInterface
	NewEntity() content.Entity
}

func WrapPonzuModelNameSpace(name tokens.RepositoryToken) string {
	return fmt.Sprintf("%s_%s", ponzuModelNameSpace, name)
}

func GetPonzuModels() []ModelInterface {
	return []ModelInterface{
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
