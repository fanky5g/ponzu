package generator

import (
	"github.com/fanky5g/ponzu/generator"
	"github.com/pkg/errors"
)

func Run(generators []generator.Generator, definition *generator.TypeDefinition, w generator.Writer) error {
	for _, generator := range generators {
		if err := generator.Initialize(definition, w); err != nil {
			return errors.Wrap(err, "Failed to initialize generator.")
		}

		if err := generator.Generate(definition, w); err != nil {
			return errors.Wrap(err, "Failed to generate.")
		}
	}

	return nil
}
