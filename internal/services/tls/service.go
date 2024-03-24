package tls

import (
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
	"github.com/fanky5g/ponzu/tokens"
)

type service struct {
	configRepository repositories.ConfigRepositoryInterface
}

type Service interface {
	Enable()
	EnableDev()
}

func New(db driver.Database) (Service, error) {
	return &service{configRepository: db.Get(tokens.ConfigRepositoryToken).(repositories.ConfigRepositoryInterface)}, nil
}
