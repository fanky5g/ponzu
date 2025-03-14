// Package analytics provides the methods to run an analytics reporting system
// for API requests which may be useful to users for measuring access and
// possibly identifying bad actors abusing requests.
package analytics

import (
	"github.com/fanky5g/ponzu/internal/database"
)

type Service struct {
	requestsRepository database.Repository
	metricsRepository  database.Repository
}

func New(db database.Database) (*Service, error) {
	return &Service{
		requestsRepository: db.GetRepositoryByToken(RequestsRepositoryToken),
		metricsRepository:  db.GetRepositoryByToken(MetricsRepositoryToken),
	}, nil
}
