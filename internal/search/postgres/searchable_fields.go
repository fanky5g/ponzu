package postgres

import (
	"fmt"
	"reflect"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/search"
	"github.com/fanky5g/ponzu/util"
)

var searchableFieldsStore = make(map[string][]string)

// getSearchableFields returns fields that are supported for search
func getSearchableFields(entity interface{}) ([]string, error) {
	entityInterface, ok := entity.(content.Entity)
	if !ok {
		return nil, ErrInvalidSearchEntity
	}

	if searchableFields, ok := searchableFieldsStore[entityInterface.EntityName()]; ok {
		return searchableFields, nil
	}

	v := reflect.Indirect(reflect.ValueOf(entity))
	t := v.Type()

	var searchableFields []string
	searchableAttributes, ok := entity.(search.CustomizableSearchAttributes)
	if ok {
		for attribute, attributeType := range searchableAttributes.GetSearchableAttributes() {
			if attributeType.Kind() != reflect.String {
				return nil, fmt.Errorf("%s is not supported for search", attributeType.Kind())
			}

			field := v.FieldByName(attribute)
			if !field.IsValid() {
				field = util.FieldByJSONTagName(entity, attribute)
			}

			if !field.IsValid() {
				return nil, fmt.Errorf("invalid field %s", attribute)
			}

			searchableFields = append(searchableFields, attribute)
		}

		searchableFieldsStore[entityInterface.EntityName()] = searchableFields
		return searchableFields, nil
	}

	for i := 0; i < v.NumField(); i++ {
		structField := t.Field(i)
		field := v.Field(i)

		if field.Kind() == reflect.String {
			fieldName := structField.Name
			if jsonTag, ok := structField.Tag.Lookup("json"); ok {
				fieldName = jsonTag
			}

			if fieldName != "-" {
				searchableFields = append(searchableFields, fieldName)
			}

		}
	}

	searchableFieldsStore[entityInterface.EntityName()] = searchableFields
	return searchableFields, nil
}
