package editor

import (
	"io"

	"github.com/fanky5g/ponzu/internal/views"
)

type ReferenceSelectDataProvider struct {
	ContentType            string
	PublicPath             string
	OptionTemplate         string
	SelectedOptionTemplate string
	SelectType
}

func (provider *ReferenceSelectDataProvider) RenderClientOptionsProvider(w io.Writer, selector string) error {
	return views.ExecuteTemplate(w, "reference_options_loader.gohtml", struct {
		PublicPath             string
		ContentType            string
		Selector               string
		OptionTemplate         string
		SelectedOptionTemplate string
		SelectType
	}{
		PublicPath:             provider.PublicPath,
		ContentType:            provider.ContentType,
		Selector:               selector,
		OptionTemplate:         provider.OptionTemplate,
		SelectedOptionTemplate: provider.SelectedOptionTemplate,
		SelectType:             provider.SelectType,
	})
}

// ReferenceSelect returns the []byte of a <select> HTML element plus internal <options> with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func ReferenceSelect(
	publicPath,
	fieldName string,
	p interface{},
	attrs map[string]string,
	args *FieldArgs,
	contentType string,
) []byte {
	return SelectWithDataProvider(fieldName, p, attrs, args, &ReferenceSelectDataProvider{
		ContentType:    contentType,
		PublicPath:     publicPath,
		OptionTemplate: SelectOptionTemplate,
		SelectType:     SingleSelect,
	})
}

// ReferenceSelectRepeater returns the []byte of a <select> HTML element plus internal <options> with a label.
// It also includes repeat controllers (+ / -) so the element can be
// dynamically multiplied or reduced.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func ReferenceSelectRepeater(
	publicPath,
	fieldName string,
	p interface{},
	attrs map[string]string,
	args *FieldArgs,
	contentType string,
) []byte {
	return MultiSelectWithDataProvider(fieldName, p, attrs, &ReferenceSelectDataProvider{
		ContentType:            contentType,
		PublicPath:             publicPath,
		OptionTemplate:         SelectOptionTemplate,
		SelectedOptionTemplate: SelectedOptionTemplate,
		SelectType:             MultipleSelect,
	})
}
