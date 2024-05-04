package content

import "github.com/fanky5g/ponzu/generator"

type (
	Builder func() interface{}
	Types   struct {
		Content          map[string]Builder
		FieldCollections map[string]Builder
		Definitions      map[string]generator.TypeDefinition
	}
)
