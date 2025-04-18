package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type AddColumn struct {
	column  *Column
	isFirst bool
	after   string
}

func NewAddColumn() *AddColumn {
	return &AddColumn{}
}

func (a *AddColumn) Column(column, t string, length, scale int, sets ...string) ksql.ColumnInterface {
	typ := ParseType(t, length, scale, sets...)
	if typ == nil {
		return nil
	}

	a.column = NewColumn(column, typ)
	return a.column
}

func (a *AddColumn) First() ksql.TableAddColumnInterface {
	a.isFirst = true
	a.after = ""
	return nil
}

func (a *AddColumn) After(column string) ksql.TableAddColumnInterface {
	a.isFirst = false
	a.after = column
	return nil
}

func (a *AddColumn) Build(builder *strings.Builder) {
	builder.WriteString("ADD COLUMN ")
	a.column.Build(builder)
	if a.isFirst {
		builder.WriteString(" FIRST")
		return
	}

	if a.after != "" {
		builder.WriteString(" AFTER")
		operator.BuildBacktickString(a.after, builder)
	}
}
