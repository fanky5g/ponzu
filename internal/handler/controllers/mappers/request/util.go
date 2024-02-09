package request

import (
	"net/http"
	"strings"
)

func getContentType(req *http.Request) string {
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, ";") {
		contentType = strings.TrimSpace(contentType[:strings.Index(contentType, ";")])
	}

	return contentType
}
