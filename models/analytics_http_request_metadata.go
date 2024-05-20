package models

import (
	"encoding/json"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/entities"
	"strings"
)

type AnalyticsHTTPRequestMetadataDocument struct {
	*entities.AnalyticsHTTPRequestMetadata
}

func (document *AnalyticsHTTPRequestMetadataDocument) Value() (interface{}, error) {
	return json.Marshal(document)
}

func (document *AnalyticsHTTPRequestMetadataDocument) Scan(src interface{}) error {
	if byteSrc, ok := src.([]byte); ok {
		return json.Unmarshal(byteSrc, &document)
	}

	if stringSrc, ok := src.(string); ok {
		return json.NewDecoder(strings.NewReader(stringSrc)).Decode(&document)
	}

	return fmt.Errorf("unsupported type %T", src)
}

type AnalyticsHTTPRequestMetadataModel struct{}

func (*AnalyticsHTTPRequestMetadataModel) Name() string {
	return WrapPonzuModelNameSpace("analytics_http_request_metadata")
}

func (*AnalyticsHTTPRequestMetadataModel) NewEntity() content.Entity {
	return new(entities.AnalyticsHTTPRequestMetadata)
}

func (model *AnalyticsHTTPRequestMetadataModel) ToDocument(entity interface{}) DocumentInterface {
	return &AnalyticsHTTPRequestMetadataDocument{
		AnalyticsHTTPRequestMetadata: entity.(*entities.AnalyticsHTTPRequestMetadata),
	}
}
