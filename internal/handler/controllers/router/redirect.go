package router

import (
	"fmt"
	"net/http"
	"strings"
)

func (r *router) Redirect(
	req *http.Request,
	res http.ResponseWriter,
	location string,
) {
	target := fmt.Sprintf("%s/%s", r.ctx.Paths().PublicPath, strings.TrimPrefix(location, "/"))
	redir := req.URL.Scheme + req.URL.Host + target

	http.Redirect(res, req, redir, http.StatusFound)
}
