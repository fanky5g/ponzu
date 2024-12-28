package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/auth"
)

type CredentialHashDocument struct {
	*auth.CredentialHash
}

func (document *CredentialHashDocument) Value() (interface{}, error) {
	return json.Marshal(document)
}

func (document *CredentialHashDocument) Scan(src interface{}) error {
	if byteSrc, ok := src.([]byte); ok {
		return json.Unmarshal(byteSrc, &document)
	}

	if stringSrc, ok := src.(string); ok {
		return json.NewDecoder(strings.NewReader(stringSrc)).Decode(&document)
	}

	return fmt.Errorf("unsupported type %T", src)
}

type CredentialHashModel struct{}

func (*CredentialHashModel) Name() string {
	return WrapPonzuModelNameSpace("credential_hashes")
}

func (*CredentialHashModel) NewEntity() interface{} {
	return new(auth.CredentialHash)
}

func (model *CredentialHashModel) ToDocument(entity interface{}) database.DocumentInterface {
	return &CredentialHashDocument{
		CredentialHash: entity.(*auth.CredentialHash),
	}
}
