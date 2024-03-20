package parser

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/generator/types"
)

type Parser interface {
	ParseField(raw string, gt *types.TypeDefinition) (*types.Field, error)
	ParseTypeDefinition(contentType content.Type, args []string) (*types.TypeDefinition, error)
}

type parser struct {
	types content.Types
}

func New(types content.Types) (Parser, error) {
	return &parser{types: types}, nil
}
