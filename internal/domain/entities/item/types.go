package item

import (
	"errors"
)

const (
	typeNotRegistered = `Error:
There is no type registered for %[1]s

Add this to the file which defines %[1]s{} in the 'content' package:


	func init() {			
		item.Types["%[1]s"] = func() interface{} { return new(%[1]s) }
	}		
				

`
)

var (
	// ErrTypeNotRegistered means content type isn't registered (not found in Types map)
	ErrTypeNotRegistered = errors.New(typeNotRegistered)

	// Types is a map used to reference a type name to its actual Editable type
	// mainly for lookups in /admin route based utilities
	Types       map[string]EntityBuilder
	Definitions map[string]TypeDefinition
)

type EntityBuilder func() interface{}

func init() {
	Types = make(map[string]EntityBuilder)
	Definitions = make(map[string]TypeDefinition)
}
