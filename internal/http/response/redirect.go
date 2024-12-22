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
	base := req.URL.Scheme + req.URL.Host
	location, err := url.JoinPath(base, r.publicPath, r.target)
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
