package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueConstraint(t *testing.T) {
	ch := &Changeset{}
	assert.Nil(t, ch.Constraints())

	UniqueConstraint(ch, "field1")
	assert.Equal(t, 1, len(ch.Constraints()))
	assert.Equal(t, UniqueConstraintKind, ch.Constraints()[0].Kind)
}
