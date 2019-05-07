package sql

import (
	"reflect"

	"github.com/Fs02/grimoire/schema"
)

// Rows is minimal rows interface for test purpose
type Rows interface {
	Scan(dest ...interface{}) error
	Columns() ([]string, error)
	Next() bool
}

// Scan rows into interface
func Scan(value interface{}, rows Rows) (int64, error) {
	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("grimoire: record parameter must be a pointer")
	}
	rv = rv.Elem()

	var (
		count   = int64(0)
		isSlice = rv.Kind() == reflect.Slice
	)

	if isSlice {
		rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
	}

	for rows.Next() {
		var elem reflect.Value
		if isSlice {
			elem = reflect.New(rv.Type().Elem()).Elem()
		} else {
			elem = rv
		}

		ptr := schema.InferScanners(elem.Addr().Interface(), columns)

		err = rows.Scan(ptr...)
		if err != nil {
			return 0, err
		}

		count++

		if isSlice {
			rv.Set(reflect.Append(rv, elem))
		} else {
			break
		}
	}

	return count, nil
}
