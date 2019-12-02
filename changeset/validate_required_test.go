package changeset

import (
	"fmt"
	"testing"
	"time"

	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequired(t *testing.T) {
	ch := &Changeset{
		values: map[string]interface{}{
			"field3": true,
		},
		changes: map[string]interface{}{
			"field1": 1,
			"field2": " 1 ",
		},
	}

	ValidateRequired(ch, []string{"field1", "field2", "field3"})
	assert.Nil(t, ch.Errors())
}

func TestValidateRequired_cast(t *testing.T) {
	type customString string
	type customType struct {
		Field1    customString
		Field2    customString
		Field3    customString
		FieldTime time.Time
	}

	ct := customType{Field1: "field1", FieldTime: time.Now()}
	ch := Cast(ct, params.Map{"field2": "field2", "field3": "field3"}, []string{"field1", "field2", "field3", "field_time"})
	ValidateRequired(ch, []string{"field1", "field2", "field3", "field_time"})

	assert.Nil(t, ch.Errors())
}

func TestValidateRequired_error(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field1": nil,
			"field2": "  ",
		},
	}

	ValidateRequired(ch, []string{"field1", "field2", "field3"})
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, 3, len(ch.Errors()))
	assert.Equal(t, "field1 is required", ch.Errors()[0].Error())
	assert.Equal(t, "field2 is required", ch.Errors()[1].Error())
	assert.Equal(t, "field3 is required", ch.Errors()[2].Error())

	// empty struct
	ch = &Changeset{
		zero: true,
		values: map[string]interface{}{
			"field1": 0,
			"field2": "",
			"field3": false,
		},
		changes: map[string]interface{}{
			"field1": 1,
		},
	}
	ValidateRequired(ch, []string{"field1", "field2", "field3"})
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, 2, len(ch.Errors()))
	assert.Equal(t, "field2 is required", ch.Errors()[0].Error())
	assert.Equal(t, "field3 is required", ch.Errors()[1].Error())
}

func TestValidateRequired_cast_error(t *testing.T) {
	type customString string
	type customType struct {
		Field1    customString
		Field2    customString
		Field3    customString
		FieldTime time.Time
	}

	ct := customType{Field1: "field1"}
	fmt.Println(ct.FieldTime)
	ch := Cast(ct, params.Map{"field2": "field2"}, []string{"field1", "field2", "field3", "field_time"})
	ValidateRequired(ch, []string{"field1", "field2", "field3", "field_time"})

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, 2, len(ch.Errors()))
	assert.Equal(t, "field3 is required", ch.Errors()[0].Error())
	assert.Equal(t, "field_time is required", ch.Errors()[1].Error())
}
