package localstorage

import (
	"fmt"
	"github.com/fanky5g/ponzu/config"
	"net/http"
	"path/filepath"
)

type Client struct {
	storageDir string
	fs         http.FileSystem
}

func New(dir string) (*Client, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %v", err)
	}

	if dir == "" {
		dir = filepath.Join(cfg.Paths.DataDir, "uploads")
	}

	return &Client{storageDir: dir, fs: justFilesFilesystem{http.Dir(dir)}}, nil
}
