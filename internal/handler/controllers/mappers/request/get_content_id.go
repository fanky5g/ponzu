package request

import (
	"net/http"
	"regexp"
	"strconv"
)

var (
	contentIdentifierRegex = regexp.MustCompile("/api/entities/(?P<identifier>[^/]+)")
)

func GetRequestContentId(req *http.Request) (bool, string) {
	if matches := contentIdentifierRegex.FindStringSubmatch(req.URL.Path); len(matches) > 0 {
		if index := contentIdentifierRegex.SubexpIndex("identifier"); index != -1 && index < len(matches) {
			identifier := matches[index]
			isSlug, _ := strconv.ParseBool(req.URL.Query().Get("slug"))
			return isSlug, identifier
		}
	}

	return false, ""
}
