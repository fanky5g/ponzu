package content

import (
	generatorTypes "github.com/fanky5g/ponzu/content/generator/types"
)

type (
	Builder func() interface{}
	Types   struct {
		Content          map[string]Builder
		FieldCollections map[string]Builder
		Definitions      map[string]generatorTypes.TypeDefinition
	}
)
