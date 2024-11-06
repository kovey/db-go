package table

import (
	"fmt"

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
}

func NewColumn(name string, t *ColumnType) *Column {
	return &Column{name: name, t: t}
}

func (c *Column) Express() string {
	null := "NULL"
	if !c.isNull {
		null = "NOT NULL"
	}

	unsigned := ""
	if c.isUnsigned {
		unsigned = "unsigned"
	}

	autoInc := ""
	if c.isAutoInc {
		autoInc = "AUTO_INCREMENT"
	}

	if c.t.Name == "BIT" && c.def != nil {
		c.def.IsByte = true
	}

	if c.def != nil {
		return fmt.Sprintf("`%s` %s %s %s %s %s %s", c.name, c.t.Express(), unsigned, null, autoInc, c.def.Express(), c.comment)
	}

	return fmt.Sprintf("`%s` %s %s %s %s %s", c.name, c.t.Express(), unsigned, null, autoInc, c.comment)
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
	c.def = &Default{Value: value, IsKeyword: ksql.IsDefaultKeyword(value)}
	return c
}

func (c *Column) Comment(comment string) ksql.ColumnInterface {
	c.comment = fmt.Sprintf("COMMENT '%s'", comment)
	return c
}
