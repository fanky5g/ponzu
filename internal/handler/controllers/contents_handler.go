package controllers

import (
	"bytes"
	"fmt"
	"math"

	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content/editor"
	// "github.com/fanky5g/ponzu/content/item"
	"net/http"

	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/tokens"

	log "github.com/sirupsen/logrus"
)

var PaginationOptions = []int{20, 50, 100}

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

		var total int
		var posts []interface{}
		total, posts, err = contentService.GetAllWithOptions(t, search)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to search")
			r.Renderer().InternalServerError(res)
			return
		}

		//		html := `<div class="col s9 card">
		//					<div class="card-content">
		//					<div class="row">
		//					</div>
		//					<form class="col s4" action="{{ .PublicPath }}/contents/search" method="get">
		//						<div class="input-field post-search inline">
		//							<label class="active">Search:</label>
		//							<i class="right material-icons search-icon">search</i>
		//							<input class="search" name="q" type="text" placeholder="Within all {{ .Data.TypeName }} fields" class="search"/>
		//							<input type="hidden" name="type" value="{{ .Data.TypeName }}" />
		//						</div>
		//                    </form>
		//					</div>`
		//

		//		html += `<ul class="posts row">`
		//
		//		_, err = b.Write([]byte(`</ul>`))
		//		if err != nil {
		//			log.WithField("Error", err).Warning("Failed to write buffer")
		//			r.Renderer().InternalServerError(res)
		//			return
		//		}
		//
		// set up pagination values
		count := search.Pagination.Count
		offset := search.Pagination.Offset
		if total < count {
			count = total
		}

		start := 1 + offset
		end := start + len(posts) - 1

		if total < end {
			end = total
		}

		//		script := `
		//	<script>
		//		$(function() {
		//			var del = $('.quick-delete-post.__ponzu span');
		//			del.on('click', function(e) {
		//				if (confirm("[Ponzu] Please confirm:\n\nAre you sure you want to delete this post?\nThis cannot be undone.")) {
		//					$(e.target).parent().submit();
		//				}
		//			});
		//		});
		//
		//		// disable link from being clicked if parent is 'disabled'
		//		$(function() {
		//			$('ul.pagination li.disabled a').on('click', function(e) {
		//				e.preventDefault();
		//			});
		//		});
		//	</script>
		//	`
		//
		//		btn := `<div class="col s3">
		//		<a href="{{ .PublicPath }}/edit?type={{ .Data.TypeName }}" class="btn new-post waves-effect waves-light">
		//			New {{ .Data.TypeName }}
		//		</a>`
		//
		//		if _, ok := pt.(item.CSVFormattable); ok {
		//			btn += `<br/>
		//				<a href="{{ .PublicPath }}/contents/export?type={{ .Data.TypeName }}&format=csv" class="green darken-4 btn export-post waves-effect waves-light">
		//					<i class="material-icons left">system_update_alt</i>
		//					CSV
		//				</a>`
		//		}
		//
		//		html += b.String() + script + btn + `</div></div>`
		buf := &bytes.Buffer{}
		tableViewTmpl := r.Renderer().TemplateFromDir("datatable")

		currentPage := int(math.Ceil(float64(start-1)/float64(count)) + 1)
		numberOfPages := int(math.Ceil(float64(total) / float64(count)))

		data := struct {
			TableName         string
			Items             []interface{}
			TypeName          string
			PublicPath        string
			RowsPerPage       int
			SortOrder         constants.SortOrder
			TotalItemCount    int
			CurrentItemStart  int
			CurrentItemEnd    int
			CurrentPage       int
			NumberOfPages     int
			PaginationOptions []int
		}{
			TableName:         fmt.Sprintf("%s Items", t),
			Items:             posts,
			TypeName:          t,
			PublicPath:        r.Context().Paths().PublicPath,
			RowsPerPage:       search.Pagination.Count,
			TotalItemCount:    total,
			CurrentItemStart:  start,
			CurrentItemEnd:    end,
			NumberOfPages:     numberOfPages,
			CurrentPage:       currentPage,
			SortOrder:         search.SortOrder,
			PaginationOptions: PaginationOptions,
		}

		if err := tableViewTmpl.Execute(buf, data); err != nil {
			log.WithField("Error", err).Warning("Failed to write buffer")
			r.Renderer().InternalServerError(res)
			return
		}

		r.Renderer().InjectInAdminView(res, buf)
	}
}
