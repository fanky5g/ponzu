package generator

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	generatorTypes "github.com/fanky5g/ponzu/content/generator/types"
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

func (gt *generator) ValidateField(field *generatorTypes.Field) error {
	for jsonName, fieldName := range reservedFieldNames {
		if field.JSONName == jsonName || field.Name == fieldName {
			return fmt.Errorf("reserved field name: %s (%s)", jsonName, fieldName)
		}
	}

	return nil
}

// set the specified view inside the editor field for a generated field for a type
func (gt *generator) setFieldView(definition *generatorTypes.TypeDefinition, index int, args ...string) error {
	// Helper function to create a template with custom delimiters
	tmplFromWithDelims := func(filename string, delim [2]string) (*template.Template, error) {
		if delim[0] == "" || delim[1] == "" {
			delim = [2]string{"{{", "}}"}
		}
		return template.New(filename).Delims(delim[0], delim[1]).ParseFiles(filepath.Join(gt.templateDir, filename))
	}

	// Initialize variables
	var tmpl *template.Template
	var err error
	buf := &bytes.Buffer{}
	field := &definition.Fields[index]
	fieldGeneratorDataVar := "v"
	fieldGeneratorArgsVar := "nil"

	// Customize data and args variables if provided
	if len(args) > 0 {
		fieldGeneratorDataVar = args[0]
	}
	if len(args) > 1 {
		fieldGeneratorArgsVar = args[1]
	}

	// Define template arguments
	templateArg := struct {
		*generatorTypes.Field
		FieldGeneratorDataVar string
		FieldGeneratorArgsVar string
		Fields                []generatorTypes.Field
		Types                 map[string]interface{}
	}{
		Field:                 field,
		FieldGeneratorDataVar: fieldGeneratorDataVar,
		FieldGeneratorArgsVar: fieldGeneratorArgsVar,
	}

	optimizeFieldView(field)
	// Define switch cases for different view types
	switch field.ViewType {
	case "checkbox", "custom", "file", "hidden", "input", "text", "richtext", "select", "textarea", "tags":
		tmpl, err = tmplFromWithDelims("gen-"+field.ViewType+".tmpl", [2]string{})
	case "input-repeater", "select-repeater", "file-repeater":
		tmpl, err = tmplFromWithDelims("gen-"+field.ViewType+".tmpl", [2]string{})
	case "reference", "reference-repeater":
		delim := [2]string{"[[", "]]"}
		tmpl, err = tmplFromWithDelims("gen-"+field.ViewType+".tmpl", delim)
		if err != nil {
			return err
		}
	case "nested", "nested-repeater":
		tmpl, err = tmplFromWithDelims("gen-"+field.ViewType+".tmpl", [2]string{})
		t, ok := gt.types.Definitions[field.ReferenceName]
		if !ok {
			return fmt.Errorf("no definition matched for %s type", field.Name)
		}
		// Update field names and call setFieldView recursively
		for i := range t.Fields {
			t.Fields[i].Name = fmt.Sprintf("%s.%s", field.Name, t.Fields[i].Name)
			t.Fields[i].Initial = templateArg.FieldGeneratorDataVar
			if err := gt.setFieldView(&t, i, templateArg.FieldGeneratorDataVar, templateArg.FieldGeneratorArgsVar); err != nil {
				return err
			}
		}
		templateArg.Fields = t.Fields
	case "field-collection":
		tmpl, err = tmplFromWithDelims("field-collection.tmpl", [2]string{})
		f, ok := gt.types.FieldCollections[field.ReferenceName]
		if !ok {
			return fmt.Errorf("no definition matched for %s type", field.Name)
		}

		// Populate template arguments for field collection
		fieldCollection, ok := f().(content.FieldCollections)
		if !ok {
			return fmt.Errorf("type %s is not a valid FieldCollection", field.Name)
		}

		templateArg.Label = fieldCollection.Name()
		templateArg.Types = make(map[string]interface{})
		fieldGeneratorArgsVar = "args"
		templateArg.FieldGeneratorArgsVar = fieldGeneratorArgsVar
		for typeName := range fieldCollection.AllowedTypes() {
			t, ok := gt.types.Definitions[typeName]
			if !ok {
				return fmt.Errorf("type %s in FieldCollection is not valid", typeName)
			}

			t.Initial = field.Initial
			for i := range t.Fields {
				t.Fields[i].Initial = fieldGeneratorDataVar
				if err := gt.setFieldView(&t, i, fieldGeneratorDataVar, fieldGeneratorArgsVar); err != nil {
					return err
				}
			}

			templateArg.Types[typeName] = struct {
				generatorTypes.TypeDefinition
				FieldGeneratorDataVar string
				FieldGeneratorArgsVar string
			}{
				TypeDefinition:        t,
				FieldGeneratorDataVar: fieldGeneratorDataVar,
				FieldGeneratorArgsVar: templateArg.FieldGeneratorArgsVar,
			}
		}
	default:
		msg := fmt.Sprintf("'%s' is not a recognized view type. Using 'input' instead.", field.ViewType)
		log.Println(msg)
		tmpl, err = tmplFromWithDelims("gen-input.tmpl", [2]string{})
	}

	if err != nil {
		return err
	}

	// Execute template
	err = tmpl.Execute(buf, templateArg)
	if err != nil {
		return err
	}

	field.View = buf.String()
	return nil
}

func optimizeFieldView(field *generatorTypes.Field) {
	field.ViewType = strings.ToLower(field.ViewType)

	if field.IsReference {
		field.ViewType = "reference"
	} else if field.IsNested {
		field.ViewType = "nested"
	} else if field.IsFieldCollection {
		field.ViewType = "field-collection"
		field.TypeName = "*" + field.TypeName
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
