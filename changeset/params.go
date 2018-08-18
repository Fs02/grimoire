package changeset

import (
	"reflect"
)

// Params is interface used by changeset when casting parameters to changeset.
type Params interface {
	Exists(name string) bool
	Value(name string, typ reflect.Type) (interface{}, bool)
	Params(name string) (Params, bool)
	ParamsSlice(name string) ([]Params, bool)
}
