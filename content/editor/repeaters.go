package editor

import (
	"fmt"
	"regexp"
	"strings"
)

var ancestorIsFieldCollectionRegexp = regexp.MustCompile("(?P<FieldCollectionName>.*).(?P<Position>\\d).Value")

func makePositionalPlaceholder(tagName string) string {
	return fmt.Sprintf("%%%s%%", strings.NewReplacer(".", "_").Replace(tagName))
}
