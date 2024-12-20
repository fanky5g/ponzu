package entities

import (
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/search"
)

type Search struct {
	Query      string
	SortOrder  constants.SortOrder
	Pagination *search.Pagination
}
