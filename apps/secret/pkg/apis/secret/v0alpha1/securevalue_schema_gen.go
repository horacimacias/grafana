//
// Code generated by grafana-app-sdk. DO NOT EDIT.
//

package v0alpha1

import (
	"github.com/grafana/grafana-app-sdk/resource"
)

// schema is unexported to prevent accidental overwrites
var (
	schemaSecureValue = resource.NewSimpleSchema("secret.grafana.app", "v0alpha1", &SecureValue{}, &SecureValueList{}, resource.WithKind("SecureValue"),
		resource.WithPlural("securevalues"), resource.WithScope(resource.NamespacedScope))
	kindSecureValue = resource.Kind{
		Schema: schemaSecureValue,
		Codecs: map[resource.KindEncoding]resource.Codec{
			resource.KindEncodingJSON: &SecureValueJSONCodec{},
		},
	}
)

// Kind returns a resource.Kind for this Schema with a JSON codec
func SecureValueKind() resource.Kind {
	return kindSecureValue
}

// Schema returns a resource.SimpleSchema representation of SecureValue
func SecureValueSchema() *resource.SimpleSchema {
	return schemaSecureValue
}

// Interface compliance checks
var _ resource.Schema = kindSecureValue
