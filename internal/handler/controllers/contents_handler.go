package controllers

import (
	"bytes"
	"fmt"
	"math"

	"github.com/fanky5g/ponzu/content/editor"
	// "github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/tokens"
	"net/http"

	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

var RowsPerPage = 50

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

		order := strings.ToLower(q.Get("order"))
		if order != "asc" {
			order = "desc"
		}

		status := q.Get("status")
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

		count, err := strconv.Atoi(q.Get("count")) // int: determines number of posts to return (10 default, -1 is all)
		if err != nil {
			if q.Get("count") == "" {
				count = RowsPerPage
			} else {
				log.WithField("Error", err).Warning("Failed to parse count")
				r.Renderer().InternalServerError(res)
				return
			}
		}

		offset, err := strconv.Atoi(q.Get("offset")) // int: multiplier of count for pagination (0 default)
		if err != nil {
			if q.Get("offset") == "" {
				offset = 0
			} else {
				log.WithField("Error", err).Warning("Failed to parse offset")
				r.Renderer().InternalServerError(res)
				return
			}
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
		//					<div class="col s8">
		//						<div class="row">
		//							<div class="col s5 input-field inline">
		//								<select class="browser-default __ponzu sort-order">
		//									<option value="DESC">New to Old</option>
		//									<option value="ASC">Old to New</option>
		//								</select>
		//								<label class="active">Sort:</label>
		//							</div>
		//							<script>
		//								$(function() {
		//									var sort = $('select.__ponzu.sort-order');
		//
		//									sort.on('change', function() {
		//										var path = window.location.pathname;
		//										var s = sort.val();
		//										var t = getParam('type');
		//
		//										if (status == "") {
		//											status = "public";
		//										}
		//
		//										window.location.replace(path + '?type=' + t + '&order=' + s);
		//									});
		//
		//									var order = getParam('order');
		//									if (order !== '') {
		//										sort.val(order);
		//									}
		//
		//								});
		//							</script>
		//						</div>
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
		//		statusDisabled := "disabled"
		prevStatus := ""
		nextStatus := ""
		//		// total may be less than 10 (default count), so reset count to match total
		//		if total < count {
		//			count = total
		//		}
		//		// nothing previous to current list
		//		if offset == 0 {
		//			prevStatus = statusDisabled
		//		}
		//		// nothing after current list
		//		if (offset+1)*count >= total {

		//			nextStatus = statusDisabled
		//		}
		//
		// set up pagination values
		urlFmt := req.URL.Path + "?count=%d&offset=%d&&order=%s&status=%s&type=%s"
		prevURL := fmt.Sprintf(urlFmt, count, offset-1, order, status, t)
		nextURL := fmt.Sprintf(urlFmt, count, offset+1, order, status, t)
		start := 1 + count*offset
		end := start + count - 1

		if total < end {
			end = total
		}

		_ = fmt.Sprintf(`
		<ul class="pagination row">
			<li class="col s2 waves-effect %s"><a href="%s"><i class="material-icons">chevron_left</i></a></li>
			<li class="col s8">%d to %d of %d</li>
			<li class="col s2 waves-effect %s"><a href="%s"><i class="material-icons">chevron_right</i></a></li>
		</ul>
		`, prevStatus, prevURL, start, end, total, nextStatus, nextURL)

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
			TableName        string
			Items            []interface{}
			TypeName         string
			PublicPath       string
			NextPageURL      string
			PrevPageURL      string
			RowsPerPage      int
			TotalItemCount   int
			CurrentItemStart int
			CurrentItemEnd   int
			CurrentPage      int
			NumberOfPages    int
		}{
			TableName:        fmt.Sprintf("%s Items", t),
			Items:            posts,
			TypeName:         t,
			PublicPath:       r.Context().Paths().PublicPath,
			NextPageURL:      nextURL,
			PrevPageURL:      prevURL,
			RowsPerPage:      count,
			TotalItemCount:   total,
			CurrentItemStart: start,
			CurrentItemEnd:   end,
			NumberOfPages:    numberOfPages,
			CurrentPage:      currentPage,
		}

		if err := tableViewTmpl.Execute(buf, data); err != nil {
			log.WithField("Error", err).Warning("Failed to write buffer")
			r.Renderer().InternalServerError(res)
			return
		}

		r.Renderer().InjectInAdminView(res, buf)
	}
}
