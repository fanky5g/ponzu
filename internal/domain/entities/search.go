package entities

import "github.com/fanky5g/ponzu/internal/domain/enum"

type Search struct {
	Query      string
	SortOrder  enum.SortOrder
	Pagination *Pagination
}
