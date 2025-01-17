package content

import (
	"errors"
	"net/http"
)

var ErrUnsupportedMethod = errors.New("http method unsupported")

type Query struct {
	ID   string
	Type string
}

type TransitionInput struct {
	Query
	TargetState string
}

type TabularDatasource interface {
	GetNumberOfRows() (int, error)
	GetColumns() ([]string, error)
	LoadData(offset int) ([]interface{}, error)
}

func MapContentQueryFromRequest(r *http.Request) (*Query, error) {
	method := r.Method

	switch method {
	case http.MethodGet:
		q := r.URL.Query()
		return &Query{
			ID:   q.Get("id"),
			Type: q.Get("type"),
		}, nil
	case http.MethodPost:
		return &Query{
			ID:   r.FormValue("id"),
			Type: r.FormValue("type"),
		}, nil
	default:
		return nil, ErrUnsupportedMethod
	}
}

func MapContentTransitionInputFromRequest(r *http.Request) (*TransitionInput, error) {
	method := r.Method

	switch method {
	case http.MethodGet:
		q := r.URL.Query()
		return &TransitionInput{
			Query: Query{
				ID:   q.Get("id"),
				Type: q.Get("type"),
			},
			TargetState: q.Get("workflow_state"),
		}, nil
	case http.MethodPost:
		return &TransitionInput{
			Query: Query{
				ID:   r.FormValue("id"),
				Type: r.FormValue("type"),
			},
			TargetState: r.FormValue("workflow_state"),
		}, nil
	default:
		return nil, ErrUnsupportedMethod
	}
}
