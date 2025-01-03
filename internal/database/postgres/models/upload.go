package models

import (
	"encoding/json"
	"fmt"
	"github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/database"
	"strings"
)

type UploadDocument struct {
	*entities.Upload
}

func (document *UploadDocument) Value() (interface{}, error) {
	return json.Marshal(document)
}

func (document *UploadDocument) Scan(src interface{}) error {
	if byteSrc, ok := src.([]byte); ok {
		return json.Unmarshal(byteSrc, &document)
	}

	if stringSrc, ok := src.(string); ok {
		return json.NewDecoder(strings.NewReader(stringSrc)).Decode(&document)
	}

	return fmt.Errorf("unsupported type %T", src)
}

type UploadModel struct{}

func (*UploadModel) Name() string {
	return WrapPonzuModelNameSpace("upload")
}

func (*UploadModel) NewEntity() interface{} {
	return new(entities.Upload)
}

func (model *UploadModel) ToDocument(entity interface{}) database.DocumentInterface {
	return &UploadDocument{
		Upload: entity.(*entities.Upload),
	}
}
