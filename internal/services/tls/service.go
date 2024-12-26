package tls

import (
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/tokens"
)

type service struct {
	configRepository database.Repository
}

type Service interface {
	Enable()
	EnableDev()
}

func New(db database.Database) (Service, error) {
	return &service{configRepository: db.GetRepositoryByToken(tokens.ConfigRepositoryToken)}, nil
}
