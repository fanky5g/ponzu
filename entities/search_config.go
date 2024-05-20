package entities

import "reflect"

type CustomizableSearchAttributes interface {
	GetSearchableAttributes() map[string]reflect.Type
}

type SearchIndexable interface {
	IndexContent() bool
}

type CustomizableSearchAttributes interface {
	GetSearchableAttributes() map[string]reflect.Type
}

type SearchIndexable interface {
	IndexContent() bool
}
