package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTablespaceDrop(t *testing.T) {
	s := NewTablespaceDrop().Tablespace("user").Engine("users").Undo()
	assert.Equal(t, "DROP UNDO TABLESPACE `user` ENGINE = users", s.Prepare())
	assert.Nil(t, s.Binds())
}
