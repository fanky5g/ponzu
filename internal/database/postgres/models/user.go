package models

import (
	"encoding/json"
	"fmt"
	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/auth"
	"strings"
)

type UserDocument struct {
	*auth.User
}

func (document *UserDocument) Value() (interface{}, error) {
	return json.Marshal(document)
}

func (document *UserDocument) Scan(src interface{}) error {
	if byteSrc, ok := src.([]byte); ok {
		return json.Unmarshal(byteSrc, &document)
	}

	if stringSrc, ok := src.(string); ok {
		return json.NewDecoder(strings.NewReader(stringSrc)).Decode(&document)
	}

	return fmt.Errorf("unsupported type %T", src)
}

type UserModel struct{}

func (*UserModel) Name() string {
	return WrapPonzuModelNameSpace("users")
}

func (*UserModel) NewEntity() interface{} {
	return new(auth.User)
}

func (model *UserModel) ToDocument(entity interface{}) database.DocumentInterface {
	return &UserDocument{
		User: entity.(*auth.User),
	}
}
