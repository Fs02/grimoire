package params

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
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
	// TODO: differentiate when values is more than one, then it should be treated as slice
	// TODO: https://guides.rubyonrails.org/security.html#unsafe-query-generation
	return formParams[name]
}

// GetWithType returns value given from given name and type.
// second return value will only be false if the type of parameter is not convertible to requested type.
// If value is not convertible to type, it'll return nil, false
// If value is not exists, it will return nil, true
func (formParams FormParams) GetWithType(name string, typ reflect.Type) (interface{}, bool) {
	return nil, false
}

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

func (formParams FormParams) convert(str string, typ reflect.Type) (interface{}, bool) {
	result := interface{}(nil)
	valid := false

	switch typ.Kind() {
	case reflect.Bool:
		if parsed, err := strconv.ParseBool(str); err == nil {
			result, valid = parsed, true
		}
	case reflect.Int:
		if parsed, err := strconv.ParseInt(str, 10, 0); err == nil {
			result, valid = int(parsed), true
		}
	case reflect.Int8:
		if parsed, err := strconv.ParseInt(str, 10, 8); err == nil {
			result, valid = int8(parsed), true
		}
	case reflect.Int16:
		if parsed, err := strconv.ParseInt(str, 10, 16); err == nil {
			result, valid = int16(parsed), true
		}
	case reflect.Int32:
		if parsed, err := strconv.ParseInt(str, 10, 32); err == nil {
			result, valid = int32(parsed), true
		}
	case reflect.Int64:
		if parsed, err := strconv.ParseInt(str, 10, 64); err == nil {
			result, valid = int64(parsed), true
		}
	case reflect.Uint:
		if parsed, err := strconv.ParseUint(str, 10, 0); err == nil {
			result, valid = uint(parsed), true
		}
	case reflect.Uint8:
		if parsed, err := strconv.ParseUint(str, 10, 8); err == nil {
			result, valid = uint8(parsed), true
		}
	case reflect.Uint16:
		if parsed, err := strconv.ParseUint(str, 10, 16); err == nil {
			result, valid = uint16(parsed), true
		}
	case reflect.Uint32:
		if parsed, err := strconv.ParseUint(str, 10, 32); err == nil {
			result, valid = uint32(parsed), true
		}
	case reflect.Uint64:
		if parsed, err := strconv.ParseUint(str, 10, 64); err == nil {
			result, valid = uint64(parsed), true
		}
	case reflect.Float32:
		if parsed, err := strconv.ParseFloat(str, 32); err == nil {
			result, valid = float32(parsed), true
		}
	case reflect.Float64:
		if parsed, err := strconv.ParseFloat(str, 64); err == nil {
			result, valid = float64(parsed), true
		}
	case reflect.String:
		result, valid = str, true
	case reflect.Struct:
		if typ == timeType {
			if parsed, err := time.Parse(time.RFC3339, str); err == nil {
				result, valid = parsed, true
			}
		}
	}

	return result, valid
}
