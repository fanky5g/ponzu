package content

import (
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/http/response"
	"log"
	"net/http"
)

func NewEditContentFormHandler(contentService *Service, cfg config.ConfigCache, publicPath string) http.HandlerFunc {
	tmpl, err := getEditPageTemplate()
	if err != nil {
		log.Fatalf("Failed to build page template: %v", err)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		contentQuery, err := MapContentQueryFromRequest(req)
		if err != nil {
			// TODO: handle error
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		entity, err := contentService.Type(contentQuery.Type)
		if err != nil {
			// TODO: handle error
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if contentQuery.ID != "" {
			entity, err = contentService.GetContent(contentQuery.Type, contentQuery.ID)
			if err != nil {
				// TODO: handle error
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		editContentForm, err := NewEditContentFormViewModel(
			entity,
			cfg,
			publicPath,
			contentService.ContentTypes(),
		)
		if err != nil {
			// TODO: handle error
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Write(
			res,
			response.NewHTMLResponse(http.StatusOK, tmpl, editContentForm),
		)
	}
}
