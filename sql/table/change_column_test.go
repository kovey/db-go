package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeColumnFirst(t *testing.T) {
	c := NewChangeColumn().Old("user_balance").First()
	c.New("balance", "DECIMAL", 20, 2).Default("0").Comment("balance")
	var builder strings.Builder
	c.Build(&builder)
	assert.Equal(t, "CHANGE COLUMN `user_balance` `balance` DECIMAL(20,2) DEFAULT '0' COMMENT 'balance' FIRST", builder.String())
}

func TestChangeColumnAfter(t *testing.T) {
	c := NewChangeColumn().Old("user_balance").After("ba_test")
	c.New("balance", "DECIMAL", 20, 2).Default("0").Comment("balance")
	var builder strings.Builder
	c.Build(&builder)
	assert.Equal(t, "CHANGE COLUMN `user_balance` `balance` DECIMAL(20,2) DEFAULT '0' COMMENT 'balance' AFTER `ba_test`", builder.String())
}
