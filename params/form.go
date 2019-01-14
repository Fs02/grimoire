package params

import (
	"net/url"
	"strconv"
	"strings"
)

// FormParams is param type alias for url.Values
type FormParams map[string][]interface{}

func fieldsExtractor(c rune) bool {
	return c == '[' || c == ']' || c == '.'
}

func assigns(form FormParams, pfield string, cfield string, index int, values interface{}) {
	if _, exist := form[pfield]; !exist && pfield != "" && cfield != "" {
		form[pfield] = []interface{}{FormParams{}}
	}

	if cfield != "" {
		if pfield != "" {
			form = form[pfield][0].(FormParams)
		}

		if stringValues, ok := values.([]string); ok {
			for _, value := range stringValues {
				form[cfield] = append(form[cfield], value)
			}
		}
	} else {
		// expand
		if index >= len(form[pfield]) {
			form[pfield] = append(form[pfield], make([]interface{}, index-len(form[pfield])+1)...)
		}

		if form[pfield][index] == nil {
			if values == nil {
				form[pfield][index] = FormParams{}
			} else {
				form[pfield][index] = values
			}
		}
	}
}

// Form parse form from url values.
func Form(raw url.Values) FormParams {
	result := make(FormParams, len(raw))

	for k, v := range raw {
		fields := strings.FieldsFunc(k, fieldsExtractor)

		pfield, cfield := "", ""
		form := result
		for i := range fields {
			pfield = cfield
			cfield = fields[i]

			if index, err := strconv.Atoi(cfield); err == nil {
				if i == len(fields)-1 {
					assigns(form, pfield, "", index, v[0])
				} else {
					assigns(form, pfield, "", index, nil)
					// if pfield != "" {
					form = form[pfield][index].(FormParams)
					// }
					cfield = "" // set cfield empty, so unnecesary nesting wont be created in the next loop
				}
			} else {
				if i == len(fields)-1 {
					assigns(form, pfield, cfield, -1, v)
				} else {
					assigns(form, pfield, cfield, -1, nil)
					if pfield != "" {
						index := len(form[pfield]) - 1
						form = form[pfield][index].(FormParams)
					}
				}
			}
		}
	}

	return result
}

// Exists returns true if key exists.
func (formParams FormParams) Exists(name string) bool {
	_, exists := formParams[name]
	return exists
}

// Get returns value as interface.
// returns nil if value doens't exists.
func (formParams FormParams) Get(name string) interface{} {
	return formParams[name]
}

// // GetWithType returns value given from given name and type.
// // second return value will only be false if the type of parameter is not convertible to requested type.
// // If value is not convertible to type, it'll return nil, false
// // If value is not exists, it will return nil, true
// func (m Map) GetWithType(name string, typ reflect.Type) (interface{}, bool) {
// 	value := m[name]

// 	if value == nil {
// 		return nil, true
// 	}

// 	rv := reflect.ValueOf(value)
// 	rt := rv.Type()

// 	if rt.Kind() == reflect.Ptr {
// 		rv = rv.Elem()
// 		rt = rt.Elem()
// 	}

// 	if !rv.IsValid() {
// 		return nil, true
// 	}

// 	if typ.Kind() == reflect.Slice && (rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array) {
// 		result := reflect.MakeSlice(typ, rv.Len(), rv.Len())
// 		elemTyp := typ.Elem()

// 		for i := 0; i < rv.Len(); i++ {
// 			elem := rv.Index(i)
// 			if elem.Kind() == reflect.Interface {
// 				elem = elem.Elem()
// 			}

// 			if elem.Type().ConvertibleTo(elemTyp) {
// 				result.Index(i).Set(elem.Convert(elemTyp))
// 			} else {
// 				return nil, false
// 			}
// 		}

// 		return result.Interface(), true
// 	}

// 	if !rt.ConvertibleTo(typ) {
// 		return nil, false
// 	}

// 	return rv.Convert(typ).Interface(), true
// }

// GetParams returns nested param
func (formParams FormParams) GetParams(name string) (Params, bool) {
	if val, exist := formParams[name]; exist && len(val) == 1 {
		if par, ok := val[0].(Params); ok {
			return par, ok
		}
	}

	return nil, false
}

// GetParamsSlice returns slice of nested param
func (formParams FormParams) GetParamsSlice(name string) ([]Params, bool) {
	if val, exist := formParams[name]; exist {
		pars := make([]Params, len(val))
		ok := true

		for i := range val {
			pars[i], ok = val[i].(Params)

			if !ok {
				return nil, false
			}
		}
	}

	return nil, false
}
