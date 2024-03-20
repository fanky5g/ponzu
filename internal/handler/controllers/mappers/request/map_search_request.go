package request

import (
	"encoding/json"
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/request"
	"net/http"
	"net/url"
	"strings"
)

func GetSearchRequestDto(req *http.Request) (*request.SearchRequestDto, error) {
	switch getContentType(req) {
	case "multipart/form-data":
		if err := req.ParseMultipartForm(1024 * 1024 * 4); err != nil {
			return nil, err
		}

		return getSearchRequestFromURL(req.PostForm)
	case "services/json":
		var searchRequest request.SearchRequestDto
		if err := json.NewDecoder(req.Body).Decode(&searchRequest); err != nil {
			return nil, err
		}

		return &searchRequest, nil
	default:
		return getSearchRequestFromURL(req.URL.Query())
	}
}

func getSearchRequestFromURL(qs url.Values) (*request.SearchRequestDto, error) {
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
		SortOrder:            constants.SortOrder(order),
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
