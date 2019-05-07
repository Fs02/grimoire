package schema

import (
	"database/sql"
	"reflect"
)

var tScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

func scannable(rt reflect.Type) bool {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	kind := rt.Kind()

	if (kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Array) &&
		kind != reflect.Uint8 && rt != Time && rt != Bytes && !rt.Implements(tScanner) {
		return false
	}

	return true
}
