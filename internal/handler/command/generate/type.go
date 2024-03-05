package generate

import (
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/domain/enum"
	"strings"
)

// blog title:string Author:string PostCategory:string content:string some_thing:int
func parseType(contentType enum.ContentType, args []string) (*item.TypeDefinition, error) {
	name, label := fieldName(args[0])

	t := &item.TypeDefinition{
		Name:  name,
		Label: label,
	}

	t.Initial = strings.ToLower(string(t.Name[0]))

	data := args[1:]
	if contentType == enum.TypeFieldCollection {
		contentBlocks := make([]item.ContentBlock, 0)
		for _, contentBlock := range data {
			contentBlockName, contentBlockLabel := fieldName(contentBlock)
			contentBlocks = append(contentBlocks, item.ContentBlock{
				TypeName: contentBlockName,
				Label:    contentBlockLabel,
			})
		}

		t.ContentBlocks = contentBlocks
	} else {
		for _, field := range data {
			f, err := parseField(field, t)
			if err != nil {
				return nil, err
			}

			f.Initial = t.Initial
			t.Fields = append(t.Fields, *f)
		}
	}

	return t, nil
}
