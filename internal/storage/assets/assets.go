package assets

import (
	"github.com/fanky5g/ponzu/internal/storage/localstorage"
	"net/http"
	"path/filepath"
	"runtime"
)

var (
	AssetStorage http.FileSystem
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	assetDirectory := filepath.Join(filepath.Dir(b), "../../..", "public/static")

	var err error
	AssetStorage, err = localstorage.New(assetDirectory)
	if err != nil {
		panic(err)
	}
}
