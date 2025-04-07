package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTypeZero(t *testing.T) {
	ct := ParseType("date", 0, 0)
	assert.Equal(t, Scale_Type_Zero, ct.Type)
	assert.False(t, ct.IsInteger())
	assert.False(t, ct.IsNumeric())
	assert.Equal(t, "DATE", ct.Express())
}

func TestParseTypeOneNum(t *testing.T) {
	ct := ParseType("int", 10, 0)
	assert.Equal(t, Scale_Type_One, ct.Type)
	assert.True(t, ct.IsInteger())
	assert.True(t, ct.IsNumeric())
	assert.Equal(t, "INT(10)", ct.Express())
}

func TestParseTypeOne(t *testing.T) {
	ct := ParseType("varchar", 20, 0)
	assert.Equal(t, Scale_Type_One, ct.Type)
	assert.False(t, ct.IsInteger())
	assert.False(t, ct.IsNumeric())
	assert.Equal(t, "VARCHAR(20)", ct.Express())
}

func TestParseTypeTwo(t *testing.T) {
	ct := ParseType("double", 20, 3)
	assert.Equal(t, Scale_Type_Two, ct.Type)
	assert.False(t, ct.IsInteger())
	assert.True(t, ct.IsNumeric())
	assert.Equal(t, "DOUBLE(20,3)", ct.Express())
}

func TestParseTypeMore(t *testing.T) {
	ct := ParseType("enum", 0, 0, "A", "B", "C")
	assert.Equal(t, Scale_Type_More, ct.Type)
	assert.False(t, ct.IsInteger())
	assert.False(t, ct.IsNumeric())
	assert.Equal(t, "ENUM('A','B','C')", ct.Express())
}
