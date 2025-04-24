package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlterColumnDefault(t *testing.T) {
	a := NewAlterColumn().Column("user_name")
	a.Default("").Invisible()
	var builder strings.Builder
	a.Build(&builder)
	assert.Equal(t, "ALTER COLUMN `user_name` SET DEFAULT '' SET INVISIBLE", builder.String())
}

func TestAlterColumnDefaultExpress(t *testing.T) {
	a := NewAlterColumn().Column("user_name")
	a.DefaultExpress("round(balance, 2)").Visible()
	var builder strings.Builder
	a.Build(&builder)
	assert.Equal(t, "ALTER COLUMN `user_name` SET DEFAULT (round(balance, 2)) SET VISIBLE", builder.String())
}

func TestAlterColumnDropDefault(t *testing.T) {
	a := NewAlterColumn().Column("user_name")
	a.DropDefault()
	var builder strings.Builder
	a.Build(&builder)
	assert.Equal(t, "ALTER COLUMN `user_name` DROP DEFAULT", builder.String())
}
