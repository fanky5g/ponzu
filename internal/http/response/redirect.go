package response

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type redirect struct {
	publicPath string
	target     string
}

func (r *redirect) Render(w http.ResponseWriter, req *http.Request) error {
	location, err := getRedirectLocation(r.publicPath, r.target)
	if err != nil {
		return err
	}

	http.Redirect(w, req, location, http.StatusFound)
	return nil
}

func NewRedirectResponse(publicPath, target string) *Response {
	return &Response{
		Renderer: &redirect{publicPath: publicPath, target: target},
	}
}

func getRedirectLocation(
	publicPath,
	target string,
) (string, error) {
	if !strings.HasPrefix(publicPath, "/") {
		publicPath = fmt.Sprintf("/%s", publicPath)
	}

	if !strings.HasPrefix(target, "/") {
		target = fmt.Sprintf("/%s", target)
	}

	u, err := url.Parse(target)
	if err != nil {
		return "", err
	}

	location, err := url.JoinPath(publicPath, u.Path)
	if err != nil {
		return "", err
	}

	if u.RawQuery != "" {
		location += "?" + u.RawQuery
	}

	return location, nil
}
