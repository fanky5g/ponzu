package analytics

import (
	"time"

	"github.com/fanky5g/ponzu/internal/entities"
)

func (s *Service) getRequests(
	t time.Time,
	currentMetrics map[string]entities.AnalyticsMetric,
) ([]entities.AnalyticsHTTPRequestMetadata, error) {
	requests := make([]entities.AnalyticsHTTPRequestMetadata, 0)
	rr, err := s.requests.FindAll()
	if err != nil {
		return nil, err
	}

	for _, r := range rr {
		request := r.(*entities.AnalyticsHTTPRequestMetadata)
		// append request to requests for analysis if its timestamp is t
		// or if its day is not already in cache, otherwise delete it
		d := time.Unix(request.Timestamp/1000, 0)
		_, inCache := currentMetrics[d.Format("01/02")]
		if !d.Before(t) || !inCache {
			requests = append(requests, *request)
		} else {
			if err = s.requests.DeleteById(request.RequestID); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}
