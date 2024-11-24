package editor

import (
	"io"

	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/views"
)

var DefaultReferenceSelectTemplate = `
<li class="mdc-list-item" role="option" data-value="@>id">
    <span class="mdc-list-item__text">"@>name"</span>
</li>
`

type ReferenceSelectDataProvider struct {
	ContentType string
	PublicPath  string
	Template    string
}

type ReferenceSelectTemplateProvider interface {
	GetTemplate() string
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
func ReferenceSelect(
	paths config.Paths,
	fieldName string,
	p interface{},
	attrs map[string]string,
	contentType string,
) []byte {
	template := DefaultReferenceSelectTemplate
	field := GetStructFieldInterface(p, fieldName)
	if field != nil {
		if customRowTemplateProvider, ok := field.(ReferenceSelectTemplateProvider); ok {
			template = customRowTemplateProvider.GetTemplate()
		}
	}

	return SelectWithDataProvider(fieldName, p, attrs, &ReferenceSelectDataProvider{
		ContentType: contentType,
		PublicPath:  paths.PublicPath,
		Template:    template,
	})
}

// ReferenceSelectRepeater returns the []byte of a <select> HTML element plus internal <options> with a label.
// It also includes repeat controllers (+ / -) so the element can be
// dynamically multiplied or reduced.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func ReferenceSelectRepeater(
	paths config.Paths,
	fieldName string,
	p interface{},
	attrs map[string]string,
	contentType,
	tmplString string,
) []byte {
	return nil
}
