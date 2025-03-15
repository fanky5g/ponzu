package content

import (
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	"github.com/fanky5g/ponzu/internal/search"
	"net/http"

	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/entities"
	log "github.com/sirupsen/logrus"
)

func NewUploadContentsHandler(publicPath string, uploadService *UploadService, tmpl layouts.Template) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		pt := new(entities.Upload)
		_, ok := interface{}(pt).(editor.Editable)
		if !ok {
			log.Warning("entities.FileUpload is not editable")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		searchRequestDto, err := search.GetSearchRequestDto(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get search request dto")
			return
		}

		s, err := search.MapSearchRequest(searchRequestDto)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request dto")
			return
		}

		uploadResultLoader := func() ([]interface{}, int, error) {
			return uploadService.GetAllWithOptions(s)
		}

		data, err := buildTableViewModel(publicPath, pt, s, uploadResultLoader)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to build tabular view params")
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
