package controllers

import (
	"bytes"
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/domain/services/management/editor"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/views"
	"github.com/fanky5g/ponzu/internal/services/config"
	"github.com/fanky5g/ponzu/internal/services/search"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/internal/util"
	"log"
	"net/http"
)

func NewUploadSearchHandler(
	pathConf conf.Paths,
	configService config.Service,
	searchService search.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		status := q.Get("status")

		appName, err := configService.GetAppName()
		if err != nil {
			log.Printf("Failed to get app name: %v\n", appName)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		searchRequest, err := request.GetSearchRequestDto(req)
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusBadRequest)
			errView, err := views.Admin(util.Html("error_400"), appName, pathConf)
			if err != nil {
				return
			}

			res.Write(errView)
			return
		}

		matches, err := searchService.Search(storage.UploadsEntityName, searchRequest.Query, searchRequest.Count, searchRequest.Offset)
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
					<form class="col s4" action="` + pathConf.PublicPath + `/uploads/search" method="get">
						<div class="input-field post-search inline">
							<label class="active">Search:</label>
							<i class="right material-icons search-icon">search</i>
							<input class="search" name="q" type="text" placeholder="Within all upload fields" class="search"/>
							<input type="hidden" name="type" value="` + storage.UploadsEntityName + `" />
						</div>
                   </form>
					</div>
					<ul class="posts row">`

		for i := range matches {
			post := PostListItem(matches[i].(editor.Editable), storage.UploadsEntityName, status, pathConf)
			_, err = b.Write(post)
			if err != nil {
				log.Println(err)

				res.WriteHeader(http.StatusInternalServerError)
				errView, err := views.Admin(util.Html("error_500"), appName, pathConf)
				if err != nil {
					log.Println(err)
				}

				res.Write(errView)
				return
			}
		}

		_, err = b.WriteString(`</ul></div></div>`)
		if err != nil {
			LogAndFail(res, err, appName, pathConf)
			return
		}

		btn := `<div class="col s3"><a href="` + pathConf.PublicPath + `/edit/upload" class="btn new-post waves-effect waves-light">New upload</a></div></div>`
		html = html + b.String() + btn

		adminView, err := views.Admin(html, appName, pathConf)
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/html")
		res.Write(adminView)
	}
}
