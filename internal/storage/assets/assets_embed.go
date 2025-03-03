//go:build embed

package assets

import (
	"github.com/fanky5g/ponzu/public/static"
	"net/http"
)

func init() {
	AssetStorage = http.FS(static.AssetFS)
}
