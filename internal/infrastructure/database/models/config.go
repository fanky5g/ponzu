package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fanky5g/ponzu/internal/entities"
	"github.com/fanky5g/ponzu/models"
)

type ConfigDocument struct {
	*entities.Config
}

func (document *ConfigDocument) Value() (interface{}, error) {
	return json.Marshal(document)
}

func (document *ConfigDocument) Scan(src interface{}) error {
	if byteSrc, ok := src.([]byte); ok {
		return json.Unmarshal(byteSrc, &document)
	}

	if stringSrc, ok := src.(string); ok {
		return json.NewDecoder(strings.NewReader(stringSrc)).Decode(&document)
	}

	return fmt.Errorf("unsupported type %T", src)
}

type ConfigModel struct{}

func (*ConfigModel) Name() string {
	return WrapPonzuModelNameSpace("config")
}

func (*ConfigModel) NewEntity() interface{} {
	return new(entities.Config)
}

func (model *ConfigModel) ToDocument(entity interface{}) models.DocumentInterface {
	return &ConfigDocument{
		Config: entity.(*entities.Config),
	}
}
