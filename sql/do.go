package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type Do struct {
	*base
	exprs []ksql.ExpressInterface
}

func NewDo() *Do {
	d := &Do{base: newBase()}
	d.opChain.Append(d._build)
	return d
}

func (d *Do) _build(builder *strings.Builder) {
	if len(d.exprs) == 0 {
		return
	}

	builder.WriteString("DO ")
	for index, expr := range d.exprs {
		if index > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(expr.Statement())
		d.binds = append(d.binds, expr.Binds()...)
	}
}

func (d *Do) Do(expr ksql.ExpressInterface) ksql.DoInterface {
	d.exprs = append(d.exprs, expr)
	return d
}
