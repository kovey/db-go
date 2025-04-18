package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaDrop(t *testing.T) {
	s := NewSchemaDrop().IfExists().Schema("user")
	assert.Equal(t, "DROP SCHEMA IF EXISTS `user`", s.Prepare())
	assert.Nil(t, s.Binds())
}
