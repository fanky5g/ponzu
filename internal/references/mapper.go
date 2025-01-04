package references

import (
	"errors"
	"net/http"

	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/exceptions"
	"github.com/fanky5g/ponzu/internal/http/request"
)

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

	searchRequestDto, err := request.GetSearchRequestDto(req)
	if err != nil {
		return nil, err
	}

	searchRequest, err := request.MapSearchRequest(searchRequestDto)
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
