package controllers

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewUploadContentsHandler(r router.Router) http.HandlerFunc {
	storageService := r.Context().Service(tokens.StorageServiceToken).(storage.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		pt := interface{}(&entities.FileUpload{})
		_, ok := pt.(editor.Editable)
		if !ok {
			log.Warning("entities.FileUpload is not editable")
			r.Renderer().InternalServerError(res)
			return
		}

		searchRequestDto, err := request.GetSearchRequestDto(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get search request dto")
			return
		}

		search, err := request.MapSearchRequest(searchRequestDto)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request dto")
			return
		}

		b := &bytes.Buffer{}
		var total int
		var posts []*entities.FileUpload

		html := `<div class="col s9 card">
					<div class="card-content">
					<div class="row">
					<div class="col s8">
						<div class="row">
							<div class="card-title col s7">Uploaded Items</div>
							<div class="col s5 input-field inline">
								<select class="browser-default __ponzu sort-order">
									<option value="DESC">New to Old</option>
									<option value="ASC">Old to New</option>
								</select>
								<label class="active">Sort:</label>
							</div>
							<script>
								$(function() {
									var sort = $('select.__ponzu.sort-order');

									sort.on('change', function() {
										var path = window.location.pathname;
										var s = sort.val();

										window.location.replace(path + '?order=' + s);
									});

									var order = getParam('order');
									if (order !== '') {
										sort.val(order);
									}

								});
							</script>
						</div>
					</div>
					<form class="col s4" action="{{ .PublicPath }}/uploads/search" method="get">
						<div class="input-field post-search inline">
							<label class="active">Find:</label>
							<i class="right material-icons search-icon">search</i>
							<input class="search" name="q" type="text" placeholder="Within all upload fields" class="search"/>
							<input type="hidden" name="type" value="__uploads" />
						</div>
                   </form>
					</div>`

		total, posts, err = storageService.GetAllWithOptions(search)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to search uploads")
			return
		}

		for i := range posts {
			p, ok := interface{}(posts[i]).(editor.Editable)
			if !ok {
				log.Printf("Invalid entry. Item %v does not implement editable interface\n", posts[i])

				post := `<li class="col s12">Error decoding data. Possible file corruption.</li>`
				_, err = b.Write([]byte(post))
				if err != nil {
					log.WithField("Error", err).Warning("Failed to write template")
					r.Renderer().InternalServerError(res)
					return
				}

				continue
			}

			contentEntryTemplate := editor.BuildContentListEntryTemplate(p, constants.UploadsEntityName)
			_, err = b.Write([]byte(contentEntryTemplate))
			if err != nil {
				log.WithField("Error", err).Warning("Failed to write template")
				r.Renderer().InternalServerError(res)
				return
			}
		}

		html += `<ul class="posts row">`

		_, err = b.Write([]byte(`</ul>`))
		if err != nil {
			log.WithField("Error", err).Warning("Failed to write template")
			r.Renderer().InternalServerError(res)
			return
		}

		statusDisabled := "disabled"
		prevStatus := ""
		nextStatus := ""
		// total may be less than 10 (default count), so reset count to match total
		if total < search.Pagination.Count {
			search.Pagination.Count = total
		}
		// nothing previous to current list
		if search.Pagination.Offset == 0 {
			prevStatus = statusDisabled
		}
		// nothing after current list
		if (search.Pagination.Offset+1)*search.Pagination.Count >= total {
			nextStatus = statusDisabled
		}

		// set up pagination values
		urlFmt := req.URL.Path + "?count=%d&offset=%d&&order=%s"
		prevURL := fmt.Sprintf(urlFmt, search.Pagination.Count, search.Pagination.Offset-1, search.SortOrder)
		nextURL := fmt.Sprintf(urlFmt, search.Pagination.Count, search.Pagination.Offset+1, search.SortOrder)
		start := 1 + search.Pagination.Count*search.Pagination.Offset
		end := start + search.Pagination.Count - 1

		if total < end {
			end = total
		}

		pagination := fmt.Sprintf(`
	<ul class="pagination row">
		<li class="col s2 waves-effect %s"><a href="%s"><i class="material-icons">chevron_left</i></a></li>
		<li class="col s8">%d to %d of %d</li>
		<li class="col s2 waves-effect %s"><a href="%s"><i class="material-icons">chevron_right</i></a></li>
	</ul>
	`, prevStatus, prevURL, start, end, total, nextStatus, nextURL)

		// show indicator that a collection of items will be listed implicitly, but
		// that none are created yet
		if total < 1 {
			pagination = `
		<ul class="pagination row">
			<li class="col s2 waves-effect disabled"><a href="#"><i class="material-icons">chevron_left</i></a></li>
			<li class="col s8">0 to 0 of 0</li>
			<li class="col s2 waves-effect disabled"><a href="#"><i class="material-icons">chevron_right</i></a></li>
		</ul>
		`
		}

		_, err = b.Write([]byte(pagination + `</div></div>`))
		if err != nil {
			log.WithField("Error", err).Warning("Failed to write template")
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

		btn := `<div class="col s3"><a href="{{ .PublicPath }}/edit/upload" class="btn new-post waves-effect waves-light">New upload</a></div></div>`
		html = html + b.String() + script + btn

		r.Renderer().InjectTemplateInAdmin(res, html, nil)
	}
}
