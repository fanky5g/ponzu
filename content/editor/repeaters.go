package editor

import (
	"fmt"
	"regexp"
	"strings"
)

// rModPath = regexp.MustCompile("^module\\s+(?P<Module>.*)$")
var parentIsFieldCollectionRegexp = regexp.MustCompile("(?P<Position>\\d).Value$")

func makePositionalPlaceholder(tagName string) string {
	return fmt.Sprintf("%%%s%%", strings.NewReplacer(".", "_").Replace(tagName))
}
