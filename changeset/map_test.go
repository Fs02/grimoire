package changeset_test

import (
	"reflect"
	"testing"

	"github.com/Fs02/grimoire/changeset"
	"github.com/stretchr/testify/assert"
)

func TestMap_Exists(t *testing.T) {
	p := changeset.Map{"exists": true}
	assert.True(t, p.Exists("exists"))
	assert.False(t, p.Exists("not-exists"))
}

func TestMap_Value(t *testing.T) {
	p := changeset.Map{
		"nil":                      (*bool)(nil),
		"incorrect type":           "some string",
		"correct type":             true,
		"slice":                    []bool{true, false},
		"slice of interface":       []interface{}{true, false},
		"slice of interface mixed": []interface{}{true, 0},
	}

	tests := []struct {
		name  string
		typ   reflect.Type
		value interface{}
		valid bool
	}{
		{
			name:  "not exist",
			typ:   reflect.TypeOf(true),
			value: (*bool)(nil),
			valid: true,
		},
		{
			name:  "nil",
			typ:   reflect.TypeOf(true),
			value: (*bool)(nil),
			valid: true,
		},
		{
			name:  "incorrect type",
			typ:   reflect.TypeOf(true),
			value: (*bool)(nil),
			valid: false,
		},
		{
			name:  "correct type",
			typ:   reflect.TypeOf(true),
			value: true,
			valid: true,
		},
		{
			name:  "slice",
			typ:   reflect.TypeOf([]bool{}),
			value: []bool{true, false},
			valid: true,
		},
		{
			name:  "slice of interface",
			typ:   reflect.TypeOf([]bool{}),
			value: []bool{true, false},
			valid: true,
		},
		{
			name:  "slice of interface mixed",
			typ:   reflect.TypeOf([]bool{}),
			value: []bool(nil),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, valid := p.Value(tt.name, tt.typ)
			assert.Equal(t, tt.value, value)
			assert.Equal(t, tt.valid, valid)
		})
	}
}

func TestMap_Param(t *testing.T) {
	p := changeset.Map{
		"changeset.Map":       changeset.Map{},
		"changeset.Map slice": []changeset.Map{},
		"map":       map[string]interface{}{},
		"map slice": []map[string]interface{}{},
		"invalid":   true,
	}

	tests := []struct {
		name  string
		param changeset.Params
		valid bool
	}{
		{
			name:  "changeset.Map",
			param: changeset.Map{},
			valid: true,
		},
		{
			name:  "changeset.Map slice",
			param: nil,
			valid: false,
		},
		{
			name:  "map",
			param: changeset.Map{},
			valid: true,
		},
		{
			name:  "map slice",
			param: nil,
			valid: false,
		},
		{
			name:  "invalid",
			param: nil,
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param, valid := p.Params(tt.name)
			assert.Equal(t, tt.param, param)
			assert.Equal(t, tt.valid, valid)
		})
	}
}

func TestMap_Params(t *testing.T) {
	p := changeset.Map{
		"changeset.Params slice": []changeset.Params{changeset.Map{}},
		"changeset.Map":          changeset.Map{},
		"changeset.Map slice":    []changeset.Map{changeset.Map{}},
		"map":       map[string]interface{}{},
		"map slice": []map[string]interface{}{map[string]interface{}{}},
		"invalid":   true,
	}

	tests := []struct {
		name   string
		params []changeset.Params
		valid  bool
	}{
		{
			name:   "changeset.Params slice",
			params: []changeset.Params{changeset.Map{}},
			valid:  true,
		},
		{
			name:   "changeset.Map",
			params: nil,
			valid:  false,
		},
		{
			name:   "changeset.Map slice",
			params: []changeset.Params{changeset.Map{}},
			valid:  true,
		},
		{
			name:   "map",
			params: nil,
			valid:  false,
		},
		{
			name:   "map slice",
			params: []changeset.Params{changeset.Map{}},
			valid:  true,
		},
		{
			name:   "invalid",
			params: nil,
			valid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, valid := p.ParamsSlice(tt.name)
			assert.Equal(t, tt.params, params)
			assert.Equal(t, tt.valid, valid)
		})
	}
}
