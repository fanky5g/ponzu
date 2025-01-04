package editor

import (
	"bytes"
	"io"

	"github.com/fanky5g/ponzu/internal/views"
	log "github.com/sirupsen/logrus"
)

type SelectClientOptionsProvider interface {
	RenderClientOptionsProvider(w io.Writer, selector string) error
}

type SelectInitialOptionsProvider interface {
	GetInitialOptions() ([]string, error)
}

// Select returns the []byte of a <select> HTML element plus internal <options> with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Select(fieldName string, p interface{}, attrs, options map[string]string) []byte {
	return SelectWithDataProvider(fieldName, p, attrs, makeGenericSelectDataProvider(options))
}

func SelectWithDataProvider(fieldName string, p interface{}, attrs map[string]string, dataProvider interface{}) []byte {
	value := ""
	fieldVal := ValueFromStructField(fieldName, p, nil)
	var ok bool
	if value, ok = fieldVal.(string); !ok {
		log.Warnf("Expected field value to be string. Got %T", fieldVal)
	}

	var err error
	options := make([]string, 0)

	templateBuffer := &bytes.Buffer{}
	if dataProvider != nil {
		if initialdataOptionsProvider, ok := dataProvider.(SelectInitialOptionsProvider); ok {
			options, err = initialdataOptionsProvider.GetInitialOptions()
			if err != nil {
				log.Fatalf("Failed to get options for %s: %v", fieldName, err)
			}
		}

		if clientDataOptionsProvider, ok := dataProvider.(SelectClientOptionsProvider); ok {
			err = clientDataOptionsProvider.RenderClientOptionsProvider(templateBuffer, fieldName)
			if err != nil {
				log.Fatalf("Failed to render client options provider: %v", err)
			}
		}
	}

	values := make([]string, len(options))
	i := 0
	for _, v := range options {
		values[i] = v
		i = i + 1
	}

	sel := struct {
		Name        string
		Label       string
		Placeholder string
		Value       string
		Options     []string
		Selector    string
	}{
		Label:       fieldName,
		Placeholder: attrs["label"],
		Selector:    fieldName,
		Name:        fieldName,
		Options:     options,
		Value:       value,
	}

	if err := views.ExecuteTemplate(templateBuffer, "select.gohtml", sel); err != nil {
		log.Fatalf("Failed to render select: %v", err)
	}

	return templateBuffer.Bytes()
}

type selectInitialOptionsProvider struct {
	options map[string]string
}

func (s *selectInitialOptionsProvider) GetInitialOptions() ([]string, error) {
	options := make([]string, len(s.options))
	i := 0
	for k := range s.options {
		options[i] = k
		i = i + 1
	}

	return options, nil
}

func makeGenericSelectDataProvider(options map[string]string) SelectInitialOptionsProvider {
	return &selectInitialOptionsProvider{
		options: options,
	}
}
