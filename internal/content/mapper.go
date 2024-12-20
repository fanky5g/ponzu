package content

import (
	"net/http"
)

func MapToContentQuery(r *http.Request) (*ContentQuery, error) {
	q := r.URL.Query()

	return &ContentQuery{
		ID:   q.Get("id"),
		Type: q.Get("type"),
	}, nil
}
