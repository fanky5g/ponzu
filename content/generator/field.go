package generator

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/generator"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"text/template"
)

type Field struct {
	Name              string
	Label             string
	Initial           string
	TypeName          string
	JSONName          string
	ViewType          string
	IsReference       bool
	IsNested          bool
	IsFieldCollection bool
	ReferenceName     string
	ReferenceJSONTags []string

	// Render Scope data
	Parent   *Field
	Scope    *ViewScope
	Children []Field
}

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

func (field *Field) Validate() error {
	for jsonName, fieldName := range reservedFieldNames {
		if field.JSONName == jsonName || field.Name == fieldName {
			return fmt.Errorf("reserved field name: %s (%s)", jsonName, fieldName)
		}
	}

	return nil
}

func (field *Field) View() (string, error) {
	// Initialize variables
	var tmpl *template.Template
	var err error
	buf := &bytes.Buffer{}

	// Define switch cases for different view types
	switch field.ViewType {
	case "checkbox", "custom", "file", "hidden", "input", "text", "richtext", "select", "textarea",
		"tags", "input-repeater", "select-repeater", "file-repeater", "nested", "nested-repeater", "field-collection":
		tmpl, err = tmplFromWithDelims(field.Scope.TemplatesDir, field.ViewType+".tmpl", [2]string{})
	case "reference", "reference-repeater":
		tmpl, err = tmplFromWithDelims(field.Scope.TemplatesDir, field.ViewType+".tmpl", [2]string{"[[", "]]"})
		if err != nil {
			return "", err
		}
	default:
		msg := fmt.Sprintf("'%s' is not a recognized view type. Using 'input' instead.", field.ViewType)
		log.Println(msg)
		tmpl, err = tmplFromWithDelims(field.Scope.TemplatesDir, "input.tmpl", [2]string{})
	}

	if err != nil {
		return "", err
	}

	// Execute template
	err = tmpl.Execute(buf, field)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func mapBlockToField(contentTypes content.Types, block generator.Block) *Field {
	typeName := block.TypeName
	isNested := false
	isFieldCollection := false

	viewType := "input"
	data := strings.Split(block.Definition.Type, ":")
	if len(data) == 2 {
		viewType = data[1]
	}

	if block.Definition.IsReference {
		if _, ok := contentTypes.Content[block.ReferenceName]; ok {
			viewType = "reference"
		} else if _, ok = contentTypes.FieldCollections[block.ReferenceName]; ok {
			isFieldCollection = true
			viewType = "field-collection"
			typeName = "*" + typeName
		} else {
			isNested = true
			viewType = "nested"
		}

		if block.Definition.IsArray && !strings.HasPrefix(typeName, "[]") {
			typeName = fmt.Sprintf("[]%s", typeName)
		}
	}

	// if we have a []T field type, automatically make the input view a repeater
	// as long as a repeater exists for the input type
	repeaterElements := []string{"input", "select", "file", "reference", "nested"}
	if strings.HasPrefix(typeName, "[]") {
		for _, el := range repeaterElements {
			// if the viewType already is declared to be a -repeater
			// the comparison below will fail but the switch will
			// still find the right generator template
			// ex. authors:"[]string":select
			// ex. authors:string:select-repeater
			if viewType == el {
				viewType = viewType + "-repeater"
			}
		}
	} else {
		// if the viewType is already declared as a -repeater, but
		// the TypeName is not of []T, add the [] prefix so the user
		// code is correct
		// ex. authors:string:select-repeater
		// ex. authors:@author:select-repeater
		if strings.HasSuffix(viewType, "-repeater") {
			typeName = "[]" + typeName
		}
	}

	return &Field{
		Name:              block.Name,
		Label:             block.Label,
		Initial:           strings.ToLower(string(block.Name[0])),
		TypeName:          typeName,
		JSONName:          block.JSONName,
		ViewType:          viewType,
		IsReference:       block.Definition.IsReference,
		IsNested:          isNested,
		IsFieldCollection: isFieldCollection,
		ReferenceName:     block.ReferenceName,
		ReferenceJSONTags: block.ReferenceJSONTags,
	}
}

func getRootMethodReceiver(field *Field, callDepth int) string {
	if field.Parent != nil {
		return getRootMethodReceiver(field.Parent, callDepth+1)
	}

	if field.IsFieldCollection && callDepth > 0 {
		// Field Collection Editor render currently works with a hardcoded receiver v which is passed during rendering
		return "v"
	}

	return field.Scope.Definition.MethodReceiverName
}

func GetRootMethodReceiver(field *Field) string {
	return getRootMethodReceiver(field, 0)
}

func GetPath(field *Field) string {
	if field.Parent != nil && !field.Parent.IsFieldCollection {
		return strings.Join([]string{GetPath(field.Parent), field.Name}, ".")
	}

	return field.Name
}

// GetFieldArgVar is currently only used with field collections. It returns one of two values nil or args.
// If the root parent is a FieldCollection it returns args which is filled during template rendering
func GetFieldArgVar(field *Field) string {
	if field.Parent != nil {
		return GetFieldArgVar(field.Parent)
	}

	if field != nil && field.IsFieldCollection {
		return "args"
	}

	return "nil"
}

func GetCollectionTypes(field *Field) (map[string]*ViewScope, error) {
	if fcType, ok := field.Scope.ContentTypes.FieldCollections[field.ReferenceName]; ok {
		var fieldCollection content.FieldCollections

		if fieldCollection, ok = fcType().(content.FieldCollections); ok {
			collectionTypes := make(map[string]*ViewScope)
			for typeName := range fieldCollection.AllowedTypes() {
				var definition generator.TypeDefinition
				definition, ok = field.Scope.ContentTypes.Definitions[typeName]
				if !ok {
					return nil, fmt.Errorf("type definition for %s not found", typeName)
				}

				collectionTypes[typeName] = newViewScope(
					&definition,
					field.Scope.ContentTypes,
					field.Scope.Target,
					field.Scope.TemplatesDir,
				)
			}

			return collectionTypes, nil
		}
	}

	return nil, fmt.Errorf("invalid field collection type: %s", field.ReferenceName)
}

var fieldFuncMaps = template.FuncMap{
	"RootMethodReceiver": GetRootMethodReceiver,
	"Path":               GetPath,
	"CollectionTypes":    GetCollectionTypes,
	"FieldArg":           GetFieldArgVar,
}

// Helper function to create a template with custom delimiters
func tmplFromWithDelims(templateDir, filename string, delim [2]string) (*template.Template, error) {
	if delim[0] == "" || delim[1] == "" {
		delim = [2]string{"{{", "}}"}
	}

	return template.
		New(filename).
		Funcs(fieldFuncMaps).
		Delims(delim[0], delim[1]).
		ParseFiles(filepath.Join(templateDir, filename))
}
