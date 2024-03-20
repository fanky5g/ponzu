package util

import (
	"fmt"
	conf "github.com/fanky5g/ponzu/config"
	"net/http"
	"strings"
)

func Redirect(
	req *http.Request,
	res http.ResponseWriter,
	paths conf.Paths,
	location string,
	status int,
) {
	target := fmt.Sprintf("%s/%s", paths.PublicPath, strings.TrimPrefix(location, "/"))
	redir := req.URL.Scheme + req.URL.Host + target

	http.Redirect(res, req, redir, status)
}
