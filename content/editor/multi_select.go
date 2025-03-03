package editor

import (
	"bytes"
	"github.com/fanky5g/ponzu/internal/templates"
	log "github.com/sirupsen/logrus"
)

func MultiSelectWithDataProvider(fieldName string, p interface{}, attrs map[string]string, dataProvider interface{}) []byte {
	selector := TagNameFromStructField(fieldName, p, nil)
	fieldVal := ValueFromStructField(fieldName, p, nil)

	var selected []string
	if fieldVal != nil {
		var ok bool
		selected, ok = fieldVal.([]string)
		if !ok {
			log.Warnf("Expected field value to be []string. Got %T", fieldVal)
			return nil
		}
	}

	templateBuffer := &bytes.Buffer{}
	var options []SelectOption
	if dataProvider != nil {
		switch dataProvider.(type) {
		case SelectInitialOptionsProvider:
			var err error
			options, err = dataProvider.(SelectInitialOptionsProvider).GetInitialOptions()
			if err != nil {
				log.Fatalf("Failed to get options for %s: %v", fieldName, err)
			}
		case SelectClientOptionsProvider:
			clientDataOptionsProvider := dataProvider.(SelectClientOptionsProvider)
			err := clientDataOptionsProvider.RenderClientOptionsProvider(templateBuffer, selector)
			if err != nil {
				log.Fatalf("Failed to render client options provider: %v", err)
			}
		default:
			log.Fatalf("Unsupported Select Options provider: %T", dataProvider)
			return nil
		}
	}

	sel := MultiSelectData{
		SelectData: SelectData{
			Label:       fieldName,
			Placeholder: attrs["label"],
			Selector:    selector,
			Name:        fieldName,
			Options:     options,
		},
		Selected: selected,
	}

	if err := templates.ExecuteTemplate(templateBuffer, "views/select/multi_select.gohtml", sel); err != nil {
		log.Fatalf("Failed to render multi-select: %v", err)
	}

	return templateBuffer.Bytes()
}
