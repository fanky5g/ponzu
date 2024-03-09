package editor

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/internal/util"
	"reflect"
)

func NestedRepeater(fieldName string, p interface{}, m func(v interface{}, f *FieldArgs) (string, []Field)) []byte {
	value := ValueByName(fieldName, p, nil)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		panic(fmt.Sprintf("Ponzu: Type '%s' for field '%s' not supported.", value.Type(), fieldName))
	}

	scope := TagNameFromStructField(fieldName, p, nil)

	tmpl := `
		<div class="col s12 __ponzu-repeat ` + scope + `">
			<label class="active">` + fieldName + `</label>
	`

	positionalPlaceHolder := "%pos%"
	fieldArgs := &FieldArgs{
		Parent: fmt.Sprintf("%s.%s", fieldName, positionalPlaceHolder),
	}

	arrayTypeName, fields := m(p, fieldArgs)
	fieldArgs.TypeName = arrayTypeName
	emptyEntryTemplate := Nested("", p, fieldArgs, fields...)

	script := &bytes.Buffer{}
	scriptTmpl := util.MakeScript("nested_repeater")
	if err := scriptTmpl.Execute(script, struct {
		Template              string
		NumItems              int
		Scope                 string
		InputSelector         string
		CloneSelector         string
		PositionalPlaceholder string
	}{
		Template:              string(emptyEntryTemplate),
		NumItems:              value.Len(),
		Scope:                 scope,
		CloneSelector:         fmt.Sprintf(".%s", arrayTypeName),
		InputSelector:         "input",
		PositionalPlaceholder: positionalPlaceHolder,
	}); err != nil {
		panic(err)
	}

	for i := 0; i < value.Len(); i++ {
		_, fields = m(p, &FieldArgs{
			Parent: fmt.Sprintf("%s.%d", fieldName, i),
		})

		fieldTemplate := Nested("", p, fieldArgs, fields...)
		tmpl += string(fieldTemplate)
	}

	tmpl += `</div>`
	return append([]byte(tmpl), script.Bytes()...)
}
