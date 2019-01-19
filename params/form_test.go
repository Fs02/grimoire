package params_test

import (
	"net/url"
	"testing"

	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

func TestForm_Exists(t *testing.T) {
	p := params.ParseForm(url.Values{
		"exists": []string{"true"},
	})

	assert.True(t, p.Exists("exists"))
	assert.False(t, p.Exists("not-exists"))
}

func TestForm_Get(t *testing.T) {
	p := params.ParseForm(url.Values{
		"exists": []string{"true"},
	})

	assert.Equal(t, []interface{}{"true"}, p.Get("exists"))
	assert.Equal(t, nil, p.Get("not-exists"))
}

func TestForm_GetParams(t *testing.T) {
	p := params.ParseForm(url.Values{
		"object[value]": []string{"true"},
		"array[0]":      []string{"0"},
		"value":         []string{"true"},
	})

	tests := []struct {
		name  string
		valid bool
	}{
		{
			name:  "object",
			valid: true,
		},
		{
			name:  "array",
			valid: false,
		},
		{
			name:  "value",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param, valid := p.GetParams(tt.name)
			assert.Equal(t, tt.valid, valid)
			assert.Equal(t, tt.valid, param != nil)
		})
	}
}

func TestForm_GetParamsSlice(t *testing.T) {
	p := params.ParseForm(url.Values{
		"array of object[0][value]": []string{"0"},
		"array of object[1][value]": []string{"1"},
		"array of value[0]":         []string{"true"},
		"array of value[1]":         []string{"false"},
		"array of mixed[0][value]":  []string{"0"},
		"array of mixed[1]":         []string{"true"},
		// "array[0]":                  []string{},
		// "object[value]":             []string{},
		"value": []string{"true"},
	})

	tests := []struct {
		name  string
		valid bool
	}{
		{
			name:  "array of object",
			valid: true,
		},
		{
			name:  "array of array",
			valid: false,
		},
		{
			name:  "array of value",
			valid: false,
		},
		{
			name:  "array of mixed",
			valid: false,
		},
		// {
		// 	name:  "array",
		// 	valid: true,
		// },
		{
			name:  "object",
			valid: false,
		},
		{
			name:  "value",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, valid := p.GetParamsSlice(tt.name)
			assert.Equal(t, tt.valid, valid)
			assert.Equal(t, tt.valid, params != nil)
		})
	}
}

func TestParseForm(t *testing.T) {
	tests := []struct {
		name   string
		values url.Values
		form   params.Form
	}{
		{
			name: "basic",
			values: url.Values{
				"name": []string{"lorem"},
				"tags": []string{"ipsum", "dolor"},
			},
			form: params.Form{
				"name": []interface{}{"lorem"},
				"tags": []interface{}{"ipsum", "dolor"},
			},
		},
		{
			name: "nested",
			values: url.Values{
				"contacts[phone]": []string{"+628123123123"},
				"contacts[email]": []string{"lorem@ipsum.dolor"},
			},
			form: params.Form{
				"contacts": []interface{}{params.Form{
					"phone": []interface{}{"+628123123123"},
					"email": []interface{}{"lorem@ipsum.dolor"},
				}},
			},
		},
		{
			name: "nested using dot",
			values: url.Values{
				"contacts.phone": []string{"+628123123123"},
				"contacts.email": []string{"lorem@ipsum.dolor"},
			},
			form: params.Form{
				"contacts": []interface{}{params.Form{
					"phone": []interface{}{"+628123123123"},
					"email": []interface{}{"lorem@ipsum.dolor"},
				}},
			},
		},
		{
			name: "deeply nested",
			values: url.Values{
				"a[b][c][d][e]": []string{"lorem"},
			},
			form: params.Form{
				"a": []interface{}{params.Form{
					"b": []interface{}{params.Form{
						"c": []interface{}{params.Form{
							"d": []interface{}{params.Form{
								"e": []interface{}{"lorem"},
							}},
						}},
					}},
				}},
			},
		},
		{
			name: "slice",
			values: url.Values{
				"tags[0]": []string{"lorem"},
				"tags[1]": []string{"ipsum"},
				"tags[3]": []string{"three"},
			},
			form: params.Form{
				"tags": []interface{}{"lorem", "ipsum", nil, "three"},
			},
		},
		{
			name: "slice of params",
			values: url.Values{
				"addresses[0][city]":                []string{"lorem"},
				"addresses[0][province]":            []string{"ipsum"},
				"addresses[0][tags][0]":             []string{"tag0"},
				"addresses[0][tags][1]":             []string{"tag1"},
				"addresses[0][refference][0][name]": []string{"refference"},
				"addresses[1][city]":                []string{"dolor"},
				"addresses[1][province]":            []string{"sit"},
				"addresses[1][tags][0]":             []string{"tag0"},
				"addresses[1][tags][1]":             []string{"tag1"},
			},
			form: params.Form{
				"addresses": []interface{}{
					params.Form{
						"city":     []interface{}{"lorem"},
						"province": []interface{}{"ipsum"},
						"tags":     []interface{}{"tag0", "tag1"},
						"refference": []interface{}{
							params.Form{"name": []interface{}{"refference"}},
						},
					},
					params.Form{
						"city":     []interface{}{"dolor"},
						"province": []interface{}{"sit"},
						"tags":     []interface{}{"tag0", "tag1"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.form, params.ParseForm(tt.values))
		})
	}
}
