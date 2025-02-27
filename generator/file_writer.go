package generator

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type fileWriter struct{}

var FormattedSourceWriter Writer = new(fileWriter)

func (writer *fileWriter) Write(filePath string, buf []byte) error {
	w := buf
	var err error
	if strings.HasSuffix(strings.TrimSuffix(filePath, ".tmpl"), ".go") {
		w, err = format.Source(buf)
		if err != nil {
			return fmt.Errorf("failed to format template: %s.\n\n%s", err.Error(), buf)
		}
	}

	targetDir := filepath.Dir(filePath)
	if _, err = os.Stat(targetDir); errors.Is(err, fs.ErrNotExist) {
		if err = os.MkdirAll(targetDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create root directory: %v", err)
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}(file)

	_, err = file.Write(w)
	if err != nil {
		return fmt.Errorf("failed to write generated file buffer: %s", err.Error())
	}

	return nil
}
