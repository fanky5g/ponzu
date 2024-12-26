// Package analytics provides the methods to run an analytics reporting system
// for API requests which may be useful to users for measuring access and
// possibly identifying bad actors abusing requests.
package analytics

import (
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/internal/analytics"
	"github.com/fanky5g/ponzu/tokens"
)

type service struct {
	requestsRepository driver.Repository
	metricsRepository  driver.Repository
}

type Service interface {
	StartRecorder()
	Record(req analytics.AnalyticsHTTPRequestMetadata)
	GetChartData() (map[string]interface{}, error)
}

func New(db driver.Database) (Service, error) {
	return &service{
		requestsRepository: db.GetRepositoryByToken(tokens.AnalyticsRequestsRepositoryToken),
		metricsRepository:  db.GetRepositoryByToken(tokens.AnalyticsMetricsRepositoryToken),
	}, nil
}
