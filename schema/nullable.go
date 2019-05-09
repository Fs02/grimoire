package schema

import (
	"database/sql"
	"reflect"
)

type nullable struct {
	dest interface{}
}

var _ sql.Scanner = (*nullable)(nil)

func (n nullable) Scan(src interface{}) error {
	return convertAssign(n.dest, src)
}

func Nullable(dest interface{}) interface{} {
	if s, ok := dest.(sql.Scanner); ok {
		return s
	}

	rt := reflect.TypeOf(dest)
	if rt.Kind() != reflect.Ptr {
		panic("grimoire: destination must be a pointer")
	}

	if rt.Elem().Kind() == reflect.Ptr {
		return dest
	}

	return nullable{
		dest: dest,
	}
}
