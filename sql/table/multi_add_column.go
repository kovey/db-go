package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type MultiAddColumn struct {
	columns []*Column
}

func NewMultiAddColumn() *MultiAddColumn {
	return &MultiAddColumn{}
}

func (a *MultiAddColumn) Column(column, t string, length, scale int, sets ...string) ksql.ColumnInterface {
	typ := ParseType(t, length, scale, sets...)
	if typ == nil {
		return nil
	}

	col := NewColumn(column, typ)
	a.columns = append(a.columns, col)
	return col
}

func (a *MultiAddColumn) Build(builder *strings.Builder) {
	builder.WriteString("ADD COLUMN (")
	for index, column := range a.columns {
		if index > 0 {
			builder.WriteString(", ")
		}
		column.Build(builder)
	}
	builder.WriteString(")")
}
