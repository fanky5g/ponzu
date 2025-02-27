package editor

import (
	"bytes"
	"io"
	"strings"

	"github.com/fanky5g/ponzu/internal/views"
	log "github.com/sirupsen/logrus"
)

type SelectType string

type SelectOption struct {
	Value string
	Label string
}

type SelectData struct {
	Name        string
	Label       string
	Placeholder string
	Value       string
	Options     []SelectOption
	Selector    string
}

type MultiSelectData struct {
	SelectData
	Selected []string
}

var (
	SingleSelect         SelectType = "single"
	MultipleSelect       SelectType = "multiple"
	SelectOptionTemplate            = `
	<li class="mdc-list-item" role="option" data-value="@>id">
		<span class="mdc-list-item__text">"@>name"</span>
	</li>
	`
	// SelectedOptionTemplate is used on the client to render selected entry as a chip.
	// Must be synced with chip_template in multi_select.gohtml
	SelectedOptionTemplate = `
    <div class="mdc-chip" role="row" data-value="@>id">
      <div class="mdc-chip__ripple"></div>
      <span role="gridcell">
        <span role="button" tabindex="0" class="mdc-chip__primary-action">
          <span class="mdc-chip__text">"@>name"</span>
        </span>
      </span>
      <span role="gridcell">
        <i class="material-icons mdc-chip__icon mdc-chip__icon--trailing" tabindex="-1" role="button">cancel</i>
      </span>
    </div>
  `
)

type SelectClientOptionsProvider interface {
	RenderClientOptionsProvider(w io.Writer, selector string) error
}

type SelectInitialOptionsProvider interface {
	GetInitialOptions() ([]SelectOption, error)
}

// Select returns the []byte of a <select> HTML element plus internal <options> with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Select(fieldName string, p interface{}, attrs map[string]string, options []string, args *FieldArgs) []byte {
	return SelectWithDataProvider(fieldName, p, attrs, args, makeGenericSelectDataProvider(options))
}

func SelectWithDataProvider(fieldName string, p interface{}, attrs map[string]string, args *FieldArgs, dataProvider interface{}) []byte {
	value := ""

	selector := TagNameFromStructField(fieldName, p, args)
	fieldVal := ValueFromStructField(fieldName, p, args)
	var ok bool
	if value, ok = fieldVal.(string); !ok {
		log.Warnf("Expected field value to be string. Got %T", fieldVal)
	}

	var err error
	options := make([]SelectOption, 0)

	templateBuffer := &bytes.Buffer{}
	if dataProvider != nil {
		switch dataProvider.(type) {
		case SelectInitialOptionsProvider:
			options, err = dataProvider.(SelectInitialOptionsProvider).GetInitialOptions()
			if err != nil {
				log.Fatalf("Failed to get options for %s: %v", fieldName, err)
			}
		case SelectClientOptionsProvider:
			clientDataOptionsProvider := dataProvider.(SelectClientOptionsProvider)
			err = clientDataOptionsProvider.RenderClientOptionsProvider(templateBuffer, selector)
			if err != nil {
				log.Fatalf("Failed to render client options provider: %v", err)
			}
		default:
			log.Fatalf("Unsupported Select Options provider: %T", dataProvider)
			return nil
		}
	}

	sel := SelectData{
		Label:       fieldName,
		Placeholder: attrs["label"],
		Selector:    selector,
		Name:        fieldName,
		Options:     options,
		Value:       value,
	}

	if err = views.ExecuteTemplate(templateBuffer, "select.gohtml", sel); err != nil {
		log.Fatalf("Failed to render select: %v", err)
	}

	return templateBuffer.Bytes()
}

type selectInitialOptionsProvider struct {
	options []string
}

func (s *selectInitialOptionsProvider) GetInitialOptions() ([]SelectOption, error) {
	options := make([]SelectOption, len(s.options))
	for i := range s.options {
		v := strings.Split(s.options[i], ":")
		selectOption := SelectOption{
			Value: v[0],
			Label: v[0],
		}

		if len(v) > 1 {
			selectOption.Label = v[1]
		}

		options[i] = selectOption
	}

	return options, nil
}

func makeGenericSelectDataProvider(options []string) SelectInitialOptionsProvider {
	return &selectInitialOptionsProvider{
		options: options,
	}
}
