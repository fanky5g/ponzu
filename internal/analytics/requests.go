package analytics

import (
	"time"
)

func (s *Service) getRequests(
	t time.Time,
	currentMetrics map[string]Metric,
) ([]HTTPRequestMetadata, error) {
	requests := make([]HTTPRequestMetadata, 0)
	rr, err := s.requestsRepository.FindAll()
	if err != nil {
		return nil, err
	}

	for _, r := range rr {
		request := r.(*HTTPRequestMetadata)
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
