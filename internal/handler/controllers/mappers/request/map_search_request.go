package request

import (
	"github.com/fanky5g/ponzu/internal/domain/entities"
	"github.com/fanky5g/ponzu/internal/domain/enum"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/request"
	"net/url"
	"strings"
)

func GetSearchRequestDto(qs url.Values) (*request.SearchRequestDto, error) {
	q, err := url.QueryUnescape(qs.Get("q"))
	if err != nil {
		return nil, err
	}

	order := strings.ToLower(qs.Get("order"))
	if order != "asc" {
		order = "desc"
	}

	paginationRequest, err := GetPaginationRequest(qs)
	if err != nil {
		return nil, err
	}

	return &request.SearchRequestDto{
		Query:                q,
		SortOrder:            enum.SortOrder(order),
		PaginationRequestDto: *paginationRequest,
	}, nil
}

func MapSearchRequest(searchRequest *request.SearchRequestDto) (*entities.Search, error) {
	return &entities.Search{
		Query:     searchRequest.Query,
		SortOrder: searchRequest.SortOrder,
		Pagination: &entities.Pagination{
			Count:  searchRequest.Count,
			Offset: searchRequest.Offset,
		},
	}, nil
}
