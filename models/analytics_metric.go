package models

import (
	"encoding/json"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/entities"
	"strings"
)

type AnalyticsMetricDocument struct {
	*entities.AnalyticsMetric
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

func (*AnalyticsMetricModel) NewEntity() content.Entity {
	return new(entities.AnalyticsMetric)
}

func (model *AnalyticsMetricModel) ToDocument(entity interface{}) DocumentInterface {
	return &AnalyticsMetricDocument{
		AnalyticsMetric: entity.(*entities.AnalyticsMetric),
	}
}
