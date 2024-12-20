package analytics

import "github.com/fanky5g/ponzu/internal/entities"

func (s *Service) getMetrics() (map[string]entities.AnalyticsMetric, error) {
	metrics, err := s.metrics.FindAll()
	if err != nil {
		return nil, err
	}

	analyticsMetrics := make(map[string]entities.AnalyticsMetric)
	for _, m := range metrics {
		metric := m.(*entities.AnalyticsMetric)
		// add metric to currentMetrics map
		analyticsMetrics[metric.Date] = *metric
	}

	return analyticsMetrics, nil
}

func (s *Service) getMetricByDate(date string) (*entities.AnalyticsMetric, error) {
	m, err := s.metrics.FindOneBy(map[string]interface{}{"date": date})
	if err != nil {
		return nil, err
	}

	if m == nil {
		return nil, nil
	}

	return m.(*entities.AnalyticsMetric), nil
}
