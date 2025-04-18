package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T) {
	c := NewCall()
	c.Call("func").Params("@input1", "out1", "input2")
	assert.Equal(t, "CALL `func` (@input1, @out1, @input2)", c.Prepare())
	assert.Nil(t, c.Binds())
}
