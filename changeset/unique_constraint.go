package changeset

import (
	"strings"
)

// UniqueConstraintMessage is the default error message for UniqueConstraint.
var UniqueConstraintMessage = "{field} has already been taken"

// UniqueConstraint adds an unique constraint to changeset.
func UniqueConstraint(ch *Changeset, field string, opts ...Option) {
	options := Options{
		message: UniqueConstraintMessage,
		name:    field,
		exact:   false,
	}
	options.apply(opts)

	ch.constraints = append(ch.constraints, Constraint{
		Name:    options.name,
		Message: strings.Replace(options.message, "{field}", field, 1),
		Exact:   options.exact,
		Kind:    UniqueConstraintKind,
	})
}
