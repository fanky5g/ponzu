package generator

import (
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/generator"
	"github.com/pkg/errors"
)

type ViewScope struct {
	Definition   *generator.TypeDefinition
	Target       generator.Target
	TemplatesDir string
	ContentTypes content.Types
}

func (scope *ViewScope) Fields(parent *Field) ([]Field, error) {
	switch scope.Definition.Type {
	case generator.Content:
		fallthrough
	case generator.Plain:
		fields := make([]Field, len(scope.Definition.Blocks))
		for i, block := range scope.Definition.Blocks {
			field := mapBlockToField(scope.ContentTypes, block)
			if err := field.Validate(); err != nil {
				return nil, errors.Wrap(err, "Failed to validate")
			}

			field.Scope = scope
			field.Parent = parent
			if field.IsNested {
				t, ok := scope.ContentTypes.Definitions[field.ReferenceName]
				if !ok {
					return nil, fmt.Errorf("no definition matched for %s type", field.Name)
				}

				nestedScope := newViewScope(&t, scope.ContentTypes, scope.Target, scope.TemplatesDir)

				var err error
				field.Children, err = nestedScope.Fields(field)
				if err != nil {
					return nil, err
				}
			}

			fields[i] = *field
		}

		return fields, nil
	default:
		return nil, errors.New("unsupported content type")
	}
}

func newViewScope(
	definition *generator.TypeDefinition,
	contentTypes content.Types,
	target generator.Target,
	templateDir string,
) *ViewScope {
	return &ViewScope{
		Target:       target,
		Definition:   definition,
		TemplatesDir: templateDir,
		ContentTypes: contentTypes,
	}
}
