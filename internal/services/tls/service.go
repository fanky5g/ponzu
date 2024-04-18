package tls

import (
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/models"
	"github.com/fanky5g/ponzu/tokens"
)

type service struct {
	configRepository driver.Repository
}

type Service interface {
	Enable()
	EnableDev()
}

func New(db driver.Database) (Service, error) {
	return &service{configRepository: db.Get(
		models.WrapPonzuModelNameSpace(tokens.ConfigRepositoryToken),
	).(driver.Repository)}, nil
}
