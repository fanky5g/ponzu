package table

import (
	"fmt"
	"math"

	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/entities"
)

var PaginationOptions = []int{20, 50, 100}

type RowLoader func() ([]interface{}, int, error)

type Table struct {
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
	Search            *entities.Search
	CSVFormattable    bool
}

func New(
	t string,
	itemType interface{},
	search *entities.Search,
	loader RowLoader,
) (*Table, error) {
	matches, total, err := loader()
	if err != nil {
		return nil, err
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

	currentPage := int(math.Ceil(float64(start-1)/float64(count)) + 1)
	numberOfPages := int(math.Ceil(float64(total) / float64(count)))

	_, csvFormattable := itemType.(item.CSVFormattable)

	return &Table{
		ItemType:          itemType,
		TableName:         fmt.Sprintf("%s Items", t),
		Items:             matches,
		TypeName:          t,
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
	}, err
}
