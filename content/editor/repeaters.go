package editor

import (
	"fmt"
	"regexp"
	"strings"
)

var parentIsFieldCollectionRegexp = regexp.MustCompile("(?P<Position>\\d).Value$")

func makePositionalPlaceholder(tagName string) string {
	return fmt.Sprintf("%%%s%%", strings.NewReplacer(".", "_").Replace(tagName))
}
