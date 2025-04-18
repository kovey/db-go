package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type Do struct {
	*base
	sqls  []ksql.QueryInterface
	exprs []ksql.ExpressInterface
}

func NewDo() *Do {
	d := &Do{base: newBase()}
	d.opChain.Append(d._build)
	return d
}

func (d *Do) _build(builder *strings.Builder) {
	builder.WriteString("DO")
	index := 0
	for _, sql := range d.sqls {
		if index > 0 {
			builder.WriteString(", ")
		} else {
			builder.WriteString(" ")
		}

		builder.WriteString(sql.Prepare())
		d.binds = append(d.binds, sql.Binds()...)
		index++
	}

	for _, expr := range d.exprs {
		if index > 0 {
			builder.WriteString(", ")
		} else {
			builder.WriteString(" ")
		}

		builder.WriteString(expr.Statement())
		d.binds = append(d.binds, expr.Binds()...)
		index++
	}
}

func (d *Do) Do(query ksql.QueryInterface) ksql.DoInterface {
	d.sqls = append(d.sqls, query)
	return d
}

func (d *Do) DoExpress(express ksql.ExpressInterface) ksql.DoInterface {
	d.exprs = append(d.exprs, express)
	return d
}
