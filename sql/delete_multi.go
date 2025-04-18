package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type DeleteMulti struct {
	*base
	where       ksql.WhereInterface
	lowPriority string
	quick       string
	ignore      string
	tables      []string
	joins       []ksql.JoinInterface
}

func NewDeleteMulti() *DeleteMulti {
	d := &DeleteMulti{base: newBase()}
	d.opChain.Append(d._keyword, d._name, d._reference, d._where)
	return d
}

func (d *DeleteMulti) _keyword(builder *strings.Builder) {
	builder.WriteString("DELETE")
	operator.BuildPureString(d.lowPriority, builder)
	operator.BuildPureString(d.quick, builder)
	operator.BuildPureString(d.ignore, builder)
}

func (d *DeleteMulti) _name(builder *strings.Builder) {
	for index, table := range d.tables {
		if index > 0 {
			builder.WriteString(",")
		}

		operator.BuildColumnString(table, builder)
	}
}

func (d *DeleteMulti) _reference(builder *strings.Builder) {
	builder.WriteString(" FROM")
	for index, join := range d.joins {
		if index > 0 {
			builder.WriteString(" ")
		}

		join.Build(builder)
		d.binds = append(d.binds, join.Binds()...)
	}
}

func (d *DeleteMulti) _where(builder *strings.Builder) {
	if d.where == nil {
		return
	}

	builder.WriteString(" ")
	d.where.Build(builder)
	d.binds = append(d.binds, d.where.Binds()...)
}

func (d *DeleteMulti) Where(where ksql.WhereInterface) ksql.DeleteMultiInterface {
	d.where = where
	return d
}

func (d *DeleteMulti) Table(table string) ksql.DeleteMultiInterface {
	d.tables = append(d.tables, table)
	return d
}

func (d *DeleteMulti) LowPriority() ksql.DeleteMultiInterface {
	d.lowPriority = "LOW_PRIORITY"
	return d
}

func (d *DeleteMulti) Quick() ksql.DeleteMultiInterface {
	d.quick = "QUICK"
	return d
}

func (d *DeleteMulti) Ignore() ksql.DeleteMultiInterface {
	d.ignore = "IGNORE"
	return d
}

func (d *DeleteMulti) _join(table string) ksql.JoinInterface {
	join := NewJoin()
	join.Table(table)
	d.joins = append(d.joins, join)
	return join
}

func (d *DeleteMulti) _joinExpress(express ksql.ExpressInterface) ksql.JoinInterface {
	join := NewJoin()
	join.Express(express)
	d.joins = append(d.joins, join)
	return join
}

func (d *DeleteMulti) Join(table string) ksql.JoinInterface {
	return d._join(table).Inner()
}

func (d *DeleteMulti) JoinExpress(express ksql.ExpressInterface) ksql.JoinInterface {
	return d._joinExpress(express)
}

func (d *DeleteMulti) LeftJoin(table string) ksql.JoinInterface {
	return d._join(table).Left()
}

func (d *DeleteMulti) RightJoin(table string) ksql.JoinInterface {
	return d._join(table).Right()
}
