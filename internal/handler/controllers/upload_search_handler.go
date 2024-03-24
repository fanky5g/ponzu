package controllers

import (
	"bytes"
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/search"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewUploadSearchHandler(r router.Router) http.HandlerFunc {
	searchService := r.Context().Service(tokens.UploadSearchServiceToken).(search.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		searchRequest, err := request.GetSearchRequestDto(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get search request dto")
			r.Renderer().InternalServerError(res)
			return
		}

		matches, err := searchService.Search(constants.UploadsEntityName, searchRequest.Query, searchRequest.Count, searchRequest.Offset)
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		b := &bytes.Buffer{}

		html := `<div class="col s9 card">
					<div class="card-content">
					<div class="row">
					<div class="card-title col s7">Uploads Results</div>
					<form class="col s4" action="{{ .PublicPath }}/uploads/search" method="get">
						<div class="input-field post-search inline">
							<label class="active">Search:</label>
							<i class="right material-icons search-icon">search</i>
							<input class="search" name="q" type="text" placeholder="Within all upload fields" class="search"/>
							<input type="hidden" name="type" value="` + constants.UploadsEntityName + `" />
						</div>
                   </form>
					</div>
					<ul class="posts row">`

		for i := range matches {
			contentEntryTemplate := editor.BuildContentListEntryTemplate(matches[i].(editor.Editable), constants.UploadsEntityName)
			_, err = b.Write([]byte(contentEntryTemplate))
			if err != nil {
				log.WithField("Error", err).Warning("Failed to write template")
				r.Renderer().InternalServerError(res)
				return
			}
		}

		_, err = b.WriteString(`</ul></div></div>`)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to write template")
			r.Renderer().InternalServerError(res)
			return
		}

		btn := `<div class="col s3"><a href="{{ .PublicPath }}/edit/upload" class="btn new-post waves-effect waves-light">New upload</a></div></div>`
		html = html + b.String() + btn

		r.Renderer().InjectTemplateInAdmin(res, html, nil)
	}
}
