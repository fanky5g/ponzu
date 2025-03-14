package search

import (
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/pkg/errors"
)

type Service struct {
	client   SearchInterface
	database database.Database
}

func New(client SearchInterface, database database.Database) (*Service, error) {
	return &Service{client: client, database: database}, nil
}

// Search conducts a search and returns a set of content documents after loading from database
// if search driver supports GetID methods on returned matches. Otherwise, plain Ponzu targets Type:ID pairs are returned
func (s *Service) Search(entity interface{}, query string, count, offset int) ([]interface{}, int, error) {
	matches, size, err := s.client.SearchWithPagination(entity, query, count, offset)
	if err != nil {
		return nil, 0, err
	}

	if len(matches) == 0 {
		return nil, 0, nil
	}

	_, ok := matches[0].(interface {
		GetID() string
	})

	if !ok {
		return matches, size, nil
	}

	persistable, ok := entity.(database.Persistable)
	if !ok {
		return matches, size, nil
	}

	repository := s.database.GetRepositoryByToken(persistable.GetRepositoryToken())
	results := make([]interface{}, len(matches))
	for i := range matches {
		identifiable := matches[i].(interface {
			GetID() string
		})

		result, err := repository.FindOneById(identifiable.GetID())
		if err != nil {
			return nil, 0, errors.Wrap(err, "Failed to fetch document")
		}

		results[i] = result
	}

	return results, size, nil
}
