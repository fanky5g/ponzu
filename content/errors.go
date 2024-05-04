package content

import (
	"errors"
)

const (
	typeNotRegistered = `Error:
	There is no type registered for %[1]s
	Add this to the file which defines %[1]s{} in the 'entities' package:
	func init() {			
		item.Types["%[1]s"] = func() interface{} { return new(%[1]s) }
	}
`
)

var (
	// ErrTypeNotRegistered means entities type isn't registered (not found in Types map)
	ErrTypeNotRegistered = errors.New(typeNotRegistered)
)
