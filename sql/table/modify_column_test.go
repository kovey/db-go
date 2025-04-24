package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModifyColumnFirst(t *testing.T) {
	m := NewModifyColumn("user_name")
	m.First().Column("varchar", 10, 0).Nullable()
	var builder strings.Builder
	m.Build(&builder)
	assert.Equal(t, "MODIFY COLUMN `user_name` VARCHAR(10) NULL FIRST", builder.String())
}

func TestModifyColumnAfter(t *testing.T) {
	m := NewModifyColumn("user_name")
	m.After("user_id").Column("varchar", 10, 0).Nullable()
	var builder strings.Builder
	m.Build(&builder)
	assert.Equal(t, "MODIFY COLUMN `user_name` VARCHAR(10) NULL AFTER `user_id`", builder.String())
}
