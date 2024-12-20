package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fanky5g/ponzu/internal/entities"
	"github.com/fanky5g/ponzu/models"
)

type UserDocument struct {
	*entities.User
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
	return new(entities.User)
}

func (model *UserModel) ToDocument(entity interface{}) models.DocumentInterface {
	return &UserDocument{
		User: entity.(*entities.User),
	}
}
