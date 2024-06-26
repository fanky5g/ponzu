package models

import (
	"encoding/json"
	"fmt"
	"github.com/fanky5g/ponzu/entities"
	"strings"
)

type UploadDocument struct {
	*entities.FileUpload
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
	return WrapPonzuModelNameSpace("uploads")
}

func (*UploadModel) NewEntity() interface{} {
	return new(entities.FileUpload)
}

func (model *UploadModel) ToDocument(entity interface{}) DocumentInterface {
	return &UploadDocument{
		FileUpload: entity.(*entities.FileUpload),
	}
}
