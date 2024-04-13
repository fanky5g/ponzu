package entities

import "reflect"

type Searchable interface {
	GetSearchableAttributes() map[string]reflect.Type
	IndexContent() bool
}
