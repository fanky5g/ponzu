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
	repository repositories.AnalyticsRepositoryInterface
}

type Service interface {
	StartRecorder(analyticsRepository repositories.AnalyticsRepositoryInterface)
	Record(req entities.AnalyticsHTTPRequestMetadata)
	GetChartData() (map[string]interface{}, error)
}

func New(db driver.Database) (Service, error) {
	return &service{repository: db.Get(tokens.AnalyticsRepositoryToken).(repositories.AnalyticsRepositoryInterface)}, nil
}
