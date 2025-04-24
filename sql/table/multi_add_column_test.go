package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiAddColumn(t *testing.T) {
	m := NewMultiAddColumn()
	m.Column("name", "varchar", 10, 0).Comment("name")
	m.Column("balance", "DECIMAL", 20, 2).Comment("balance").Default("0.00")
	m.Column("count", "bigint", 20, 0).Comment("balance").Default("0")
	var builder strings.Builder
	m.Build(&builder)
	assert.Equal(t, "ADD COLUMN (`name` VARCHAR(10) COMMENT 'name', `balance` DECIMAL(20,2) DEFAULT '0.00' COMMENT 'balance', `count` BIGINT(20) DEFAULT '0' COMMENT 'balance')", builder.String())
}
