package config

import (
	"os"
	"path/filepath"
	"runtime"
)

var rootPath string

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootPath = filepath.Join(filepath.Dir(b), "..")
}

func TlsDir() string {
	tlsDir := os.Getenv("PONZU_TLS_DIR")
	if tlsDir == "" {
		tlsDir = filepath.Join(rootPath, "internal", "tls")
	}

	return tlsDir
}
