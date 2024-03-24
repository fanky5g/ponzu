package util

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func Html(templates ...string) string {
	var tmpl string
	for _, name := range templates {
		htmlString, err := getTemplateString(
			fmt.Sprintf("%s%s", strings.TrimSuffix(name, ".gohtml"), ".gohtml"),
		)
		if err != nil {
			panic(err)
		}

		tmpl += htmlString
	}

	return tmpl
}

func getTemplateString(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open template file: %s. Error: %v", path, err)
	}

	htmlBytes, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	return string(htmlBytes), err
}
