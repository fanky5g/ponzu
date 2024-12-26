package analytics

import (
	"github.com/fanky5g/ponzu/internal/analytics"
)

func (s *service) getMetrics() (map[string]analytics.AnalyticsMetric, error) {
	metrics, err := s.metricsRepository.FindAll()
	if err != nil {
		return nil, err
	}

	analyticsMetrics := make(map[string]analytics.AnalyticsMetric)
	for _, m := range metrics {
		metric := m.(*analytics.AnalyticsMetric)
		// add metric to currentMetrics map
		analyticsMetrics[metric.Date] = *metric
	}

	return analyticsMetrics, nil
}

func (s *service) getMetricByDate(date string) (*analytics.AnalyticsMetric, error) {
	m, err := s.metricsRepository.FindOneBy(map[string]interface{}{"date": date})
	if err != nil {
		return nil, err
	}

	if m == nil {
		return nil, nil
	}

	return m.(*analytics.AnalyticsMetric), nil
}
