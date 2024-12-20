package database

import (
	"errors"

	postgres "github.com/fanky5g/ponzu-driver-postgres/database"
	"github.com/fanky5g/ponzu/config"
	databaseModels "github.com/fanky5g/ponzu/internal/infrastructure/database/models"
	"github.com/fanky5g/ponzu/models"
)

func New(contentModels []models.ModelInterface) (interface{}, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	m := make([]models.ModelInterface, 0)
	m = append(m, contentModels...)
	m = append(m, databaseModels.GetPonzuModels()...)

	switch cfg.DatabaseDriver {
	case "postgres":
		return postgres.New(m)
	default:
		return nil, errors.New("invalid driver")
	}
}
