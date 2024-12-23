package response

import (
	"net/http"
	"net/url"
)

type redirect struct {
	publicPath string
	target     string
}

func (r *redirect) Render(w http.ResponseWriter, req *http.Request) error {
	u, err := url.Parse(r.target)
	if err != nil {
		return err
	}

	p, err := url.JoinPath(r.publicPath, u.Path)
	if err != nil {
		return err
	}

	location := req.URL.Scheme + req.URL.Host + p
	if u.RawQuery != "" {
		location += "?" + u.RawQuery
	}

	http.Redirect(w, req, location, http.StatusFound)
	return nil
}

func NewRedirectResponse(publicPath, target string) *Response {
	return &Response{
		Renderer: &redirect{publicPath: publicPath, target: target},
	}
}
