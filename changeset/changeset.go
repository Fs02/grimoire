// Package changeset used to cast and validate data before saving it to the database.
package changeset

import (
	"reflect"
)

// ConstraintKind defines contraint type.
type ConstraintKind int

const (
	// Invalid constraint
	Invalid ConstraintKind = iota
	// UniqueConstraintKind for unique constraint.
	UniqueConstraintKind
	// ForeignKeyConstraintKind for foreign key constraint.
	ForeignKeyConstraintKind
	// CheckConstraintKind for check constraint.
	CheckConstraintKind
)

// Constraint defines information to infer constraint error.
type Constraint struct {
	Name    string
	Message string
	Exact   bool
	Kind    ConstraintKind
}

// Changeset used to cast and validate data before saving it to the database.
type Changeset struct {
	errors      []error
	params      map[string]interface{}
	changes     map[string]interface{}
	values      map[string]interface{}
	types       map[string]reflect.Type
	constraints []Constraint
}

// Errors of changeset.
func (changeset *Changeset) Errors() []error {
	return changeset.errors
}

// Error of changeset, returns the first error if any.
func (changeset *Changeset) Error() error {
	if changeset.errors != nil {
		return changeset.errors[0]
	}
	return nil
}

// Changes of changeset.
func (changeset *Changeset) Changes() map[string]interface{} {
	return changeset.changes
}

// Values of changeset.
func (changeset *Changeset) Values() map[string]interface{} {
	return changeset.values
}

// Types of changeset.
func (changeset *Changeset) Types() map[string]reflect.Type {
	return changeset.types
}

// Constraints of changeset.
func (changeset *Changeset) Constraints() []Constraint {
	return changeset.constraints
}
