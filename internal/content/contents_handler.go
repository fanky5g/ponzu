package content

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	"github.com/fanky5g/ponzu/internal/search"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewContentsHandler(
	publicPath string,
	contentService *Service,
	contentTypes map[string]content.Builder,
	tmpl layouts.Template) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, ok := contentTypes[t]; !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		pt := contentTypes[t]()
		if _, ok := pt.(editor.Editable); !ok {
			log.Warnf("item %s does not implement editable interface", t)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		entity, ok := pt.(content.Entity)
		if !ok {
			log.Warnf("item %s does not implement entity interface", t)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		searchRequestDto, err := search.GetSearchRequestDto(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get SearchRequestDto")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		s, err := search.MapSearchRequest(searchRequestDto)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		contentLoader := func() ([]interface{}, int, error) {
			return contentService.GetAllWithOptions(t, s)
		}

		data, err := buildTableViewModel(publicPath, entity, s, contentLoader)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to build table params")
			res.WriteHeader(http.StatusInternalServerError)
			return

		}

		response.Respond(
			res,
			req,
			response.NewHTMLResponse(
				http.StatusOK,
				tmpl,
				data,
			),
		)
	}
}
