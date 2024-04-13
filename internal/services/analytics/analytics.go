// Package analytics provides the methods to run an analytics reporting system
// for API requests which may be useful to users for measuring access and
// possibly identifying bad actors abusing requests.
package analytics

import (
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
	"github.com/fanky5g/ponzu/tokens"
)

type service struct {
	requestsRepository repositories.GenericRepositoryInterface
	metricsRepository  repositories.GenericRepositoryInterface
}

type Service interface {
	StartRecorder()
	Record(req entities.AnalyticsHTTPRequestMetadata)
	GetChartData() (map[string]interface{}, error)
}

func New(db driver.Database) (Service, error) {
	return &service{
		requestsRepository: db.Get(tokens.AnalyticsRequestsRepositoryToken).(repositories.GenericRepositoryInterface),
		metricsRepository:  db.Get(tokens.AnalyticsMetricsRepositoryToken).(repositories.GenericRepositoryInterface),
	}, nil
}
