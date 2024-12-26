package infrastructure

import (
	"errors"

	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/internal/database/postgres/models"
)

func getDatabaseDriver(name string, contentModels []database.ModelInterface) (driver.Database, error) {
	m := make([]database.ModelInterface, 0)
	m = append(m, contentModels...)
	m = append(m, models.GetPonzuModels()...)

	switch name {
	default:
		return nil, errors.New("invalid driver")
	}
}
