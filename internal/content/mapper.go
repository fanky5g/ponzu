package content

import (
	"errors"
	"github.com/fanky5g/ponzu/internal/search"
	"net/http"

	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/exceptions"
)

var ErrUnsupportedMethod = errors.New("http method unsupported")

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) MapEntityToReference(entity interface{}) (*Reference, error) {
	if entity == nil {
		return nil, nil
	}

	identifiable, ok := entity.(item.Identifiable)
	if !ok {
		return nil, errors.New("reference entity invalid")
	}

	itemName := identifiable.ItemID()
	readable, ok := entity.(item.Readable)
	if ok {
		itemName = readable.GetTitle()
	}

	return &Reference{
		ID:   identifiable.ItemID(),
		Name: itemName,
	}, nil
}

func (m *Mapper) MapReferenceMatchesToReferencesResourceOutput(
	matches []interface{},
	size int,
) (*ListReferencesOutputResource, error) {
	rs := make([]interface{}, len(matches))
	var err error
	for i := range matches {
		rs[i], err = m.MapEntityToReference(matches[i])
		if err != nil {
			return nil, err
		}
	}

	return &ListReferencesOutputResource{
		References: rs,
		Size:       size,
	}, nil
}

func (m *Mapper) MapReqToListReferencesInputResource(req *http.Request) (*ListReferencesInputResource, error) {
	referenceType := req.URL.Query().Get("type")
	if referenceType == "" {
		return nil, errors.New("invalid reference type")
	}

	searchRequestDto, err := search.GetSearchRequestDto(req)
	if err != nil {
		return nil, err
	}

	searchRequest, err := search.MapSearchRequest(searchRequestDto)
	if err != nil {
		return nil, err
	}

	return &ListReferencesInputResource{
		Type:   referenceType,
		Search: searchRequest,
	}, nil
}

func (m *Mapper) MapRequestToGetReferenceInputResource(req *http.Request) (*GetReferenceInputResource, error) {
	q := req.URL.Query()
	id := req.PathValue("id")

	typeName := q.Get("type")
	if typeName == "" {
		return nil, exceptions.NewClientException(
			"Invalid reference type",
			"",
			exceptions.INFO,
			nil,
			nil,
		)
	}

	return &GetReferenceInputResource{
		ID:   id,
		Type: typeName,
	}, nil
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
