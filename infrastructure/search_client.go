package infrastructure

import (
	bleveSearch "github.com/fanky5g/ponzu-driver-bleve"

	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/driver"
	"github.com/pkg/errors"
)

var (
	ErrSearchDriverMissing     = errors.New("Invalid configuration: missing search_driver")
	ErrUnsupportedSearchDriver = errors.New("Unsupported search driver")
)

type searchConfig struct {
	database     driver.Database
	contentTypes map[string]func() interface{}
	cfg          *config.Config
}

func getSearchClient() (driver.SearchInterface, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	if cfg.SearchDriver == "" {
		return nil, ErrSearchDriverMissing
	}

	var searchClient driver.SearchInterface
	switch cfg.SearchDriver {
	case "postgres":
		// TODO: define integration for postgres search
		return nil, ErrUnsupportedSearchDriver
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
