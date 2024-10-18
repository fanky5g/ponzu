package editor

func Nested(fieldName string, p interface{}, args *FieldArgs, fields ...Field) []byte {
	name := fieldName

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
