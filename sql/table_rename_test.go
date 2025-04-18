package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableRename(t *testing.T) {
	r := NewTableRename().Table("user", "users").Table("ext", "user_ext")
	assert.Equal(t, "RENAME TABLE `user` TO `users`, `ext` TO `user_ext`", r.Prepare())
	assert.Nil(t, r.Binds())
}
