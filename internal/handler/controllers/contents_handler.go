package controllers

import (
	"net/http"

	"github.com/fanky5g/ponzu/content/editor"

	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/tokens"

	log "github.com/sirupsen/logrus"
)

func NewContentsHandler(r router.Router) http.HandlerFunc {
	contentService := r.Context().Service(tokens.ContentServiceToken).(content.Service)
	contentTypes := r.Context().Types().Content

	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")
		if t == "" {
			r.Renderer().BadRequest(res)
			return
		}

		if _, ok := contentTypes[t]; !ok {
			r.Renderer().BadRequest(res)
			return
		}

		pt := contentTypes[t]()
		if _, ok := pt.(editor.Editable); !ok {
			log.Warnf("item %s does not implement editable interface", t)
			r.Renderer().InternalServerError(res)
			return
		}

		searchRequestDto, err := request.GetSearchRequestDto(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get SearchRequestDto")
			r.Renderer().InternalServerError(res)
			return
		}

		search, err := request.MapSearchRequest(searchRequestDto)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request")
			r.Renderer().InternalServerError(res)
			return
		}

		renderContentList(r, res, t, search, pt, func() ([]interface{}, int, error) {
			return contentService.GetAllWithOptions(t, search)
		})
	}
}
