package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type ModifyColumn struct {
	column  *Column
	isFirst bool
	after   string
	opChain *operator.Chain
	colName string
}

func NewModifyColumn(name string) *ModifyColumn {
	m := &ModifyColumn{opChain: operator.NewChain(), colName: name}
	m.opChain.Append(m._keyword, m._column, m._first)
	return m
}

func (m *ModifyColumn) _keyword(builder *strings.Builder) {
	builder.WriteString("MODIFY COLUMN")
}

func (m *ModifyColumn) _column(builder *strings.Builder) {
	if m.column == nil {
		return
	}

	builder.WriteString(" ")
	m.column.Build(builder)
}

func (m *ModifyColumn) _first(builder *strings.Builder) {
	if m.isFirst {
		builder.WriteString(" FIRST")
		return
	}

	if m.after != "" {
		builder.WriteString(" AFTER")
		operator.BuildBacktickString(m.after, builder)
	}
}

func (m *ModifyColumn) Column(t string, length, scale int, sets ...string) ksql.ColumnInterface {
	typ := ParseType(t, length, scale, sets...)
	if typ == nil {
		return nil
	}

	m.column = NewColumn(m.colName, typ)
	return m.column
}

func (m *ModifyColumn) First() ksql.ModifyColumnInterface {
	m.isFirst = true
	m.after = ""
	return m
}

func (m *ModifyColumn) After(column string) ksql.ModifyColumnInterface {
	m.isFirst = false
	m.after = column
	return m
}

func (m *ModifyColumn) Build(builder *strings.Builder) {
	m.opChain.Call(builder)
}
