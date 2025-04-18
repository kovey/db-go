package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventDrop(t *testing.T) {
	s := NewEventDrop().IfExists().Event("user")
	assert.Equal(t, "DROP EVENT IF EXISTS `user`", s.Prepare())
	assert.Nil(t, s.Binds())
}
