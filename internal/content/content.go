package content

import "net/http"

type ContentQuery struct {
	ID   string
	Type string
}

type TabularDatasource interface {
	GetNumberOfRows() (int, error)
	GetColumns() ([]string, error)
	LoadData(offset int) ([]interface{}, error)
}

func MapContentQueryFromRequest(r *http.Request) (*ContentQuery, error) {
	q := r.URL.Query()

	return &ContentQuery{
		ID:   q.Get("id"),
		Type: q.Get("type"),
	}, nil
}
