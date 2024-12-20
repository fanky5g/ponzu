package storage

import (
	localStorage "github.com/fanky5g/ponzu-driver-local-storage"

	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/storage"
)

func NewAssetsStorage() (storage.Client, error) {
	return localStorage.New(config.AssetStaticDir())
}
