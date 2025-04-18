package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type RenameIndex struct {
	oldIndex string
	newIndex string
	typ      ksql.IndexSubType
}

func (r *RenameIndex) Old(column string) ksql.RenameIndexInterface {
	r.oldIndex = column
	return r
}

func (r *RenameIndex) New(column string) ksql.RenameIndexInterface {
	r.newIndex = column
	return r
}

func (r *RenameIndex) Type(typ ksql.IndexSubType) ksql.RenameIndexInterface {
	r.typ = typ
	return r
}

func (r *RenameIndex) Build(builder *strings.Builder) {
	builder.WriteString("RENAME ")
	builder.WriteString(string(r.typ))
	operator.BuildBacktickString(r.oldIndex, builder)
	builder.WriteString(" TO")
	operator.BuildBacktickString(r.newIndex, builder)
}
