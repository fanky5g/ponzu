package content

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/http/response"
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
			res.WriteHeader(http.StatusBadRequest)
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

		response.Respond(
			res,
			req,
			response.NewHTMLResponse(http.StatusOK, tmpl, editContentForm),
		)
	}
}

func NewSaveContentHandler(contentService *Service, publicPath string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		contentQuery, err := MapContentQueryFromRequest(req)

		entity, err := contentService.Type(contentQuery.Type)
		if err != nil {
			// TODO: handle error
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		entity, err = request.GetEntityFromFormData(entity, req.PostForm)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map request entity")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if contentQuery.ID == "" {
			contentQuery.ID, err = contentService.CreateContent(contentQuery.Type, entity)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to create content")
				// TODO: handle error
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else {
			_, err = contentService.UpdateContent(contentQuery.Type, contentQuery.ID, entity)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to update content")
				res.WriteHeader(http.StatusInternalServerError)
				// TODO: handle error
				return
			}
		}

		response.Respond(
			res,
			req,
			response.NewRedirectResponse(
				publicPath,
				req.URL.Path+"?type="+contentQuery.Type+"&id="+contentQuery.ID,
			),
		)
	}
}
