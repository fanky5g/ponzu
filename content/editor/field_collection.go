package editor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/views"
)

func FieldCollection(fieldName, label string, p interface{}, types map[string]func(interface{}, *FieldArgs, ...Field) []byte) []byte {
	scope := TagNameFromStructField(fieldName, p, nil)
	tmpl := `
		<div class="control-block __ponzu-field-collection ` + scope + `">
			<label class="active">` + label + `</label>
	`

	value := ValueByName(fieldName, p, nil)
	fieldCollections, ok := value.Interface().(content.FieldCollections)
	if !ok {
		panic(fmt.Sprintf("Ponzu: '%s' is not a valid FieldCollections type", value.Type()))
	}

	positionalPlaceHolder := "%pos%"
	parentFieldPath := fmt.Sprintf("%s.%s.Value", fieldName, positionalPlaceHolder)

	typeTemplateMap := make(map[string]string)
	for typeName := range fieldCollections.AllowedTypes() {
		var emptyType interface{}
		emptyType, err := makeTypeWithEmptyAllowedTypes(p, fieldName, typeName)
		if err != nil {
			panic(err)
		}

		fieldCollectionTemplate := types[typeName](
			emptyType,
			&FieldArgs{
				Parent:   parentFieldPath,
				TypeName: typeName,
			},
			Field{
				View: []byte(
					fmt.Sprintf(
						`<input name="%s" type="hidden" value="%s" />`,
						fmt.Sprintf("%s.%s.type", scope, positionalPlaceHolder),
						typeName,
					),
				),
			},
		)

		typeTemplateMap[typeName] = string(fieldCollectionTemplate)
	}

	templatesBytes, err := json.Marshal(typeTemplateMap)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse templates: %s", err))
	}

	data := fieldCollections.Data()
	for i := 0; i < len(data); i++ {
		typeName := data[i].Type

		fieldCollectionTemplate := types[typeName](
			p,
			&FieldArgs{
				Parent:   fmt.Sprintf("%s.%d.Value", fieldName, i),
				TypeName: typeName,
			},
			Field{
				View: []byte(
					fmt.Sprintf(
						`<input name="%s" type="hidden" value="%s" />`,
						fmt.Sprintf("%s.%d.type", scope, i),
						typeName,
					),
				),
			})

		tmpl += string(fieldCollectionTemplate)
	}

	blockSelector := "block-select"
	tmpl += string(getBlockSelector(blockSelector, types))
	tmpl += `</div>`

	script := &bytes.Buffer{}
	scriptTmpl := makeScript("field_collection")

	if err = scriptTmpl.Execute(script, struct {
		Scope                        string
		BlockSelector                string
		FieldCollectionTemplates     string
		PositionalPlaceholder        string
		FieldCollectionInputTypeName string
		NumItems                     int
	}{
		Scope:                        scope,
		BlockSelector:                blockSelector,
		FieldCollectionTemplates:     string(templatesBytes),
		PositionalPlaceholder:        positionalPlaceHolder,
		NumItems:                     len(data),
		FieldCollectionInputTypeName: fmt.Sprintf("%s.%s.type", scope, positionalPlaceHolder),
	}); err != nil {
		panic(err)
	}

	return append([]byte(tmpl), script.Bytes()...)
}

func makeTypeWithEmptyAllowedTypes(p interface{}, fieldName, typeName string) (interface{}, error) {
	pType := reflect.TypeOf(p)
	if pType.Kind() == reflect.Pointer {
		pType = pType.Elem()
	}

	emptyType := reflect.New(pType).Interface()
	value := ValueByName(fieldName, emptyType, nil)

	iface := reflect.New(value.Type().Elem()).Interface()
	fieldCollections, ok := (iface).(content.FieldCollections)
	if !ok {
		return nil, fmt.Errorf("ponzu: '%s' is not a valid FieldCollections type", value.Type())
	}

	allowedTypes := fieldCollections.AllowedTypes()
	t, ok := allowedTypes[typeName]
	if !ok {
		return nil, fmt.Errorf("invalid type %s", typeName)
	}

	fieldCollections.Add(content.FieldCollection{
		Type:  typeName,
		Value: t(),
	})

	value.Set(reflect.ValueOf(iface))

	return emptyType, nil
}

func getBlockSelector(
	name string,
	types map[string]func(interface{}, *FieldArgs, ...Field) []byte) []byte {
	options := make([]string, len(types))
	i := 0
	for k := range types {
		options[i] = k
		i += 1
	}

	sel := struct {
		Name     string
		Label    string
		Value    string
		Options  []string
		Selector string
	}{
		Label:    "Select a block...",
		Selector: name,
		Name:     name,
		Options:  options,
		Value:    "",
	}

	w := &bytes.Buffer{}
	if err := views.ExecuteTemplate(w, "select.gohtml", sel); err != nil {
		panic(err)
	}

	return w.Bytes()
}
