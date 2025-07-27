package editor

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

func InputRepeater(publicPath, fieldName string, p interface{}, attrs map[string]string, args *FieldArgs) []byte {
	value := ValueByName(fieldName, p, args)

	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		panic(fmt.Sprintf("Ponzu: Type '%s' for field '%s' not supported.", value.Type(), fieldName))
	}

	scope := TagNameFromStructField(fieldName, p, args)

	tmpl := `
		<div class="control-block __ponzu-repeat ` + scope + `">
			<label class="active">` + fieldName + `</label>
	`

	positionalPlaceHolder := makePositionalPlaceholder(fieldName)
	fieldArgs := &FieldArgs{
		Parent:                 fmt.Sprintf("%s.%s", fieldName, positionalPlaceHolder),
		PositionalPlaceHolders: []string{positionalPlaceHolder},
	}

	emptyType := makeEmptyType(p)
	if value.IsZero() {
		emptyType = p
	}

	if args != nil && args.Parent != "" {
		fieldArgs.Parent = fmt.Sprintf("%s.%s", args.Parent, fieldArgs.Parent)
		fieldArgs.PositionalPlaceHolders = append(
			fieldArgs.PositionalPlaceHolders,
			args.PositionalPlaceHolders...,
		)

		matches := ancestorIsFieldCollectionRegexp.FindStringSubmatch(args.Parent)
		if len(matches) > 0 {
			positionMatchIndex := ancestorIsFieldCollectionRegexp.SubexpIndex("Position")
			if positionMatchIndex == -1 {
				panic("Parent path is invalid")
			}

			fieldCollectionNameIndex := ancestorIsFieldCollectionRegexp.SubexpIndex("FieldCollectionName")
			if fieldCollectionNameIndex == -1 {
				panic("Parent path is invalid")
			}

			var err error
			position, err := strconv.Atoi(matches[positionMatchIndex])
			if err != nil {
				panic(err)
			}

			emptyType, err = makeFieldCollectionAtPosition(p, matches[fieldCollectionNameIndex], args.TypeName, position)
			if err != nil {
				panic(err)
			}
		}
	}

	elementClass := fmt.Sprintf("%s-element", fieldName)
	extendAttributeValue(attrs, "controlClass", elementClass)
	entryTemplate := Input("", emptyType, attrs, fieldArgs)

	script := &bytes.Buffer{}
	scriptTmpl := makeScript("input_repeater")
	if err := scriptTmpl.Execute(script, struct {
		Template              string
		NumItems              int
		Scope                 string
		CloneSelector         string
		PositionalPlaceholder string
		PublicPath            string
		EntityName            string
	}{
		Template:              string(entryTemplate),
		NumItems:              value.Len(),
		Scope:                 scope,
		CloneSelector:         fmt.Sprintf(".%s", elementClass),
		PositionalPlaceholder: positionalPlaceHolder,
		PublicPath:            publicPath,
		EntityName:            fieldName,
	}); err != nil {
		panic(err)
	}

	for i := 0; i < value.Len(); i++ {
		entryArgs := &FieldArgs{
			Parent:                 fmt.Sprintf("%s.%d", fieldName, i),
			PositionalPlaceHolders: []string{positionalPlaceHolder},
		}

		if args != nil && args.Parent != "" {
			entryArgs.Parent = fmt.Sprintf("%s.%s", args.Parent, entryArgs.Parent)
			entryArgs.PositionalPlaceHolders = append(
				entryArgs.PositionalPlaceHolders,
				args.PositionalPlaceHolders...,
			)
		}

		fieldTemplate := Input("", p, attrs, entryArgs)
		tmpl += string(fieldTemplate)
	}

	tmpl += `</div>`
	return append([]byte(tmpl), script.Bytes()...)
}
