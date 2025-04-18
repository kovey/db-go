package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnTypeZero(t *testing.T) {
	var builder strings.Builder
	c := &ColumnType{Name: Type_Int, Type: Scale_Type_Zero}
	c.Build(&builder)
	assert.Equal(t, " INT", builder.String())
	assert.True(t, c.IsInteger())
	assert.True(t, c.IsNumeric())
}

func TestColumnTypeZeroFloat(t *testing.T) {
	var builder strings.Builder
	c := &ColumnType{Name: Type_Float, Type: Scale_Type_Zero}
	c.Build(&builder)
	assert.Equal(t, " FLOAT", builder.String())
	assert.False(t, c.IsInteger())
	assert.True(t, c.IsNumeric())
}

func TestColumnTypeOne(t *testing.T) {
	var builder strings.Builder
	c := &ColumnType{Name: Type_VarChar, Type: Scale_Type_One, Length: 20}
	c.Build(&builder)
	assert.Equal(t, " VARCHAR(20)", builder.String())
	assert.False(t, c.IsInteger())
	assert.False(t, c.IsNumeric())
}

func TestColumnTypeOneBigInt(t *testing.T) {
	var builder strings.Builder
	c := &ColumnType{Name: Type_BigInt, Type: Scale_Type_One, Length: 20}
	c.Build(&builder)
	assert.Equal(t, " BIGINT(20)", builder.String())
	assert.True(t, c.IsInteger())
	assert.True(t, c.IsNumeric())
}

func TestColumnTypeTwo(t *testing.T) {
	var builder strings.Builder
	c := &ColumnType{Name: Type_Decimal, Type: Scale_Type_Two, Length: 20, Scale: 2}
	c.Build(&builder)
	assert.Equal(t, " DECIMAL(20,2)", builder.String())
	assert.False(t, c.IsInteger())
	assert.True(t, c.IsNumeric())
}

func TestColumnTypeMore(t *testing.T) {
	var builder strings.Builder
	c := &ColumnType{Name: Type_Set, Type: Scale_Type_More, Length: 20, Scale: 2}
	c.Set("1", "3", "5")
	c.Build(&builder)
	assert.Equal(t, " SET('1','3','5')", builder.String())
	assert.False(t, c.IsInteger())
	assert.False(t, c.IsNumeric())
}
