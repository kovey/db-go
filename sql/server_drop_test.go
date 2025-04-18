package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerDrop(t *testing.T) {
	s := NewServerDrop().IfExists().Server("user")
	assert.Equal(t, "DROP SERVER IF EXISTS `user`", s.Prepare())
	assert.Nil(t, s.Binds())
}
