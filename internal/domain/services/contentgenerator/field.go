package contentgenerator

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	reservedFieldNames = map[string]string{
		"uuid":      "UUID",
		"item":      "Item",
		"id":        "ID",
		"slug":      "Slug",
		"timestamp": "Timestamp",
		"updated":   "Updated",
	}
)

func (gt *generator) ValidateField(field *item.Field) error {
	for jsonName, fieldName := range reservedFieldNames {
		if field.JSONName == jsonName || field.Name == fieldName {
			return fmt.Errorf("reserved field name: %s (%s)", jsonName, fieldName)
		}
	}

	return nil
}

// set the specified view inside the editor field for a generated field for a type
func (gt *generator) setFieldView(definition *item.TypeDefinition, index int, args ...string) error {
	var err error
	var tmpl *template.Template
	buf := &bytes.Buffer{}
	field := &definition.Fields[index]

	fieldGeneratorDataVar := "v"
	fieldGeneratorArgsVar := "nil"

	if len(args) > 0 {
		fieldGeneratorDataVar = args[0]
	}

	if len(args) > 1 {
		fieldGeneratorArgsVar = args[1]
	}

	var templateArg = struct {
		*item.Field
		FieldGeneratorDataVar string
		FieldGeneratorArgsVar string
		Fields                []item.Field
	}{
		Field:                 field,
		FieldGeneratorDataVar: fieldGeneratorDataVar,
		FieldGeneratorArgsVar: fieldGeneratorArgsVar,
	}

	tmplFromWithDelims := func(filename string, delim [2]string) (*template.Template, error) {
		if delim[0] == "" || delim[1] == "" {
			delim = [2]string{"{{", "}}"}
		}

		return template.New(filename).Delims(delim[0], delim[1]).ParseFiles(filepath.Join(gt.templateDir, filename))
	}

	optimizeFieldView(field)
	switch field.ViewType {
	case "checkbox":
		tmpl, err = tmplFromWithDelims("gen-checkbox.tmpl", [2]string{})
	case "custom":
		tmpl, err = tmplFromWithDelims("gen-custom.tmpl", [2]string{})
	case "file":
		tmpl, err = tmplFromWithDelims("gen-file.tmpl", [2]string{})
	case "hidden":
		tmpl, err = tmplFromWithDelims("gen-hidden.tmpl", [2]string{})
	case "input", "text":
		tmpl, err = tmplFromWithDelims("gen-input.tmpl", [2]string{})
	case "richtext":
		tmpl, err = tmplFromWithDelims("gen-richtext.tmpl", [2]string{})
	case "select":
		tmpl, err = tmplFromWithDelims("gen-select.tmpl", [2]string{})
	case "textarea":
		tmpl, err = tmplFromWithDelims("gen-textarea.tmpl", [2]string{})
	case "tags":
		tmpl, err = tmplFromWithDelims("gen-tags.tmpl", [2]string{})

	case "input-repeater":
		tmpl, err = tmplFromWithDelims("gen-input-repeater.tmpl", [2]string{})
	case "select-repeater":
		tmpl, err = tmplFromWithDelims("gen-select-repeater.tmpl", [2]string{})
	case "file-repeater":
		tmpl, err = tmplFromWithDelims("gen-file-repeater.tmpl", [2]string{})

	// use [[ and ]] as delimiters since reference html need to generate
	// display names containing {{ and }}
	case "reference":
		tmpl, err = tmplFromWithDelims("gen-reference.tmpl", [2]string{"[[", "]]"})
		if err != nil {
			return err
		}
	case "reference-repeater":
		tmpl, err = tmplFromWithDelims("gen-reference-repeater.tmpl", [2]string{"[[", "]]"})
		if err != nil {
			return err
		}
	case "nested":
		tmpl, err = tmplFromWithDelims("gen-nested.tmpl", [2]string{})
		t, ok := item.Definitions[field.ReferenceName]
		if !ok {
			return fmt.Errorf("no definition matched for %s type", field.Name)
		}

		for i := range t.Fields {
			t.Fields[i].Name = fmt.Sprintf("%s.%s", field.Name, t.Fields[i].Name)
			t.Fields[i].Initial = definition.Initial
			if err = gt.setFieldView(&t, i); err != nil {
				return err
			}
		}

		templateArg.Fields = t.Fields
	case "nested-repeater":
		tmpl, err = tmplFromWithDelims("gen-nested-repeater.tmpl", [2]string{})
		t, ok := item.Definitions[field.ReferenceName]
		if !ok {
			return fmt.Errorf("no definition matched for %s type", field.Name)
		}

		templateArg.FieldGeneratorArgsVar = "args"
		for i := range t.Fields {
			t.Fields[i].Initial = templateArg.FieldGeneratorDataVar

			if err = gt.setFieldView(
				&t,
				i,
				templateArg.FieldGeneratorDataVar,
				templateArg.FieldGeneratorArgsVar,
			); err != nil {
				return err
			}
		}

		templateArg.Fields = t.Fields
	default:
		msg := fmt.Sprintf("'%s' is not a recognized view type. Using 'input' instead.", field.ViewType)
		log.Println(msg)
		tmpl, err = tmplFromWithDelims("gen-input.tmpl", [2]string{})
	}

	if err != nil {
		return err
	}

	err = tmpl.Execute(buf, templateArg)
	if err != nil {
		return err
	}

	field.View = buf.String()

	return nil
}

func optimizeFieldView(field *item.Field) {
	field.ViewType = strings.ToLower(field.ViewType)

	if field.IsReference {
		field.ViewType = "reference"
	} else if field.IsNested {
		field.ViewType = "nested"
	}

	// if we have a []T field type, automatically make the input view a repeater
	// as long as a repeater exists for the input type
	repeaterElements := []string{"input", "select", "file", "reference", "nested"}
	if strings.HasPrefix(field.TypeName, "[]") {
		for _, el := range repeaterElements {
			// if the viewType already is declared to be a -repeater
			// the comparison below will fail but the switch will
			// still find the right generator template
			// ex. authors:"[]string":select
			// ex. authors:string:select-repeater
			if field.ViewType == el {
				field.ViewType = field.ViewType + "-repeater"
			}
		}
	} else {
		// if the viewType is already declared as a -repeater, but
		// the TypeName is not of []T, add the [] prefix so the user
		// code is correct
		// ex. authors:string:select-repeater
		// ex. authors:@author:select-repeater
		if strings.HasSuffix(field.ViewType, "-repeater") {
			field.TypeName = "[]" + field.TypeName
		}
	}
}
