// Package manager contains the controllers UI to the CMS which wraps all entities editor
// interfaces to manage the create/edit/delete capabilities of Ponzu entities.
package manager

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
	"html/template"

	"github.com/gofrs/uuid"
)

const managerHTML = `
<form class="form-view-root" method="post" action="{{ .PublicPath }}/edit" enctype="multipart/form-data">
	<input type="hidden" name="uuid" value="{{.UUID}}"/>
	<input type="hidden" name="id" value="{{.ID}}"/>
	<input type="hidden" name="type" value="{{.Kind}}"/>
	<input type="hidden" name="slug" value="{{.Slug}}"/>
	{{ .Editor }}
</form>
`

var managerTmpl = template.Must(template.New("manager").Parse(managerHTML))

type manager struct {
	ID         string
	UUID       uuid.UUID
	Kind       string
	Slug       string
	Editor     template.HTML
	PublicPath string
}

func Manage(e editor.Editable, paths config.Paths, typeName string) ([]byte, error) {
	v, err := e.MarshalEditor(paths)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal editor for entities %s. %s", typeName, err.Error())
	}

	i, ok := e.(item.Identifiable)
	if !ok {
		return nil, fmt.Errorf("entities type %s does not implement item.Identifiable", typeName)
	}

	s, ok := e.(item.Sluggable)
	if !ok {
		return nil, fmt.Errorf("entities type %s does not implement item.Sluggable", typeName)
	}

	m := manager{
		ID:         i.ItemID(),
		Kind:       typeName,
		Slug:       s.ItemSlug(),
		Editor:     template.HTML(v),
		PublicPath: paths.PublicPath,
	}

	// execute templates template into buffer for func return val
	buf := &bytes.Buffer{}
	if err = managerTmpl.Execute(buf, m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
