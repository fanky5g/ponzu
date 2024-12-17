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
	tabIdentifier := getTabIdentifier(tab.Name)

	return fmt.Sprintf(
		`<div id="%s" role="tab-panel" aria-labelledby="%s" class="%s">%s</div>`,
		tabIdentifier,
		tabIdentifier,
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

	viewBuf.WriteString(`<div class="form-view-content">`)

	if len(tabs) > 1 {
		_, err = viewBuf.WriteString(`
        <div class="mdc-tab-bar" role="tablist">
            <div class="mdc-tab-scroller">
                <div class="mdc-tab-scroller__scroll-area">
                    <div class="mdc-tab-scroller__scroll-content">
        `)

		if err != nil {
			return nil, fmt.Errorf("failed to write HTML string to editor Form buffer: %v", err)
		}

		// write tab header
		for _, tab := range tabs {
			_, err = viewBuf.WriteString(fmt.Sprintf(`
        <button type="button" class="mdc-tab" role="tab" aria-selected="true" tabindex="0">
          <span class="mdc-tab__content">
            <span class="mdc-tab__icon material-icons" aria-hidden="true">%s</span>
            <span class="mdc-tab__text-label">%s</span>
          </span>
          <span class="mdc-tab-indicator">
            <span class="mdc-tab-indicator__content mdc-tab-indicator__content--underline"></span>
          </span>
          <span class="mdc-tab__ripple"></span>
        </button>`,
				tab.Icon,
				getTabIdentifier(tab.Name),
			))

			if err != nil {
				return nil, fmt.Errorf("failed to write HTML string to editor Form buffer: %v", err)
			}
		}

		// tab header closing tag
		_, err = viewBuf.WriteString(`</div></div></div></div>`)
		if err != nil {
			return nil, fmt.Errorf("failed to write HTML string to editor Form buffer: %v", err)
		}

	}

	// write content
	for _, tab := range tabs {
		_, err = viewBuf.WriteString(getTabContent(tab))
		if err != nil {
			return nil, fmt.Errorf("failed to write HTML string to editor Form buffer: %v", err)
		}
	}

	if _, err = viewBuf.WriteString(`</div></div>`); err != nil {
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
	iface := (interface{})(e)
	viewBuf := &bytes.Buffer{}
	if sluggable, ok := iface.(item.Sluggable); ok {
		if sluggable.ItemSlug() != "" {
			_, err := viewBuf.Write(Input("Slug", e, map[string]string{
				"label":       "URL Slug",
				"type":        "text",
				"disabled":    "true",
				"placeholder": "Will be set automatically",
			}, nil),
			)
			if err != nil {
				return nil, fmt.Errorf("error writing field view to editor view buffer: %v", err)
			}
		}
	}

	if viewBuf.Len() == 0 {
		return nil, nil
	}

	return &Tab{
		Name:       "Properties",
		Icon:       "tune",
		Content:    viewBuf.Bytes(),
		ClassNames: []string{"editor-content"},
	}, nil
}
