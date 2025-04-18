package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type RenameColumn struct {
	oldColumn string
	newColumn string
}

func (r *RenameColumn) Old(column string) ksql.RenameColumnInterface {
	r.oldColumn = column
	return r
}

func (r *RenameColumn) New(column string) ksql.RenameColumnInterface {
	r.newColumn = column
	return r
}

func (r *RenameColumn) Build(builder *strings.Builder) {
	builder.WriteString("RENAME COLUMN")
	operator.BuildBacktickString(r.oldColumn, builder)
	builder.WriteString(" TO")
	operator.BuildBacktickString(r.newColumn, builder)
}
