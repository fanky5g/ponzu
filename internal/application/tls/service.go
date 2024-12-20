package tls

import "github.com/fanky5g/ponzu/database"

type Service struct {
	config database.Repository
}

func New(config database.Repository) (*Service, error) {
	return &Service{config: config}, nil
}
