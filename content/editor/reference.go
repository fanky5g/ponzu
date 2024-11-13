package editor

import (
	"io"

	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/views"
)

type ReferenceSelectDataProvider struct {
	ContentType string
	PublicPath  string
	Template    string
}

func (provider *ReferenceSelectDataProvider) RenderClientOptionsProvider(w io.Writer, selector string) error {
	return views.ExecuteTemplate(w, "reference_options_loader.gohtml", struct {
		PublicPath  string
		ContentType string
		Selector    string
		Template    string
	}{
		PublicPath:  provider.PublicPath,
		ContentType: provider.ContentType,
		Selector:    selector,
		Template:    provider.Template,
	})
}

// ReferenceSelect returns the []byte of a <select> HTML element plus internal <options> with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func ReferenceSelect(paths config.Paths, fieldName string, p interface{}, attrs map[string]string, contentType, tmplString string) []byte {
	return Select(fieldName, p, attrs, &ReferenceSelectDataProvider{
		ContentType: contentType,
		PublicPath:  paths.PublicPath,
		// TODO: remove tmplString definition from ReferenceSelect and Repeater
		// Support getting option template string definition from content type (i.e get type from p.FieldName)
		// and check if type supports an interface GetSelectOptionTemplate, else default to string below
		// Select component must also support Option templates
		Template: `
<li class="mdc-list-item" role="option" data-value="@>id">
    <span class="mdc-list-item__text">"@>name"</span>
</li>
		`,
	})
}

// ReferenceSelectRepeater returns the []byte of a <select> HTML element plus internal <options> with a label.
// It also includes repeat controllers (+ / -) so the element can be
// dynamically multiplied or reduced.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func ReferenceSelectRepeater(paths config.Paths, fieldName string, p interface{}, attrs map[string]string, contentType, tmplString string) []byte {
	return nil
}
