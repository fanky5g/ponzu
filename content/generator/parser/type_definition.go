package parser

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/generator/types"
	"strings"
)

// ParseTypeDefinition e.g. blog title:string Author:string PostCategory:string entities:string some_thing:int
func (p *parser) ParseTypeDefinition(contentType content.Type, args []string) (*types.TypeDefinition, error) {
	name, label := fieldName(args[0])

	t := &types.TypeDefinition{
		Name:  name,
		Label: label,
	}

	t.Initial = strings.ToLower(string(t.Name[0]))

	data := args[1:]
	if contentType == content.TypeFieldCollection {
		contentBlocks := make([]types.ContentBlock, 0)
		for _, contentBlock := range data {
			contentBlockName, contentBlockLabel := fieldName(contentBlock)
			contentBlocks = append(contentBlocks, types.ContentBlock{
				TypeName: contentBlockName,
				Label:    contentBlockLabel,
			})
		}

		t.ContentBlocks = contentBlocks
	} else {
		for _, field := range data {
			f, err := p.ParseField(field, t)
			if err != nil {
				return nil, err
			}

			f.Initial = t.Initial
			t.Fields = append(t.Fields, *f)
		}
	}

	return t, nil
}
