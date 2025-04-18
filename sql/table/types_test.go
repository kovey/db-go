package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTypeZero(t *testing.T) {
	ct := ParseType("date", 0, 0)
	assert.Equal(t, Scale_Type_Zero, ct.Type)
	assert.False(t, ct.IsInteger())
	assert.False(t, ct.IsNumeric())
	var builder strings.Builder
	ct.Build(&builder)
	assert.Equal(t, " DATE", builder.String())
}

func TestParseTypeOneNum(t *testing.T) {
	ct := ParseType("int", 10, 0)
	assert.Equal(t, Scale_Type_One, ct.Type)
	assert.True(t, ct.IsInteger())
	assert.True(t, ct.IsNumeric())
	var builder strings.Builder
	ct.Build(&builder)
	assert.Equal(t, " INT(10)", builder.String())
}

func TestParseTypeOne(t *testing.T) {
	ct := ParseType("varchar", 20, 0)
	assert.Equal(t, Scale_Type_One, ct.Type)
	assert.False(t, ct.IsInteger())
	assert.False(t, ct.IsNumeric())
	var builder strings.Builder
	ct.Build(&builder)
	assert.Equal(t, " VARCHAR(20)", builder.String())
}

func TestParseTypeTwo(t *testing.T) {
	ct := ParseType("double", 20, 3)
	assert.Equal(t, Scale_Type_Two, ct.Type)
	assert.False(t, ct.IsInteger())
	assert.True(t, ct.IsNumeric())
	var builder strings.Builder
	ct.Build(&builder)
	assert.Equal(t, " DOUBLE(20,3)", builder.String())
}

func TestParseTypeMore(t *testing.T) {
	ct := ParseType("enum", 0, 0, "A", "B", "C")
	assert.Equal(t, Scale_Type_More, ct.Type)
	assert.False(t, ct.IsInteger())
	assert.False(t, ct.IsNumeric())
	var builder strings.Builder
	ct.Build(&builder)
	assert.Equal(t, " ENUM('A','B','C')", builder.String())
}
