package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForeignKeyConstraint(t *testing.T) {
	ch := &Changeset{}
	assert.Nil(t, ch.Constraints())

	ForeignKeyConstraint(ch, "field1")
	assert.Equal(t, 1, len(ch.Constraints()))
	assert.Equal(t, ForeignKeyConstraintKind, ch.Constraints()[0].Kind)
}
