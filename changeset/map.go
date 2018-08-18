package changeset

import (
	"reflect"
)

// Map is param type alias for map[string]interface{}
type Map map[string]interface{}

// Exists returns true if key exists.
func (m Map) Exists(name string) bool {
	_, exists := m[name]
	return exists
}

// Value returns value given from given name and type.
// If value is not convertible to type, it'll return nil, false
// If value is not exists, it will return nil, true
func (m Map) Value(name string, typ reflect.Type) (interface{}, bool) {
	value := m[name]

	if value == nil {
		return reflect.Zero(reflect.PtrTo(typ)).Interface(), true
	}

	rv := reflect.ValueOf(value)
	rt := rv.Type()

	if rt.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if !rv.IsValid() {
		return reflect.Zero(reflect.PtrTo(typ)).Interface(), true
	}

	if (typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array) && (rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array) && rt.Elem().Kind() == reflect.Interface {
		result := reflect.Zero(typ)
		elemTyp := typ.Elem()

		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i)
			if elem.Elem().Type().ConvertibleTo(elemTyp) {
				result = reflect.Append(result, elem.Elem().Convert(elemTyp))
			} else {
				return reflect.Zero(typ).Interface(), false
			}
		}

		return result.Interface(), true
	}

	if !rt.ConvertibleTo(typ) {
		return reflect.Zero(reflect.PtrTo(typ)).Interface(), false
	}

	return rv.Convert(typ).Interface(), true
}

// Params returns nested param
func (m Map) Params(name string) (Params, bool) {
	if val, exist := m[name]; exist {
		if par, ok := val.(Params); ok {
			return par, ok
		}

		if par, ok := val.(map[string]interface{}); ok {
			return Map(par), ok
		}
	}

	return nil, false
}

// ParamsSlice returns slice of nested param
func (m Map) ParamsSlice(name string) ([]Params, bool) {
	if val, exist := m[name]; exist {
		if pars, ok := val.([]Params); ok {
			return pars, ok
		}

		if pars, ok := val.([]Map); ok {
			mpar := make([]Params, len(pars))
			for i, par := range pars {
				mpar[i] = Map(par)
			}
			return mpar, ok
		}

		if pars, ok := val.([]map[string]interface{}); ok {
			mpar := make([]Params, len(pars))
			for i, par := range pars {
				mpar[i] = Map(par)
			}
			return mpar, ok
		}
	}

	return nil, false
}
