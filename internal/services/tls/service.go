package tls

import (
	"github.com/fanky5g/ponzu/internal/domain/interfaces"
	"github.com/fanky5g/ponzu/internal/services"
)

var ServiceToken services.ServiceToken = "TlsService"

type service struct {
	configRepository interfaces.ConfigRepositoryInterface
}

type Service interface {
	Enable()
	EnableDev()
}

func New(repository interfaces.ConfigRepositoryInterface) (Service, error) {
	return &service{configRepository: repository}, nil
}
