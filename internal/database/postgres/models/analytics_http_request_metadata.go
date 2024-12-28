package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/analytics"
)

type AnalyticsHTTPRequestMetadataDocument struct {
	*analytics.AnalyticsHTTPRequestMetadata
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

func (*AnalyticsHTTPRequestMetadataModel) NewEntity() interface{} {
	return new(analytics.AnalyticsHTTPRequestMetadata)
}

func (model *AnalyticsHTTPRequestMetadataModel) ToDocument(entity interface{}) database.DocumentInterface {
	return &AnalyticsHTTPRequestMetadataDocument{
		AnalyticsHTTPRequestMetadata: entity.(*analytics.AnalyticsHTTPRequestMetadata),
	}
}
