package infrastructure

import (
	bleveSearch "github.com/fanky5g/ponzu-driver-bleve"

	postgres "github.com/fanky5g/ponzu-driver-postgres/database"
	postgresSearch "github.com/fanky5g/ponzu-driver-postgres/search"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/search"
	"github.com/pkg/errors"
)

var (
	ErrSearchDriverMissing     = errors.New("Invalid configuration: missing search_driver")
	ErrUnsupportedSearchDriver = errors.New("Unsupported search driver")
)

func New(db database.Database) (search.Client, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	if cfg.SearchDriver == "" {
		return nil, ErrSearchDriverMissing
	}

	var searchClient search.Client 
	switch cfg.SearchDriver {
	case "postgres":
		postgresDb, ok := db.(*postgres.Database)
		if !ok {
			return nil, errors.New("database driver incompatible with postgres search driver")
		}

		return postgresSearch.New(postgresDb)
	case "bleve":
		searchClient, err = bleveSearch.New(cfg.Paths.DataDir)
		if err != nil {
			return nil, errors.Wrap(err, "failed to initialize search client")
		}
	default:
		return nil, ErrUnsupportedSearchDriver
	}

	return searchClient, nil
}
