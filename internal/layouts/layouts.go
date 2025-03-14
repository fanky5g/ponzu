package layouts

import (
	"github.com/fanky5g/ponzu/internal/layouts/layout"
	"io"
)

type Template interface {
	Child(names ...string) (*layout.Layout, error)
	Execute(w io.Writer, data interface{}) error
}
