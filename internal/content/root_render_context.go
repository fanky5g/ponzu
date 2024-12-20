package content

import "github.com/fanky5g/ponzu/content"

type RootRenderContext struct {
	PublicPath string
	AppName    string
	Logo       string
	Types      map[string]content.Builder
}
