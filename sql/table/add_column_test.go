package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddColumnFirst(t *testing.T) {
	c := NewAddColumn()
	c.Column("name", "varchar", 10, 0).Default("").Comment("name")
	c.Column("bad", "bad", 10, 0)
	c.First()
	var builder strings.Builder
	c.Build(&builder)
	assert.Equal(t, "ADD COLUMN `name` VARCHAR(10) DEFAULT '' COMMENT 'name' FIRST", builder.String())
}

func TestAddColumnAfter(t *testing.T) {
	c := NewAddColumn()
	c.Column("name", "varchar", 10, 0).Default("").Comment("name")
	c.After("user_id")
	var builder strings.Builder
	c.Build(&builder)
	assert.Equal(t, "ADD COLUMN `name` VARCHAR(10) DEFAULT '' COMMENT 'name' AFTER `user_id`", builder.String())
}
