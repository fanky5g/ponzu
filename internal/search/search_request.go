package search

import "github.com/fanky5g/ponzu/internal/constants"

type PaginationRequestDto struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
}

type RequestDto struct {
	Query     string              `json:"query"`
	SortOrder constants.SortOrder `json:"order"`
	PaginationRequestDto
}
