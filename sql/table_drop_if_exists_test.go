package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableDropIfExists(t *testing.T) {
	d := NewDropTableIfExists().Table("user")
	assert.Equal(t, "DROP TABLE IF EXISTS `user`", d.Prepare())
	assert.Nil(t, d.Binds())
	assert.Equal(t, "DROP TABLE IF EXISTS `user`", d.Prepare())
}
