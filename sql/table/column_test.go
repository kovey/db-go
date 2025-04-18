package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultIsKeyword(t *testing.T) {
	var builder strings.Builder
	d := &Default{IsKeyword: true, Value: "aaa", IsByte: false}
	d.Build(&builder)
	assert.Equal(t, " DEFAULT aaa", builder.String())
	builder.Reset()
	d = &Default{IsKeyword: true, Value: "aaa", IsByte: true}
	d.Build(&builder)
	assert.Equal(t, " DEFAULT aaa", builder.String())
}

func TestDefaultIsNotKeyword(t *testing.T) {
	var builder strings.Builder
	d := &Default{IsKeyword: false, Value: "aaa", IsByte: false}
	d.Build(&builder)
	assert.Equal(t, " DEFAULT 'aaa'", builder.String())
	builder.Reset()
	d = &Default{IsKeyword: false, Value: "aaa", IsByte: true}
	d.Build(&builder)
	assert.Equal(t, " DEFAULT b'aaa'", builder.String())
}

func TestColumnNullable(t *testing.T) {
	var builder strings.Builder
	c := NewColumn("user_name", &ColumnType{Name: Type_VarChar, Type: Scale_Type_One, Length: 20})
	c.Nullable().Comment("user name").Default("").Unsigned()
	c.Build(&builder)
	assert.Equal(t, "`user_name` VARCHAR(20) NULL DEFAULT '' COMMENT 'user name'", builder.String())
}

func TestColumnNotNullable(t *testing.T) {
	var builder strings.Builder
	c := NewColumn("user_name", &ColumnType{Name: Type_VarChar, Type: Scale_Type_One, Length: 20})
	c.Comment("user name").Default("").Unsigned().NotNullable()
	c.Build(&builder)
	assert.Equal(t, "`user_name` VARCHAR(20) NOT NULL DEFAULT '' COMMENT 'user name'", builder.String())
}

func TestColumnIsByte(t *testing.T) {
	var builder strings.Builder
	c := NewColumn("user_name", &ColumnType{Name: Type_Bit, Type: Scale_Type_One, Length: 20})
	c.Comment("user name").DefaultBit("").Unsigned()
	c.Build(&builder)
	assert.Equal(t, "`user_name` BIT(20) DEFAULT b'' COMMENT 'user name'", builder.String())
}

func TestColumnAutoInc(t *testing.T) {
	var builder strings.Builder
	c := NewColumn("user_id", &ColumnType{Name: "BIGINT", Length: 20})
	c.Nullable().Comment("main key").Unsigned().AutoIncrement()
	c.Build(&builder)
	assert.Equal(t, "`user_id` BIGINT UNSIGNED AUTO_INCREMENT COMMENT 'main key'", builder.String())
	c = NewColumn("user_name", &ColumnType{Name: Type_VarChar, Length: 20})
	c.Nullable().Comment("main key").Default("").Unsigned().AutoIncrement()
	builder.Reset()
	c.Build(&builder)
	assert.Equal(t, "`user_name` VARCHAR NULL DEFAULT '' COMMENT 'main key'", builder.String())
}

func TestColumnUseCurrent(t *testing.T) {
	var builder strings.Builder
	c := NewColumn("user_date", &ColumnType{Name: Type_DateTime})
	c.Nullable().Comment("date").UseCurrent().Default("").Unsigned().AutoIncrement()
	c.Build(&builder)
	assert.Equal(t, "`user_date` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'date'", builder.String())
}

func TestColumnUseCurrentOnUpdate(t *testing.T) {
	var builder strings.Builder
	c := NewColumn("user_date", &ColumnType{Name: Type_DateTime})
	c.Nullable().Comment("date").UseCurrentOnUpdate().Default("").Unsigned().AutoIncrement()
	c.Build(&builder)
	assert.Equal(t, "`user_date` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'date'", builder.String())
}
