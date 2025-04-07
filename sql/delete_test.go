package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetele(t *testing.T) {
	d := NewDelete()
	w := NewWhere()
	w.Where("id", "=", 1)
	d.Table("user").Where(w)

	assert.Equal(t, "DELETE FROM `user` WHERE `id` = ?", d.Prepare())
	assert.Equal(t, []any{1}, d.Binds())
	assert.Equal(t, "DELETE FROM `user` WHERE `id` = ?", d.Prepare())
}
