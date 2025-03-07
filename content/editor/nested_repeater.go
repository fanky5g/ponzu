package editor

import (
	"bytes"
	"fmt"
	"reflect"
)

var positionalPlaceHolder = "%pos%"

type NestedFieldGenerator func(v interface{}, f *FieldArgs) (string, []Field)

func NestedRepeater(fieldName string, p interface{}, args *FieldArgs, nestedFieldGenerator NestedFieldGenerator) []byte {
	value := ValueByName(fieldName, p, args)

	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		panic(fmt.Sprintf("Ponzu: Type '%s' for field '%s' not supported.", value.Type(), fieldName))
	}

	scope := TagNameFromStructField(fieldName, p, args)

	tmpl := `
		<div class="control-block __ponzu-nested __ponzu-repeat ` + scope + `">
			<label class="active">` + fieldName + `</label>
	`

	fieldArgs := &FieldArgs{
		Parent: fmt.Sprintf("%s.%s", fieldName, positionalPlaceHolder),
	}

	emptyType := makeEmptyType(p)
	if args != nil && args.Parent != "" {
		fieldArgs.Parent = fmt.Sprintf("%s.%s", args.Parent, fieldArgs.Parent)
		// we want to keep type representation intact
		// e.g. in cases where NestedRepeater is a child of FieldCollection
		emptyType = p
	}

	arrayTypeName, entryTemplate := generateNestedTemplate(nestedFieldGenerator, emptyType, fieldArgs)

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

func generateNestedTemplate(nestedFieldGenerator NestedFieldGenerator, emptyType interface{}, fieldArgs *FieldArgs) (string, string) {
	arrayTypeName, fields := nestedFieldGenerator(emptyType, fieldArgs)
	fieldArgs.TypeName = arrayTypeName

	return arrayTypeName, string(Nested("", emptyType, fieldArgs, fields...))
}
