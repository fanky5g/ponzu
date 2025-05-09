package models

import (
	"encoding/json"
	"fmt"
	"github.com/fanky5g/ponzu/internal/content"
	"strings"

	"github.com/fanky5g/ponzu/database"
)

type SlugDocument struct {
	*content.Slug
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

func (*SlugModel) NewEntity() interface{} {
	return new(content.Slug)
}

func (model *SlugModel) ToDocument(entity interface{}) database.DocumentInterface {
	return &SlugDocument{
		Slug: entity.(*content.Slug),
	}
}
