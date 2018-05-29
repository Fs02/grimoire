package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckConstraint(t *testing.T) {
	ch := &Changeset{}
	assert.Nil(t, ch.Constraints())

	CheckConstraint(ch, "field1")
	assert.Equal(t, 1, len(ch.Constraints()))
	assert.Equal(t, CheckConstraintKind, ch.Constraints()[0].Kind)
}
