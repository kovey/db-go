package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpatialReferenceSystemDrop(t *testing.T) {
	s := NewSpatialReferenceSystemDrop().IfExists().Srid("user")
	assert.Equal(t, "DROP SPATIAL REFERENCE SYSTEM IF EXISTS `user`", s.Prepare())
	assert.Nil(t, s.Binds())
}
