package search

import (
	"github.com/fanky5g/ponzu/internal/constants"
	"reflect"
)

type Search struct {
	Query     string
	SortOrder constants.SortOrder
	Count     int
	Offset    int
}

type SearchInterface interface {
	Update(id string, data interface{}) error
	Delete(entityName, entityId string) error
	Search(entityDefinition interface{}, query string, count, offset int) ([]interface{}, error)
	SearchWithPagination(entityDefinition interface{}, query string, count, offset int) ([]interface{}, int, error)
}

type CustomizableSearchAttributes interface {
	GetSearchableAttributes() map[string]reflect.Type
}

type SearchIndexable interface {
	IndexContent() bool
}
