package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableDrop(t *testing.T) {
	d := NewDropTable().Table("user")
	assert.Equal(t, "DROP TABLE `user`", d.Prepare())
	assert.Nil(t, d.Binds())
	assert.Equal(t, "DROP TABLE `user`", d.Prepare())
}
