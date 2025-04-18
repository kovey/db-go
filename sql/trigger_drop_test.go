package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTriggerDrop(t *testing.T) {
	s := NewTriggerDrop().Trigger("user").Schema("test").IfExists()
	assert.Equal(t, "DROP TRIGGER IF EXISTS `test`.`user`", s.Prepare())
	assert.Nil(t, s.Binds())
}
