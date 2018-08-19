package changeset

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/Fs02/grimoire/errors"
	"github.com/Fs02/grimoire/params"
)

// CastAssocErrorMessage is the default error message for CastAssoc.
var CastAssocErrorMessage = "{field} is invalid"

// ChangeFunc is changeset function.
type ChangeFunc func(interface{}, params.Params) *Changeset

// CastAssoc casts association changes using changeset function.
// Repo insert or update won't persist any changes generated by CastAssoc.
func CastAssoc(ch *Changeset, field string, fn ChangeFunc, opts ...Option) {
	options := Options{
		message: CastAssocErrorMessage,
	}
	options.apply(opts)

	typ, texist := ch.types[field]
	valid := true

	if texist && ch.params.Exists(field) {
		if typ.Kind() == reflect.Struct {
			valid = castOne(ch, field, fn)
		} else if typ.Kind() == reflect.Slice && typ.Elem().Kind() == reflect.Struct {
			valid = castMany(ch, field, fn)
		}
	}

	if !valid {
		msg := strings.Replace(options.message, "{field}", field, 1)
		AddError(ch, field, msg)
	}
}

func castOne(ch *Changeset, field string, fn ChangeFunc) bool {
	par, valid := ch.params.GetParams(field)
	if !valid {
		return false
	}

	var innerch *Changeset

	if val, exist := ch.values[field]; exist && val != nil {
		innerch = fn(val, par)
	} else {
		innerch = fn(reflect.Zero(ch.types[field]).Interface(), par)
	}

	ch.changes[field] = innerch

	// add errors to main errors
	mergeErrors(ch, innerch, field+".")

	return true
}

func castMany(ch *Changeset, field string, fn ChangeFunc) bool {
	spar, valid := ch.params.GetParamsSlice(field)
	if !valid {
		return false
	}

	data := reflect.Zero(ch.types[field].Elem()).Interface()

	chs := make([]*Changeset, len(spar))
	for i, par := range spar {
		innerch := fn(data, par)
		chs[i] = innerch

		// add errors to main errors
		mergeErrors(ch, innerch, field+"["+strconv.Itoa(i)+"].")
	}
	ch.changes[field] = chs

	return true
}

func mergeErrors(parent *Changeset, child *Changeset, prefix string) {
	for _, err := range child.errors {
		e := err.(errors.Error)
		AddError(parent, prefix+e.Field, e.Message)
	}
}
