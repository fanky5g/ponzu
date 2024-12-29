package content

import (
	"errors"
	"net/http"
)

var ErrUnsupportedMethod = errors.New("http method unsupported")

type ContentQuery struct {
	ID   string
	Type string
}

type ContentTransitionInput struct {
	ContentQuery
	TargetState string
}

type TabularDatasource interface {
	GetNumberOfRows() (int, error)
	GetColumns() ([]string, error)
	LoadData(offset int) ([]interface{}, error)
}

func MapContentQueryFromRequest(r *http.Request) (*ContentQuery, error) {
	method := r.Method

	switch method {
	case http.MethodGet:
		q := r.URL.Query()
		return &ContentQuery{
			ID:   q.Get("id"),
			Type: q.Get("type"),
		}, nil
	case http.MethodPost:
		return &ContentQuery{
			ID:   r.FormValue("id"),
			Type: r.FormValue("type"),
		}, nil
	default:
		return nil, ErrUnsupportedMethod
	}
}

func MapContentTransitionInputFromRequest(r *http.Request) (*ContentTransitionInput, error) {
	method := r.Method

	switch method {
	case http.MethodGet:
		q := r.URL.Query()
		return &ContentTransitionInput{
			ContentQuery: ContentQuery{
				ID:   q.Get("id"),
				Type: q.Get("type"),
			},
			TargetState: q.Get("workflow_state"),
		}, nil
	case http.MethodPost:
		return &ContentTransitionInput{
			ContentQuery: ContentQuery{
				ID:   r.FormValue("id"),
				Type: r.FormValue("type"),
			},
			TargetState: r.FormValue("workflow_state"),
		}, nil
	default:
		return nil, ErrUnsupportedMethod
	}
}
