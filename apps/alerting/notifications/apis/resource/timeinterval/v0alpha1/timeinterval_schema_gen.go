//
// Code generated by grafana-app-sdk. DO NOT EDIT.
//

package v0alpha1

import (
	"github.com/grafana/grafana-app-sdk/resource"
)

// schema is unexported to prevent accidental overwrites
var (
	schemaTimeInterval = resource.NewSimpleSchema("notifications.alerting.grafana.app", "v0alpha1", &TimeInterval{}, &TimeIntervalList{}, resource.WithKind("TimeInterval"),
		resource.WithPlural("timeintervals"), resource.WithScope(resource.NamespacedScope))
	kindTimeInterval = resource.Kind{
		Schema: schemaTimeInterval,
		Codecs: map[resource.KindEncoding]resource.Codec{
			resource.KindEncodingJSON: &JSONCodec{},
		},
	}
)

// Kind returns a resource.Kind for this Schema with a JSON codec
func Kind() resource.Kind {
	return kindTimeInterval
}

// Schema returns a resource.SimpleSchema representation of TimeInterval
func Schema() *resource.SimpleSchema {
	return schemaTimeInterval
}

// Interface compliance checks
var _ resource.Schema = kindTimeInterval
