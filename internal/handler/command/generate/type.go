package generate

import (
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"strings"
)

// blog title:string Author:string PostCategory:string content:string some_thing:int
func parseType(args []string) (*item.TypeDefinition, error) {
	name := fieldName(args[0])
	t := &item.TypeDefinition{
		Name:  name,
		Label: name,
	}
	t.Initial = strings.ToLower(string(t.Name[0]))

	fields := args[1:]
	for _, field := range fields {
		f, err := parseField(field, t)
		if err != nil {
			return nil, err
		}

		f.Initial = t.Initial
		t.Fields = append(t.Fields, *f)
	}

	return t, nil
}
