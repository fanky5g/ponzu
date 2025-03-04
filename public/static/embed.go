//go:build embed

package static

import "embed"

//go:embed all:*
var AssetFS embed.FS
