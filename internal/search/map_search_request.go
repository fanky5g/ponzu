package search

import (
	"encoding/json"
	request2 "github.com/fanky5g/ponzu/internal/http/request"
	"net/http"
	"net/url"
	"strings"

	"github.com/fanky5g/ponzu/internal/constants"
)

func GetSearchRequestDto(req *http.Request) (*RequestDto, error) {
	switch request2.GetContentType(req) {
	case "multipart/form-data":
		if err := req.ParseMultipartForm(1024 * 1024 * 4); err != nil {
			return nil, err
		}

		return getSearchRequestFromURL(req.PostForm)
	case "application/json":
		var searchRequest RequestDto
		if err := json.NewDecoder(req.Body).Decode(&searchRequest); err != nil {
			return nil, err
		}

		return &searchRequest, nil
	default:
		return getSearchRequestFromURL(req.URL.Query())
	}
}

func getSearchRequestFromURL(qs url.Values) (*RequestDto, error) {
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

	return &RequestDto{
		Query:                q,
		SortOrder:            constants.SortOrder(order),
		PaginationRequestDto: *paginationRequest,
	}, nil
}

func MapSearchRequest(searchRequest *RequestDto) (*Search, error) {
	return &Search{
		Query:     searchRequest.Query,
		SortOrder: searchRequest.SortOrder,
		Count:     searchRequest.Count,
		Offset:    searchRequest.Offset,
	}, nil
}
