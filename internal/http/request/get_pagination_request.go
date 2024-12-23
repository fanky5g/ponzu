package request

import (
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/request"
	"net/url"
	"strconv"
)

var RowsPerPage = 20

func GetPaginationRequest(qs url.Values) (*request.PaginationRequestDto, error) {
	count, err := strconv.Atoi(qs.Get("count"))
	if err != nil {
		if qs.Get("count") == "" {
			count = RowsPerPage
		} else {
			return nil, err
		}
	}

	offset, err := strconv.Atoi(qs.Get("offset"))
	if err != nil {
		if qs.Get("offset") == "" {
			offset = 0
		} else {
			return nil, err
		}
	}

	return &request.PaginationRequestDto{
		Count:  count,
		Offset: offset,
	}, nil
}
