package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/analytics"
)

type AnalyticsMetricDocument struct {
	*analytics.Metric
}

func (document *AnalyticsMetricDocument) Value() (interface{}, error) {
	return json.Marshal(document)
}

func (document *AnalyticsMetricDocument) Scan(src interface{}) error {
	if byteSrc, ok := src.([]byte); ok {
		return json.Unmarshal(byteSrc, &document)
	}

	if stringSrc, ok := src.(string); ok {
		return json.NewDecoder(strings.NewReader(stringSrc)).Decode(&document)
	}

	return fmt.Errorf("unsupported type %T", src)
}

type AnalyticsMetricModel struct{}

func (*AnalyticsMetricModel) Name() string {
	return WrapPonzuModelNameSpace("analytics_metrics")
}

func (*AnalyticsMetricModel) NewEntity() interface{} {
	return new(analytics.Metric)
}

func (model *AnalyticsMetricModel) ToDocument(entity interface{}) database.DocumentInterface {
	return &AnalyticsMetricDocument{
		Metric: entity.(*analytics.Metric),
	}
}
