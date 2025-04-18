package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type ChangeColumn struct {
	oldColumn string
	newColumn *Column
	isFirst   bool
	after     string
	opChain   *operator.Chain
}

func NewChangeColumn() *ChangeColumn {
	c := &ChangeColumn{opChain: operator.NewChain()}
	c.opChain.Append(c._keyword, c._name, c._first)
	return c
}

func (c *ChangeColumn) _keyword(builder *strings.Builder) {
	builder.WriteString("CHANGE COLUMN")
}

func (c *ChangeColumn) _name(builder *strings.Builder) {
	if c.newColumn == nil {
		return
	}

	operator.BuildBacktickString(c.oldColumn, builder)
	builder.WriteString(" ")
	c.newColumn.Build(builder)
}

func (c *ChangeColumn) _first(builder *strings.Builder) {
	if c.isFirst {
		builder.WriteString(" FIRST")
		return
	}

	if c.after != "" {
		builder.WriteString(" AFTER")
		operator.BuildBacktickString(c.after, builder)
	}
}

func (c *ChangeColumn) Old(column string) ksql.ChangeColumnInterface {
	c.oldColumn = column
	return c
}

func (c *ChangeColumn) New(column, t string, length, scale int, sets ...string) ksql.ColumnInterface {
	typ := ParseType(t, length, scale, sets...)
	if typ == nil {
		return nil
	}

	c.newColumn = NewColumn(column, typ)
	return c.newColumn
}

func (c *ChangeColumn) First() ksql.ChangeColumnInterface {
	c.isFirst = true
	c.after = ""
	return c
}

func (c *ChangeColumn) After(column string) ksql.ChangeColumnInterface {
	c.isFirst = false
	c.after = column
	return c
}

func (c *ChangeColumn) Build(builder *strings.Builder) {
	c.opChain.Call(builder)
}
