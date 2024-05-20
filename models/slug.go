package models

import (
	"encoding/json"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/entities"
	"strings"
)

type SlugDocument struct {
	*entities.Slug
}

func (document *SlugDocument) Value() (interface{}, error) {
	return json.Marshal(document)
}

func (document *SlugDocument) Scan(src interface{}) error {
	if byteSrc, ok := src.([]byte); ok {
		return json.Unmarshal(byteSrc, &document)
	}

	if stringSrc, ok := src.(string); ok {
		return json.NewDecoder(strings.NewReader(stringSrc)).Decode(&document)
	}

	return fmt.Errorf("unsupported type %T", src)
}

type SlugModel struct{}

func (*SlugModel) Name() string {
	return WrapPonzuModelNameSpace("slugs")
}

func (*SlugModel) NewEntity() content.Entity {
	return new(entities.Slug)
}

func (model *SlugModel) ToDocument(entity interface{}) DocumentInterface {
	return &SlugDocument{
		Slug: entity.(*entities.Slug),
	}
}
