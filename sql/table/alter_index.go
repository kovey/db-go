package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type AlterIndex struct {
	index   string
	visible string
}

func NewAlterIndex() *AlterIndex {
	return &AlterIndex{}
}

func (a *AlterIndex) Index(index string) ksql.AlterIndexInterface {
	a.index = index
	return a
}

func (a *AlterIndex) Visible() ksql.AlterIndexInterface {
	a.visible = "VISIBLE"
	return a
}

func (a *AlterIndex) Invisible() ksql.AlterIndexInterface {
	a.visible = "INVISIBLE"
	return a
}

func (a *AlterIndex) Build(builder *strings.Builder) {
	builder.WriteString("ALTER INDEX ")
	operator.Column(a.index, builder)
	operator.BuildPureString(a.visible, builder)
}
