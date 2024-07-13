package search

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/pkg/errors"
)

type service struct {
	client   driver.SearchInterface
	database driver.Database
}

type Service interface {
	Search(entity interface{}, query string, count, offset int) ([]interface{}, int, error)
}

func New(client driver.SearchInterface, database driver.Database) (Service, error) {
	return &service{client: client, database: database}, nil
}

// Search conducts a search and returns a set of content documents after loading from database
// if search driver supports GetID methods on returned matches. Otherwise, plain Ponzu targets Type:ID pairs are returned
func (s *service) Search(entity interface{}, query string, count, offset int) ([]interface{}, int, error) {
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

	e, ok := entity.(content.Entity)
	if !ok {
		return matches, size, nil
	}

	repository := s.database.GetRepositoryByToken(tokens.RepositoryToken(e.EntityName()))
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
