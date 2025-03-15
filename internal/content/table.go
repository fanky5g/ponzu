package content

import (
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/fanky5g/ponzu/internal/datasource"
	"github.com/fanky5g/ponzu/internal/search"
	"math"
)

var PaginationOptions = []int{20, 50, 100}

type searchFunc func() ([]interface{}, int, error)

type TableViewModel struct {
	ItemType          interface{}
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
	Search            *search.Search
	CSVFormattable    bool
}

func buildTableViewModel(
	publicPath string,
	entity content.Entity,
	searchRequest *search.Search,
	s searchFunc,
) (*TableViewModel, error) {
	entityName := entity.EntityName()
	matches, total, err := s()
	if err != nil {
		return nil, err
	}

	// set up pagination values
	count := searchRequest.Count
	offset := searchRequest.Offset
	if total < count {
		count = total
	}

	start := 1 + offset
	end := start + len(matches) - 1

	if total < end {
		end = total
	}

	currentPage := int(math.Ceil(float64(start-1)/float64(count)) + 1)
	numberOfPages := int(math.Ceil(float64(total) / float64(count)))

	_, csvFormattable := entity.(datasource.Row)

	return &TableViewModel{
		PublicPath:        publicPath,
		ItemType:          entity,
		TableName:         fmt.Sprintf("%s Items", entityName),
		Items:             matches,
		TypeName:          entityName,
		RowsPerPage:       searchRequest.Count,
		TotalItemCount:    total,
		CurrentItemStart:  start,
		CurrentItemEnd:    end,
		NumberOfPages:     numberOfPages,
		CurrentPage:       currentPage,
		SortOrder:         searchRequest.SortOrder,
		PaginationOptions: PaginationOptions,
		Search:            searchRequest,
		CSVFormattable:    csvFormattable,
	}, nil
}
