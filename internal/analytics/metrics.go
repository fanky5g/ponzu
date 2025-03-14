package analytics

func (s *Service) getMetrics() (map[string]Metric, error) {
	metrics, err := s.metricsRepository.FindAll()
	if err != nil {
		return nil, err
	}

	analyticsMetrics := make(map[string]Metric)
	for _, m := range metrics {
		metric := m.(*Metric)
		// add metric to currentMetrics map
		analyticsMetrics[metric.Date] = *metric
	}

	return analyticsMetrics, nil
}

func (s *Service) getMetricByDate(date string) (*Metric, error) {
	m, err := s.metricsRepository.FindOneBy(map[string]interface{}{"date": date})
	if err != nil {
		return nil, err
	}

	if m == nil {
		return nil, nil
	}

	return m.(*Metric), nil
}
