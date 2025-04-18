package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type TableRename struct {
	*base
	froms []string
	toes  []string
}

func NewTableRename() *TableRename {
	t := &TableRename{base: newBase()}
	t.opChain.Append(t._build)
	return t
}

func (t *TableRename) _build(builder *strings.Builder) {
	builder.WriteString("RENAME TABLE")
	for index, from := range t.froms {
		if index > 0 {
			builder.WriteString(",")
		}

		operator.BuildColumnString(from, builder)
		builder.WriteString(" TO")
		operator.BuildColumnString(t.toes[index], builder)
	}
}

func (t *TableRename) Table(from, to string) ksql.RenameTableInterface {
	t.froms = append(t.froms, from)
	t.toes = append(t.toes, to)
	return t
}
