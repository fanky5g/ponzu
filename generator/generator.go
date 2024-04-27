package generator

import "github.com/pkg/errors"

var (
	Generators = make([]Generator, 0)
)

type Generator interface {
	Generate(definition *TypeDefinition, writer Writer) error
	Initialize(definition *TypeDefinition, writer Writer) error
}

func Register(generator Generator) {
	Generators = append(Generators, generator)
}

func Run(itemType Type, arguments []string) error {
	typeDefinition, err := NewTypeDefinition(itemType, arguments)
	if err != nil {
		return err
	}

	for _, generator := range Generators {
		if err = generator.Initialize(typeDefinition, FormattedSourceWriter); err != nil {
			return errors.Wrap(err, "Failed to initialize generator.")
		}

		if err = generator.Generate(typeDefinition, FormattedSourceWriter); err != nil {
			return errors.Wrap(err, "Failed to generate.")
		}
	}

	return nil
}
