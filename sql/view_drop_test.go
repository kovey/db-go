package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewDrop(t *testing.T) {
	d := NewDropView().View("user")
	assert.Equal(t, "DROP VIEW `user`", d.Prepare())
	assert.Nil(t, d.Binds())
	assert.Equal(t, "DROP VIEW `user`", d.Prepare())
}

func TestViewDropIfExists(t *testing.T) {
	d := NewDropView().View("user").IfExists().View("users")
	assert.Equal(t, "DROP VIEW IF EXISTS `user`, `users`", d.Prepare())
	assert.Nil(t, d.Binds())
	assert.Equal(t, "DROP VIEW IF EXISTS `user`, `users`", d.Prepare())
}
