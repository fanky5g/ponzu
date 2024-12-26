package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/auth"
)

type RecoveryKeyDocument struct {
	*auth.RecoveryKey
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
	return new(auth.RecoveryKey)
}

func (model *RecoveryKeyModel) ToDocument(entity interface{}) database.DocumentInterface {
	return &RecoveryKeyDocument{
		RecoveryKey: entity.(*auth.RecoveryKey),
	}
}
