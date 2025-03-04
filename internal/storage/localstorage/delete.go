package localstorage

import (
	"os"
	"path/filepath"
	"strings"
)

func (c *Client) Delete(path string) error {
	// split and rebuild path in OS friendly way
	// use path to delete the physical file from disk
	pathSplit := strings.Split(strings.TrimPrefix(path, "/api/"), "/")
	return os.Remove(filepath.Join(pathSplit...))
}
