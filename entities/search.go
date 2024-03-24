package entities

import (
	"github.com/fanky5g/ponzu/constants"
)

type Search struct {
	Query      string
	SortOrder  constants.SortOrder
	Pagination *Pagination
}
