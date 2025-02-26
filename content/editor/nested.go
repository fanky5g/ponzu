package editor

import "fmt"

func Nested(fieldName string, p interface{}, args *FieldArgs, fields ...Field) []byte {
	name := fieldName

	fmt.Println("Nested", fieldName)
	if name == "" && args != nil {
		name = args.TypeName
	}

	tmpl := `
		<fieldset class="control-block ` + name + `" name="` + name + `">
			<label>` + name + `</label>
	`

	for _, field := range fields {
		tmpl += string(field.View)
	}

	tmpl += `</fieldset>`
	return []byte(tmpl)
}
