package request

import (
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/request"
	"net/url"
	"strconv"
)

func GetPaginationRequest(qs url.Values) (*request.PaginationRequestDto, error) {
	count, err := strconv.Atoi(qs.Get("count")) // int: determines number of posts to return (10 default, -1 is all)
	if err != nil {
		if qs.Get("count") == "" {
			count = 10
		} else {
			return nil, err
		}
	}

	offset, err := strconv.Atoi(qs.Get("offset")) // int: multiplier of count for pagination (0 default)
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
