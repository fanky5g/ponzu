package editor

import (
	"bytes"
	"fmt"
	"reflect"
)

var positionalPlaceHolder = "%pos%"

type NestedFieldGenerator func(v interface{}, f *FieldArgs) (string, []Field)

func NestedRepeater(fieldName string, p interface{}, nestedFieldGenerator NestedFieldGenerator) []byte {
	value := ValueByName(fieldName, p, nil)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		panic(fmt.Sprintf("Ponzu: Type '%s' for field '%s' not supported.", value.Type(), fieldName))
	}

	scope := TagNameFromStructField(fieldName, p, nil)

	tmpl := `
		<div class="control-block __ponzu-nested __ponzu-repeat ` + scope + `">
			<label class="active">` + fieldName + `</label>
	`

	fieldArgs := &FieldArgs{
		Parent: fmt.Sprintf("%s.%s", fieldName, positionalPlaceHolder),
	}

	arrayTypeName, entryTemplate := generateNestedTemplate(nestedFieldGenerator, p, fieldArgs)

	script := &bytes.Buffer{}
	scriptTmpl := makeScript("nested_repeater")
	if err := scriptTmpl.Execute(script, struct {
		Template              string
		NumItems              int
		Scope                 string
		InputSelector         string
		CloneSelector         string
		PositionalPlaceholder string
	}{
		Template:              entryTemplate,
		NumItems:              value.Len(),
		Scope:                 scope,
		CloneSelector:         fmt.Sprintf(".%s", arrayTypeName),
		InputSelector:         "input",
		PositionalPlaceholder: positionalPlaceHolder,
	}); err != nil {
		panic(err)
	}

	for i := 0; i < value.Len(); i++ {
		_, fields := nestedFieldGenerator(p, &FieldArgs{
			Parent: fmt.Sprintf("%s.%d", fieldName, i),
		})

		fieldTemplate := Nested("", p, fieldArgs, fields...)
		tmpl += string(fieldTemplate)
	}

	tmpl += `</div>`
	return append([]byte(tmpl), script.Bytes()...)
}

func generateNestedTemplate(nestedFieldGenerator NestedFieldGenerator, entity interface{}, fieldArgs *FieldArgs) (string, string) {
	emptyType := makeEmptyType(entity)

	arrayTypeName, fields := nestedFieldGenerator(emptyType, fieldArgs)
	fieldArgs.TypeName = arrayTypeName

	return arrayTypeName, string(Nested("", emptyType, fieldArgs, fields...))
}
