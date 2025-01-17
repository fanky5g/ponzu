package references

import (
	"fmt"
	"github.com/fanky5g/ponzu/content/item"
	"reflect"
	"slices"
)

var ReferenceLoaderChunkSize = 20

type ReferenceLoaderInterface interface {
	GetReferences(entityName string, entityIds ...string) ([]interface{}, error)
}

type Store struct {
	data       map[string]interface{}
	references map[string][]string
}

func NewStore(entity interface{}) *Store {
	return &Store{
		references: buildReferences(entity),
		data:       make(map[string]interface{}),
	}
}

func (s *Store) PopulateReferences(referenceLoader ReferenceLoaderInterface) error {
	if len(s.references) == 0 {
		return nil
	}

	// TODO(B.B): will benefit from parallel execution
	for entityName, entityIds := range s.references {
		for chunk := range slices.Chunk(entityIds, ReferenceLoaderChunkSize) {
			data, err := referenceLoader.GetReferences(entityName, chunk...)
			if err != nil {
				return err
			}

			for i := range data {
				reference := data[i]
				identifiable, ok := reference.(item.Identifiable)
				if !ok {
					return fmt.Errorf("reference %s is not an Identifiable", reference)
				}

				s.data[s.key(entityName, identifiable.ItemID())] = reference
			}
		}
	}

	return nil
}

func (s *Store) GetReference(entityName, entityId string) interface{} {
	if entity, ok := s.data[s.key(entityName, entityId)]; ok {
		return entity
	}

	return nil
}

func (s *Store) key(entityName, entityId string) string {
	return fmt.Sprintf("%s:%s", entityName, entityId)
}

func buildReferences(entity interface{}) map[string][]string {
	if entity == nil {
		return nil
	}

	r := make(map[string][]string)
	buildReferenceMap(entity, r)

	return r
}

func buildReferenceMap(entity interface{}, collector map[string][]string) {
	v := reflect.ValueOf(entity)

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			buildReferenceMap(v.Index(i).Interface(), collector)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if !field.IsExported() {
				continue
			}

			fieldValue := v.Field(i)
			fieldTag := field.Tag

			if isComplexType(fieldValue.Kind()) {
				buildReferenceMap(fieldValue.Interface(), collector)
			}

			if referenceTag, ok := fieldTag.Lookup("reference"); ok {
				if str, ok := fieldValue.Interface().(string); ok {
					collector[referenceTag] = append(collector[referenceTag], str)
				}
			}
		}
	default:
	}
}

func isComplexType(kind reflect.Kind) bool {
	return kind == reflect.Pointer ||
		kind == reflect.Array ||
		kind == reflect.Slice ||
		kind == reflect.Struct ||
		kind == reflect.Interface
}
