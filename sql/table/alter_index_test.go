package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlterIndexInvisble(t *testing.T) {
	a := NewAlterIndex()
	a.Index("idx_name").Invisible()
	var builder strings.Builder
	a.Build(&builder)
	assert.Equal(t, "ALTER INDEX `idx_name` INVISIBLE", builder.String())
}

func TestAlterIndexVisble(t *testing.T) {
	a := NewAlterIndex()
	a.Index("idx_name").Visible()
	var builder strings.Builder
	a.Build(&builder)
	assert.Equal(t, "ALTER INDEX `idx_name` VISIBLE", builder.String())
}
