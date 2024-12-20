// Package analytics provides the methods to run an analytics reporting system
// for API requests which may be useful to users for measuring access and
// possibly identifying bad actors abusing requests.
package analytics

import (
	"github.com/fanky5g/ponzu/database"
)

type Service struct {
	requests database.Repository
	metrics  database.Repository
}

func New(requests database.Repository, metrics database.Repository) (*Service, error) {
	return &Service{
		requests: requests,
		metrics:  metrics,
	}, nil
}
