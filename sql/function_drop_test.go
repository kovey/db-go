package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionDrop(t *testing.T) {
	s := NewFunctionDrop().IfExists().Function("user")
	assert.Equal(t, "DROP FUNCTION IF EXISTS `user`", s.Prepare())
	assert.Nil(t, s.Binds())
}
