package request

import "github.com/fanky5g/ponzu/internal/domain/enum"

type PaginationRequestDto struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
}

type SearchRequestDto struct {
	Query     string         `json:"query"`
	SortOrder enum.SortOrder `json:"order"`
	PaginationRequestDto
}
