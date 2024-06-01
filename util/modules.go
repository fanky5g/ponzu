package util

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

var (
	ErrMustBeAGoModule = errors.New(`
Target application MUST be a go module. Please check that your working directory is a go module application
`)

	rModPath = regexp.MustCompile("^module\\s+(?P<Module>.*)$")
)

func GetModulePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	goModFile := fmt.Sprintf("%s/go.mod", wd)
	fileInfo, err := os.Stat(goModFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", ErrMustBeAGoModule
		}

		return "", errors.Wrap(err, "Failed to read target package module info.")
	}

	f, err := os.Open(fmt.Sprintf("%s/%s", wd, fileInfo.Name()))
	if err != nil {
		return "", errors.Wrap(err, "Failed to read target package module info.")
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.WithField("Error", err).Warning("Failed to close module file")
		}
	}()

	c, err := io.ReadAll(f)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read target package module info.")
	}

	modFileContent := strings.TrimSpace(string(c))
	modPathLine, _, ok := strings.Cut(modFileContent, "\n")
	if !ok {
		return "", errors.Wrap(err, "Failed to read target package module info.")
	}

	matches := rModPath.FindStringSubmatch(modPathLine)
	if len(matches) == 0 {
		return "", ErrMustBeAGoModule
	}

	index := rModPath.SubexpIndex("Module")
	if index == -1 {
		return "", ErrMustBeAGoModule
	}

	return matches[index], nil
}
