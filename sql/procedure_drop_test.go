package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcedureDrop(t *testing.T) {
	s := NewProcedureDrop().IfExists().Procedure("user")
	assert.Equal(t, "DROP PROCEDURE IF EXISTS `user`", s.Prepare())
	assert.Nil(t, s.Binds())
}
