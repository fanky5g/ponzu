package editor

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content/item"
	"strings"
)

type Tab struct {
	Name       string
	Icon       string
	ClassNames []string
	Content    []byte
}

func getTabIdentifier(name string) string {
	return strings.ToLower(name)
}

func getTabContent(tab Tab) string {
	return fmt.Sprintf(
		`<div id="%s" class="col s12 card-content %s">%s</div>`,
		getTabIdentifier(tab.Name),
		strings.Join(tab.ClassNames, " "),
		string(tab.Content),
	)
}

// Form takes editable entities and any number of Field funcs to describe the edit
// page for any entities struct added by a user
func Form(post Editable, paths config.Paths, fields ...Field) ([]byte, error) {
	viewBuf := &bytes.Buffer{}

	tabs, err := getTabs(post, fields)
	if err != nil {
		return nil, err
	}

	if len(tabs) == 0 {
		return nil, nil
	}

	if len(tabs) == 1 {
		tab := tabs[0]
		_, err = viewBuf.WriteString(fmt.Sprintf(`<div class="row">%s`, getTabContent(tab)))
		if err != nil {
			return nil, fmt.Errorf("failed to write HTML string to editor Form buffer: %v", err)
		}
	} else {
		_, err = viewBuf.WriteString(`<div class="row"><div class="col s12"><ul class="tabs">`)
		if err != nil {
			return nil, fmt.Errorf("failed to write HTML string to editor Form buffer: %v", err)
		}

		// write tab header
		for _, tab := range tabs {
			_, err = viewBuf.WriteString(fmt.Sprintf(`
				<li class="tab col s3"><a href="#%s"><i class="material-icons">%s</i>%s</a></li>`,
				getTabIdentifier(tab.Name),
				tab.Icon,
				tab.Name,
			))
			if err != nil {
				return nil, fmt.Errorf("failed to write HTML string to editor Form buffer: %v", err)
			}
		}

		// tab header closing tag
		_, err = viewBuf.WriteString(`</ul></div>`)
		if err != nil {
			return nil, fmt.Errorf("failed to write HTML string to editor Form buffer: %v", err)
		}

		// write content
		for _, tab := range tabs {
			_, err = viewBuf.WriteString(getTabContent(tab))
			if err != nil {
				return nil, fmt.Errorf("failed to write HTML string to editor Form buffer: %v", err)
			}
		}
	}

	script := &bytes.Buffer{}
	scriptTmpl := makeScript("editor")
	if err = scriptTmpl.Execute(script, paths); err != nil {
		panic(err)
	}

	editorControls := makeHtml("editor_controls")
	_, err = viewBuf.WriteString(editorControls + script.String() + `</div>`)
	if err != nil {
		return nil, fmt.Errorf("failed to write HTML string to editor Form buffer: %v", err)
	}

	return viewBuf.Bytes(), nil
}

func getTabs(e Editable, contentFields []Field) ([]Tab, error) {
	tabs := make([]Tab, 0)
	if len(contentFields) > 0 {
		contentTab, err := getContentTab(contentFields)
		if err != nil {
			return nil, err
		}

		if contentTab != nil {
			tabs = append(tabs, *contentTab)
		}
	}

	properties, err := getPropertiesTab(e)
	if err != nil {
		return nil, err
	}

	if properties != nil {
		tabs = append(tabs, *properties)
	}

	return tabs, nil
}

func getContentTab(fields []Field) (*Tab, error) {
	viewBuf := &bytes.Buffer{}
	for _, f := range fields {
		_, err := viewBuf.Write(f.View)
		if err != nil {
			return nil, fmt.Errorf("error writing field view to editor view buffer: %v", err)
		}
	}

	return &Tab{
		Name:       "Edit",
		Icon:       "edit",
		Content:    viewBuf.Bytes(),
		ClassNames: []string{"editor-content"},
	}, nil
}

func getPropertiesTab(e Editable) (*Tab, error) {
	properties := getDefaultFields(e)
	if len(properties) == 0 {
		return nil, nil
	}

	viewBuf := &bytes.Buffer{}
	for _, f := range properties {
		_, err := viewBuf.Write(f.View)
		if err != nil {
			return nil, fmt.Errorf("error writing field view to editor view buffer: %v", err)
		}
	}

	_, err := viewBuf.WriteString(makeHtml("editor_timestamp"))
	if err != nil {
		return nil, fmt.Errorf("error writing field view to editor view buffer: %v", err)
	}

	return &Tab{
		Name:       "Properties",
		Icon:       "tune",
		Content:    viewBuf.Bytes(),
		ClassNames: []string{"editor-metadata"},
	}, nil
}

// Default fields (properties) are system generated, and mostly non-editable. Most system entities that do not
// have properties by default. E.g. System Configuration entities. As such, we only render properties
// for existing entities that already have these properties. This
// allows us to omit rendering properties unnecessarily for auto-generated system entities.
func getDefaultFields(e Editable) []Field {
	iface := (interface{})(e)
	properties := make([]Field, 0)
	if sluggable, ok := iface.(item.Sluggable); ok {
		if sluggable.ItemSlug() != "" {
			properties = append(properties, Field{
				View: Input("Slug", e, map[string]string{
					"label":       "URL Slug",
					"type":        "text",
					"disabled":    "true",
					"placeholder": "Will be set automatically",
				}, nil),
			})
		}
	}

	if sortable, ok := iface.(item.Sortable); ok {
		if sortable.Time() != 0 {
			properties = append(properties, []Field{
				{
					View: Timestamp("Timestamp", e, map[string]string{
						"type":  "hidden",
						"class": "__ponzu timestamp",
					}),
				},
				{
					View: Timestamp("Updated", e, map[string]string{
						"type":  "hidden",
						"class": "__ponzu updated",
					}),
				},
			}...)
		}
	}

	return properties
}
