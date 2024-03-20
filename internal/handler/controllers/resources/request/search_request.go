package request

import "github.com/fanky5g/ponzu/constants"

type PaginationRequestDto struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
}

type SearchRequestDto struct {
	Query     string              `json:"query"`
	SortOrder constants.SortOrder `json:"order"`
	PaginationRequestDto
}
