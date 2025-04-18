package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type AlterColumn struct {
	column      string
	def         *Default
	visible     string
	dropDefault string
	opChain     *operator.Chain
}

func NewAlterColumn() *AlterColumn {
	a := &AlterColumn{opChain: operator.NewChain()}
	a.opChain.Append(a._keyword, a._column, a._default, a._visible, a._drop)
	return a
}

func (a *AlterColumn) _keyword(builder *strings.Builder) {
	builder.WriteString("ALTER COLUMN")
}

func (a *AlterColumn) _column(builder *strings.Builder) {
	builder.WriteString(" ")
	operator.Column(a.column, builder)
}

func (a *AlterColumn) _default(builder *strings.Builder) {
	if a.def == nil {
		return
	}

	builder.WriteString(" SET")
	a.def.Build(builder)
}

func (a *AlterColumn) _visible(builder *strings.Builder) {
	if a.visible == "" {
		return
	}

	builder.WriteString(" SET ")
	builder.WriteString(a.visible)
}

func (a *AlterColumn) _drop(builder *strings.Builder) {
	if a.dropDefault == "" {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(a.dropDefault)
}

func (a *AlterColumn) Column(column string) ksql.AlterColumnInterface {
	a.column = column
	return a
}

func (a *AlterColumn) Default(value string) ksql.AlterColumnInterface {
	a.def = &Default{Value: value, IsKeyword: ksql.IsDefaultKeyword(value)}
	return a
}

func (a *AlterColumn) DefaultExpress(expr string) ksql.AlterColumnInterface {
	a.def = &Default{Value: expr, IsExpr: true}
	return a
}

func (a *AlterColumn) Visible() ksql.AlterColumnInterface {
	a.visible = "VISIBLE"
	return a
}

func (a *AlterColumn) Invisible() ksql.AlterColumnInterface {
	a.visible = "INVISIBLE"
	return a
}

func (a *AlterColumn) DropDefault() ksql.AlterColumnInterface {
	a.dropDefault = "DROP DEFAULT"
	return a
}

func (a *AlterColumn) Build(builder *strings.Builder) {
	a.opChain.Call(builder)
}
