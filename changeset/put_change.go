package changeset

import (
	"reflect"
	"strings"
)

// PutChangeErrorMessage is the default error message for PutChange.
var PutChangeErrorMessage = "{field} is invalid"

// PutChange to changeset.
func PutChange(ch *Changeset, field string, value interface{}, opts ...Option) {
	options := Options{
		message: PutChangeErrorMessage,
	}
	options.apply(opts)

	valTyp := reflect.TypeOf(value)
	if valTyp.Kind() == reflect.Ptr {
		valTyp = valTyp.Elem()
	}

	if typ, exist := ch.types[field]; exist && valTyp.ConvertibleTo(typ) {
		ch.changes[field] = value
	} else {
		msg := strings.Replace(options.message, "{field}", field, 1)
		AddError(ch, field, msg)
	}
}
