package generate

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go/format"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func writeTemplateFile(fileName, templateName string, data interface{}) error {
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: go.mod. Error: %v", err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
				"File":  f.Name(),
			}).Warning("Failed to close file")
		}
	}()

	return writeTemplate(templateName, data, f)
}

func writeTemplate(name string, data interface{}, w io.Writer) error {
	templateRoot := fmt.Sprintf("%s/cmd/generate/templates", rootPath)
	templ := template.Must(template.New(name).ParseFiles(filepath.Join(templateRoot, name)))

	buf := &bytes.Buffer{}
	err := templ.Execute(buf, data)

	if err != nil {
		return err
	}

	if strings.HasSuffix(strings.TrimSuffix(name, ".tmpl"), ".go") {
		var fmtBuf []byte
		fmtBuf, err = format.Source(buf.Bytes())
		if err != nil {
			return fmt.Errorf("failed to format template: %s.\n\n%s", err.Error(), buf)
		}

		_, err = io.Copy(w, bytes.NewBuffer(fmtBuf))
	} else {
		_, err = io.Copy(w, buf)
	}

	return err
}

func copyFiles(src, target string) error {
	dirName := filepath.Clean(src)
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			targetPath := filepath.Join(target, strings.TrimPrefix(path, dirName))
			return copyFile(path, targetPath)
		}

		return nil
	})
}

func copyFile(src, target string) error {
	targetDir := filepath.Dir(target)

	_, err := os.Stat(targetDir)
	if errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	var srcFile *os.File
	srcFile, err = os.Open(src)
	if err != nil {
		return err
	}

	defer func() {
		if err = srcFile.Close(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
				"Path":  src,
			}).Warning("Failed to close src file")
		}
	}()

	var targetFile *os.File
	targetFile, err = os.Create(target)
	if err != nil {
		return err
	}

	defer func() {
		if err = targetFile.Close(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
				"Path":  target,
			}).Warning("Failed to close target file")
		}
	}()

	_, err = io.Copy(targetFile, srcFile)
	return err
}
