package params_test

import (
	"net/url"
	"testing"

	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

func TestForm(t *testing.T) {
	tests := []struct {
		name   string
		values url.Values
		form   params.FormParams
	}{
		{
			name: "basic",
			values: url.Values{
				"name": []string{"lorem"},
				"tags": []string{"ipsum", "dolor"},
			},
			form: params.FormParams{
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
			form: params.FormParams{
				"contacts": []interface{}{params.FormParams{
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
			form: params.FormParams{
				"contacts": []interface{}{params.FormParams{
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
			form: params.FormParams{
				"a": []interface{}{params.FormParams{
					"b": []interface{}{params.FormParams{
						"c": []interface{}{params.FormParams{
							"d": []interface{}{params.FormParams{
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
			form: params.FormParams{
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
			form: params.FormParams{
				"addresses": []interface{}{
					params.FormParams{
						"city":     []interface{}{"lorem"},
						"province": []interface{}{"ipsum"},
						"tags":     []interface{}{"tag0", "tag1"},
						"refference": []interface{}{
							params.FormParams{"name": []interface{}{"refference"}},
						},
					},
					params.FormParams{
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
			assert.Equal(t, tt.form, params.Form(tt.values))
		})
	}
}
