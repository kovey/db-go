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

func (r *RenameIndex) Old(index string) ksql.RenameIndexInterface {
	r.oldIndex = index
	return r
}

func (r *RenameIndex) New(index string) ksql.RenameIndexInterface {
	r.newIndex = index
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
