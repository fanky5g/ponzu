package controllers

import (
	"bytes"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/search"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewSearchHandler(r router.Router) http.HandlerFunc {
	searchService := r.Context().Service(tokens.ContentSearchServiceToken).(search.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")

		searchRequest, err := request.GetSearchRequestDto(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request DTO")
			r.Renderer().InternalServerError(res)
			return
		}

		// Query must be set
		if searchRequest.Query == "" {
			r.Renderer().BadRequest(res)
			return
		}

		if t == "" {
			r.Redirect(req, res, "/admin")
			return
		}

        contentTypeConstructor, ok := r.Context().Types().Content[t]
        if !ok {
            r.Redirect(req, res, "/admin")
            return
        }

		// TODO: implement pagination with response size
		matches, _, err := searchService.Search(contentTypeConstructor(), searchRequest.Query, searchRequest.Count, searchRequest.Offset)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to search")
			return
		}

		b := &bytes.Buffer{}
		html := `<div class="col s9 card">
					<div class="card-content">
					<div class="row">
					<div class="card-title col s7">` + t + ` Results</div>
					<form class="col s4" action="{{ .PublicPath }}/contents/search" method="get">
						<div class="input-field post-search inline">
							<label class="active">Find:</label>
							<i class="right material-icons search-icon">search</i>
							<input class="search" name="q" type="text" placeholder="Within all ` + t + ` fields" class="search"/>
							<input type="hidden" name="type" value="` + t + `" />
						</div>
                   </form>
					</div>
					<ul class="posts row">`

		for i := range matches {
			contentEntryTemplate := editor.BuildContentListEntryTemplate(matches[i].(editor.Editable), t)
			_, err = b.Write([]byte(contentEntryTemplate))
			if err != nil {
				log.WithField("Error", err).Warning("Failed to write template")
				r.Renderer().InternalServerError(res)
				return
			}
		}

		_, err = b.WriteString(`</ul></div></div>`)
		if err != nil {
			r.Renderer().InternalServerError(res)
			return
		}

		script := `
	<script>
		$(function() {
			var del = $('.quick-delete-post.__ponzu span');
			del.on('click', function(e) {
				if (confirm("[Ponzu] Please confirm:\n\nAre you sure you want to delete this post?\nThis cannot be undone.")) {
					$(e.target).parent().submit();
				}
			});
		});

		// disable link from being clicked if parent is 'disabled'
		$(function() {
			$('ul.pagination li.disabled a').on('click', function(e) {
				e.preventDefault();
			});
		});
	</script>
	`

		btn := `<div class="col s3">
		<a href="{{ .PublicPath }}/edit?type=` + t + `" class="btn new-post waves-effect waves-light">
			New ` + t + `
		</a>`

		html += b.String() + script + btn + `</div></div>`
		r.Renderer().InjectTemplateInAdmin(res, html, nil)
	}
}
