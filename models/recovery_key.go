package models

import (
	"encoding/json"
	"fmt"
	"github.com/fanky5g/ponzu/entities"
	"strings"
)

type RecoveryKeyDocument struct {
	*entities.RecoveryKey
}

func (document *RecoveryKeyDocument) Value() (interface{}, error) {
	return json.Marshal(document)
}

func (document *RecoveryKeyDocument) Scan(src interface{}) error {
	if byteSrc, ok := src.([]byte); ok {
		return json.Unmarshal(byteSrc, &document)
	}

	if stringSrc, ok := src.(string); ok {
		return json.NewDecoder(strings.NewReader(stringSrc)).Decode(&document)
	}

	return fmt.Errorf("unsupported type %T", src)
}

type RecoveryKeyModel struct{}

func (*RecoveryKeyModel) Name() string {
	return WrapPonzuModelNameSpace("recovery_keys")
}

func (*RecoveryKeyModel) NewEntity() interface{} {
	return new(entities.RecoveryKey)
}

func (model *RecoveryKeyModel) ToDocument(entity interface{}) DocumentInterface {
	return &RecoveryKeyDocument{
		RecoveryKey: entity.(*entities.RecoveryKey),
	}
}
