package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnTypeZero(t *testing.T) {
	c := &ColumnType{Name: Type_Int, Type: Scale_Type_Zero}
	assert.Equal(t, "INT", c.Express())
	assert.True(t, c.IsInteger())
	assert.True(t, c.IsNumeric())
}

func TestColumnTypeZeroFloat(t *testing.T) {
	c := &ColumnType{Name: Type_Float, Type: Scale_Type_Zero}
	assert.Equal(t, "FLOAT", c.Express())
	assert.False(t, c.IsInteger())
	assert.True(t, c.IsNumeric())
}

func TestColumnTypeOne(t *testing.T) {
	c := &ColumnType{Name: Type_VarChar, Type: Scale_Type_One, Length: 20}
	assert.Equal(t, "VARCHAR(20)", c.Express())
	assert.False(t, c.IsInteger())
	assert.False(t, c.IsNumeric())
}

func TestColumnTypeOneBigInt(t *testing.T) {
	c := &ColumnType{Name: Type_BigInt, Type: Scale_Type_One, Length: 20}
	assert.Equal(t, "BIGINT(20)", c.Express())
	assert.True(t, c.IsInteger())
	assert.True(t, c.IsNumeric())
}

func TestColumnTypeTwo(t *testing.T) {
	c := &ColumnType{Name: Type_Decimal, Type: Scale_Type_Two, Length: 20, Scale: 2}
	assert.Equal(t, "DECIMAL(20,2)", c.Express())
	assert.False(t, c.IsInteger())
	assert.True(t, c.IsNumeric())
}

func TestColumnTypeMore(t *testing.T) {
	c := &ColumnType{Name: Type_Set, Type: Scale_Type_More, Length: 20, Scale: 2}
	c.Set("1", "3", "5")
	assert.Equal(t, "SET('1','3','5')", c.Express())
	assert.False(t, c.IsInteger())
	assert.False(t, c.IsNumeric())
}
