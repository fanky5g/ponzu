package controllers

import (
	"bytes"
	"fmt"
	"math"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
)

var PaginationOptions = []int{20, 50, 100}

type ResultLoader func() ([]interface{}, int, error)

func renderContentList(
	r router.Router,
	res http.ResponseWriter,
	t string,
	search *entities.Search,
	itemType interface{},
	loader ResultLoader,
) {
	matches, total, err := loader()
	if err != nil {
		log.WithField("Error", err).Warning("Failed to search")
		r.Renderer().InternalServerError(res)
		return

	}

	// set up pagination values
	count := search.Pagination.Count
	offset := search.Pagination.Offset
	if total < count {
		count = total
	}

	start := 1 + offset
	end := start + len(matches) - 1

	if total < end {
		end = total
	}

	buf := &bytes.Buffer{}
	tableViewTmpl := r.Renderer().TemplateFromDir("datatable")

	currentPage := int(math.Ceil(float64(start-1)/float64(count)) + 1)
	numberOfPages := int(math.Ceil(float64(total) / float64(count)))

	_, csvFormattable := itemType.(item.CSVFormattable)

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
		Search            *entities.Search
		CSVFormattable    bool
	}{
		TableName:         fmt.Sprintf("%s Items", t),
		Items:             matches,
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
		Search:            search,
		CSVFormattable:    csvFormattable,
	}

	if err := tableViewTmpl.Execute(buf, data); err != nil {
		log.WithField("Error", err).Warning("Failed to write buffer")
		r.Renderer().InternalServerError(res)
		return
	}

	r.Renderer().InjectInAdminView(res, buf)
}
