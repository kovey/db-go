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

func TestTableDropTemporary(t *testing.T) {
	d := NewDropTable().Table("user").Temporary()
	assert.Equal(t, "DROP TEMPORARY TABLE `user`", d.Prepare())
	assert.Nil(t, d.Binds())
	assert.Equal(t, "DROP TEMPORARY TABLE `user`", d.Prepare())
}

func TestTableDropIfExists(t *testing.T) {
	d := NewDropTable().Table("user").IfExists().Table("users")
	assert.Equal(t, "DROP TABLE IF EXISTS `user`, `users`", d.Prepare())
	assert.Nil(t, d.Binds())
	assert.Equal(t, "DROP TABLE IF EXISTS `user`, `users`", d.Prepare())
}
