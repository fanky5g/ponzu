package analytics

import (
	"github.com/fanky5g/ponzu/internal/analytics"
	"time"
)

func (s *service) getRequests(
	t time.Time,
	currentMetrics map[string]analytics.AnalyticsMetric,
) ([]analytics.AnalyticsHTTPRequestMetadata, error) {
	requests := make([]analytics.AnalyticsHTTPRequestMetadata, 0)
	rr, err := s.requestsRepository.FindAll()
	if err != nil {
		return nil, err
	}

	for _, r := range rr {
		request := r.(*analytics.AnalyticsHTTPRequestMetadata)
		// append request to requests for analysis if its timestamp is t
		// or if its day is not already in cache, otherwise delete it
		d := time.Unix(request.Timestamp/1000, 0)
		_, inCache := currentMetrics[d.Format("01/02")]
		if !d.Before(t) || !inCache {
			requests = append(requests, *request)
		} else {
			if err = s.requestsRepository.DeleteById(request.RequestID); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}
