package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableTruncate(t *testing.T) {
	r := NewTableTruncate().Table("user")
	assert.Equal(t, "TRUNCATE TABLE `user`", r.Prepare())
	assert.Nil(t, r.Binds())
}
