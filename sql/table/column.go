package table

import (
	"fmt"
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type Default struct {
	IsKeyword bool
	Value     string
	IsByte    bool
}

func (d *Default) Express() string {
	if d.IsKeyword {
		return fmt.Sprintf("DEFAULT %s", d.Value)
	}

	if d.IsByte {
		return fmt.Sprintf("DEFAULT b'%s'", d.Value)
	}

	return fmt.Sprintf("DEFAULT '%s'", d.Value)
}

type Column struct {
	name       string
	t          *ColumnType
	isNull     bool
	def        *Default
	comment    string
	isAutoInc  bool
	isUnsigned bool
	isBuild    bool
	builder    strings.Builder
}

func NewColumn(name string, t *ColumnType) *Column {
	return &Column{name: name, t: t}
}

func (c *Column) Express() string {
	if c.isBuild {
		return c.builder.String()
	}

	c.isBuild = true
	null := "NULL"
	if !c.isNull {
		null = "NOT NULL"
	}

	unsigned := ""
	if c.isUnsigned && c.t.IsNumeric() {
		unsigned = "unsigned"
	}

	autoInc := ""
	if c.isAutoInc && c.t.IsInteger() {
		autoInc = "AUTO_INCREMENT"
		null = "NOT NULL"
		c.def = nil
	}

	if c.t.Name == "BIT" && c.def != nil {
		c.def.IsByte = true
	}

	builder := strings.Builder{}
	builder.WriteString("`")
	builder.WriteString(c.name)
	builder.WriteString("` ")
	builder.WriteString(c.t.Express())
	if unsigned != "" {
		builder.WriteString(" ")
		builder.WriteString(unsigned)
	}
	builder.WriteString(" ")
	builder.WriteString(null)
	if autoInc != "" {
		builder.WriteString(" ")
		builder.WriteString(autoInc)
	}
	if c.def != nil {
		builder.WriteString(" ")
		builder.WriteString(c.def.Express())
	}
	if c.comment != "" {
		builder.WriteString(" ")
		builder.WriteString(c.comment)
	}

	return builder.String()
}

func (c *Column) Nullable() ksql.ColumnInterface {
	c.isNull = true
	return c
}

func (c *Column) AutoIncrement() ksql.ColumnInterface {
	c.isAutoInc = true
	return c
}

func (c *Column) Unsigned() ksql.ColumnInterface {
	c.isUnsigned = true
	return c
}

func (c *Column) UseCurrent() ksql.ColumnInterface {
	c.Default(ksql.CURRENT_TIMESTAMP)
	return c
}

func (c *Column) UseCurrentOnUpdate() ksql.ColumnInterface {
	c.Default(ksql.CURRENT_TIMESTAMP_ON_UPDATE_CURRENT_TIMESTAMP)
	return c
}

func (c *Column) Default(value string) ksql.ColumnInterface {
	if c.def != nil {
		return c
	}

	c.def = &Default{Value: value, IsKeyword: ksql.IsDefaultKeyword(value)}
	return c
}

func (c *Column) Comment(comment string) ksql.ColumnInterface {
	c.comment = fmt.Sprintf("COMMENT '%s'", comment)
	return c
}
