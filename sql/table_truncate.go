package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type TableTruncate struct {
	*base
	table string
}

func NewTableTruncate() *TableTruncate {
	t := &TableTruncate{base: newBase()}
	t.opChain.Append(t._build)
	return t
}

func (t *TableTruncate) _build(builder *strings.Builder) {
	builder.WriteString("TRUNCATE TABLE")
	operator.BuildColumnString(t.table, builder)
}

func (t *TableTruncate) Table(table string) ksql.TruncateTableInterface {
	t.table = table
	return t
}
