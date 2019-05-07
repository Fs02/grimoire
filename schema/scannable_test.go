package schema

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScannable(t *testing.T) {
	var (
		vint    = 0
		vstring = ""
		vbool   = false
		vtime   = time.Now()
		vbytes  = []byte{}
		vstruct = struct{}{}
	)

	assert.True(t, scannable(reflect.TypeOf(vint)))
	assert.True(t, scannable(reflect.TypeOf(vstring)))
	assert.True(t, scannable(reflect.TypeOf(vbool)))
	assert.True(t, scannable(reflect.TypeOf(vtime)))
	assert.True(t, scannable(reflect.TypeOf(vbytes)))
	assert.False(t, scannable(reflect.TypeOf(vstruct)))

	assert.True(t, scannable(reflect.TypeOf(&vint)))
	assert.True(t, scannable(reflect.TypeOf(&vstring)))
	assert.True(t, scannable(reflect.TypeOf(&vbool)))
	assert.True(t, scannable(reflect.TypeOf(&vtime)))
	assert.True(t, scannable(reflect.TypeOf(&vbytes)))
	assert.False(t, scannable(reflect.TypeOf(&vstruct)))
}
