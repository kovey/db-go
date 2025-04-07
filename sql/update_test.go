package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	u := NewUpdate()
	w := NewWhere()
	w.Where("id", "=", 1)
	u.Table("user").Set("name", "kovey").Set("age", 10).Where(w)

	assert.Equal(t, "UPDATE `user` SET `name` = ?,`age` = ? WHERE `id` = ?", u.Prepare())
	assert.Equal(t, []any{"kovey", 10, 1}, u.Binds())
	assert.Equal(t, "UPDATE `user` SET `name` = ?,`age` = ? WHERE `id` = ?", u.Prepare())
}
