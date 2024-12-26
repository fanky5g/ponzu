package search

import (
	"reflect"

	"github.com/fanky5g/ponzu/constants"
)

type Pagination struct {
	Count  int
	Offset int
}

type Search struct {
	Query      string
	SortOrder  constants.SortOrder
	Pagination *Pagination
}

type CustomizableSearchAttributes interface {
	GetSearchableAttributes() map[string]reflect.Type
}

type SearchIndexable interface {
	IndexContent() bool
}
