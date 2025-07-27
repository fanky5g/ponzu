package editor

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

type NestedFieldGenerator func(v interface{}, f *FieldArgs) (string, []Field)

func NestedRepeater(publicPath, fieldName string, p interface{}, args *FieldArgs, nestedFieldGenerator NestedFieldGenerator) []byte {
	value := ValueByName(fieldName, p, args)

	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		panic(fmt.Sprintf("Ponzu: Type '%s' for field '%s' not supported.", value.Type(), fieldName))
	}

	scope := TagNameFromStructField(fieldName, p, args)
	tmpl := `
		<div class="control-block __ponzu-nested __ponzu-repeat ` + scope + `">
			<label class="active">` + fieldName + `</label>
	`

	positionalPlaceHolder := makePositionalPlaceholder(fieldName)
	parent := fmt.Sprintf("%s.%s", fieldName, positionalPlaceHolder)
	fieldArgs := &FieldArgs{
		Parent:                 parent,
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

	arrayTypeName, entryTemplate := generateNestedTemplate(nestedFieldGenerator, emptyType, fieldArgs)

	script := &bytes.Buffer{}
	scriptTmpl := makeScript("nested_repeater")
	if err := scriptTmpl.Execute(script, struct {
		Template              string
		NumItems              int
		Scope                 string
		CloneSelector         string
		PositionalPlaceholder string
		PublicPath            string
		EntityName            string
	}{
		Template:              entryTemplate,
		NumItems:              value.Len(),
		Scope:                 scope,
		CloneSelector:         fmt.Sprintf(".%s", arrayTypeName),
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

		_, fields := nestedFieldGenerator(p, entryArgs)

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
