package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultIsKeyword(t *testing.T) {
	d := &Default{IsKeyword: true, Value: "aaa", IsByte: false}
	assert.Equal(t, "DEFAULT aaa", d.Express())
	d = &Default{IsKeyword: true, Value: "aaa", IsByte: true}
	assert.Equal(t, "DEFAULT aaa", d.Express())
}

func TestDefaultIsNotKeyword(t *testing.T) {
	d := &Default{IsKeyword: false, Value: "aaa", IsByte: false}
	assert.Equal(t, "DEFAULT 'aaa'", d.Express())
	d = &Default{IsKeyword: false, Value: "aaa", IsByte: true}
	assert.Equal(t, "DEFAULT b'aaa'", d.Express())
}

func TestColumnNullable(t *testing.T) {
	c := NewColumn("user_name", &ColumnType{Name: Type_VarChar, Type: Scale_Type_One, Length: 20})
	c.Nullable().Comment("user name").Default("").Unsigned()
	assert.Equal(t, "`user_name` VARCHAR(20) NULL DEFAULT '' COMMENT 'user name'", c.Express())
}

func TestColumnAutoInc(t *testing.T) {
	c := NewColumn("user_id", &ColumnType{Name: "BIGINT", Length: 20})
	c.Nullable().Comment("main key").Default("").Unsigned().AutoIncrement()
	assert.Equal(t, "`user_id` BIGINT unsigned NOT NULL AUTO_INCREMENT COMMENT 'main key'", c.Express())
	c = NewColumn("user_name", &ColumnType{Name: Type_VarChar, Length: 20})
	c.Nullable().Comment("main key").Default("").Unsigned().AutoIncrement()
	assert.Equal(t, "`user_name` VARCHAR NULL DEFAULT '' COMMENT 'main key'", c.Express())
}

func TestColumnUseCurrent(t *testing.T) {
	c := NewColumn("user_date", &ColumnType{Name: Type_DateTime})
	c.Nullable().Comment("date").UseCurrent().Default("").Unsigned().AutoIncrement()
	assert.Equal(t, "`user_date` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'date'", c.Express())
}

func TestColumnUseCurrentOnUpdate(t *testing.T) {
	c := NewColumn("user_date", &ColumnType{Name: Type_DateTime})
	c.Nullable().Comment("date").UseCurrentOnUpdate().Default("").Unsigned().AutoIncrement()
	assert.Equal(t, "`user_date` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'date'", c.Express())
}
