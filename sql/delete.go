package sql

import (
	"strconv"
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type Delete struct {
	*base
	where       ksql.WhereInterface
	lowPriority string
	quick       string
	ignore      string
	as          string
	partitions  []string
	order       *orderInfo
	limit       string
	table       string
}

func NewDelete() *Delete {
	d := &Delete{base: newBase(), order: &orderInfo{}}
	d.opChain.Append(d._keyword, d._name, d._partition, d._where, d._order, d._limit)
	return d
}

func (d *Delete) _keyword(builder *strings.Builder) {
	builder.WriteString("DELETE")
	operator.BuildPureString(d.lowPriority, builder)
	operator.BuildPureString(d.quick, builder)
	operator.BuildPureString(d.ignore, builder)
	builder.WriteString(" FROM")
}

func (d *Delete) _name(builder *strings.Builder) {
	operator.BuildColumnString(d.table, builder)
	if d.as != "" {
		builder.WriteString(" AS")
		operator.BuildColumnString(d.as, builder)
	}
}

func (d *Delete) _partition(builder *strings.Builder) {
	if len(d.partitions) == 0 {
		return
	}

	builder.WriteString(" PARTITION ")
	for index, partition := range d.partitions {
		if index > 0 {
			builder.WriteString(",")
		}
		operator.BuildColumnString(partition, builder)
	}
}

func (d *Delete) _where(builder *strings.Builder) {
	if d.where == nil {
		return
	}

	builder.WriteString(" ")
	d.where.Build(builder)
	d.binds = append(d.binds, d.where.Binds()...)
}

func (d *Delete) _order(builder *strings.Builder) {
	if d.order.Empty() {
		return
	}

	d.order.Build(builder)
}

func (d *Delete) _limit(builder *strings.Builder) {
	if d.limit == "" {
		return
	}

	builder.WriteString("LIMIT")
	operator.BuildPureString(d.limit, builder)
}

func (d *Delete) Where(where ksql.WhereInterface) ksql.DeleteInterface {
	d.where = where
	return d
}

func (d *Delete) Table(table string) ksql.DeleteInterface {
	d.table = table
	return d
}

func (d *Delete) LowPriority() ksql.DeleteInterface {
	d.lowPriority = "LOW_PRIORITY"
	return d
}

func (d *Delete) Quick() ksql.DeleteInterface {
	d.quick = "QUICK"
	return d
}

func (d *Delete) Ignore() ksql.DeleteInterface {
	d.ignore = "IGNORE"
	return d
}

func (d *Delete) As(as string) ksql.DeleteInterface {
	d.as = as
	return d
}

func (d *Delete) Partitions(names ...string) ksql.DeleteInterface {
	d.partitions = append(d.partitions, names...)
	return d
}

func (d *Delete) OrderByDesc(columns ...string) ksql.DeleteInterface {
	for _, column := range columns {
		d.order.Append(&orderMeta{column: &columnInfo{column: column}, typ: "DESC"})
	}
	return d
}

func (d *Delete) OrderByAsc(columns ...string) ksql.DeleteInterface {
	for _, column := range columns {
		d.order.Append(&orderMeta{column: &columnInfo{column: column}, typ: "ASC"})
	}
	return d
}

func (d *Delete) Limit(limit int) ksql.DeleteInterface {
	d.limit = strconv.Itoa(limit)
	return d
}
