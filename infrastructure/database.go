package infrastructure

import (
	"errors"

	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/models"

	postgres "github.com/fanky5g/ponzu-driver-postgres/database"
)

func getDatabaseDriver(name string, contentModels []models.ModelInterface) (driver.Database, error) {
	m := make([]models.ModelInterface, 0)
	m = append(m, contentModels...)
	m = append(m, models.GetPonzuModels()...)

	switch name {
	case "postgres":
		return postgres.New(m)
	default:
		return nil, errors.New("invalid driver")
	}
}
