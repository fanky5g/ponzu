package entities

import "reflect"

type Searchable interface {
	GetSearchableAttributes() map[string]reflect.Type
	IndexContent() bool
}

type CustomizableSearchAttributes interface {
	GetSearchableAttributes() map[string]reflect.Type
}

type SearchIndexable interface {
	IndexContent() bool
}
