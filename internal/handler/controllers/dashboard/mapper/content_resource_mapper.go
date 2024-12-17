package mapper

import (
	"net/http"

	"github.com/fanky5g/ponzu/internal/handler/controllers/dashboard/resources"
)

func MapToContentIdentifier(r *http.Request) (*resources.ContentIdentifier, error) {
	q := r.URL.Query()

	return &resources.ContentIdentifier{
		ID:   q.Get("id"),
		Type: q.Get("type"),
	}, nil
}
