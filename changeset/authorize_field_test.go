package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizeField(t *testing.T) {
	ch := &Changeset{
		values: map[string]interface{}{
			"field1": 1,
			"field2": 2,
			"field3": 3,
		},
		changes: map[string]interface{}{
			"field1": 4,
			"field2": 5,
			"field3": 6,
		},
	}

	AuthorizedField(ch, []string{"field1", "field4"}, true)
	assert.Nil(t, ch.Errors())

	AuthorizedField(ch, []string{"field1", "field2", "field3"}, false)
	assert.NotNil(t, ch.Errors())
}
