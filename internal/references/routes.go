package references

import (
	"net/http"

	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/uploads"
	"github.com/fanky5g/ponzu/tokens"
)

func RegisterRoutes(r router.Router) {
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)
	uploadService := r.Context().Service(tokens.UploadServiceToken).(*uploads.Service)

	service, err := New(contentService, uploadService)
	if err != nil {
		panic(err)
	}

	mapper := NewMapper()

	r.APIAuthorizedRoute("GET /api/references", func(r router.Router) http.HandlerFunc {
		return NewListReferencesHandler(service, mapper)
	})

	r.APIAuthorizedRoute("GET /api/references/{id}", func(r router.Router) http.HandlerFunc {
		return NewGetReferenceHandler(service, mapper)
	})
}
