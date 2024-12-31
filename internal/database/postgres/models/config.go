package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/config"
)

type ConfigDocument struct {
	*config.Config
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
	return new(config.Config)
}

func (model *ConfigModel) ToDocument(entity interface{}) database.DocumentInterface {
	return &ConfigDocument{
		Config: entity.(*config.Config),
	}
}
