package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type Query struct {
	*base
	table     string
	as        string
	columns   strings.Builder
	where     ksql.WhereInterface
	join      []ksql.JoinInterface
	limit     int
	offset    int
	order     strings.Builder
	group     strings.Builder
	having    ksql.HavingInterface
	forUpdate bool
}

func NewQuery() *Query {
	q := &Query{base: &base{hasPrepared: false}, where: NewWhere(), having: NewHaving()}
	q.keyword("SELECT ")
	return q
}

func (o *Query) Clone() ksql.QueryInterface {
	q := &Query{
		base: &base{hasPrepared: false, binds: o.binds}, where: o.where, having: o.having,
		table: o.table, as: o.as, join: o.join, group: o.group,
	}
	q.keyword("SELECT ")
	return q
}

func (o *Query) TableBy(operater ksql.QueryInterface, as string) ksql.QueryInterface {
	var tmp = make([]any, len(operater.Binds()))
	copy(tmp, operater.Binds())
	o.binds = append(tmp, o.binds...)
	o.as = as
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(operater.Prepare())
	builder.WriteString(")")
	o.table = builder.String()
	return o
}

func (o *Query) Table(table string) ksql.QueryInterface {
	o.table = table
	return o
}

func (o *Query) As(as string) ksql.QueryInterface {
	o.as = as
	return o
}

func (o *Query) Func(fun, column, as string) ksql.QueryInterface {
	if o.columns.Len() > 0 {
		o.columns.WriteString(",")
	}

	o.columns.WriteString(fun)
	o.columns.WriteString("(")
	Column(column, &o.columns)
	o.columns.WriteString(")")
	o.columns.WriteString(" AS ")
	Backtick(as, &o.columns)

	return o
}

func (o *Query) Column(column, as string) ksql.QueryInterface {
	if o.columns.Len() > 0 {
		o.columns.WriteString(",")
	}

	Column(column, &o.columns)
	o.columns.WriteString(" AS ")
	Backtick(as, &o.columns)
	return o
}

func (o *Query) Columns(columns ...string) ksql.QueryInterface {
	for _, column := range columns {
		if o.columns.Len() > 0 {
			o.columns.WriteString(",")
		}

		Column(column, &o.columns)
	}

	return o
}

func (o *Query) ColumnsExpress(columns ...ksql.ExpressInterface) ksql.QueryInterface {
	for _, column := range columns {
		if o.columns.Len() > 0 {
			o.columns.WriteString(",")
		}

		o.columns.WriteString(column.Statement())
	}

	return o
}

func (o *Query) _data(data ...any) {
	o.binds = append(o.binds, data...)
}

func (o *Query) WhereIsNull(column string) ksql.QueryInterface {
	o.where.IsNull(column)
	return o
}

func (o *Query) WhereIsNotNull(column string) ksql.QueryInterface {
	o.where.IsNotNull(column)
	return o
}

func (o *Query) WhereIn(column string, val []any) ksql.QueryInterface {
	o.where.In(column, val)
	return o
}

func (o *Query) WhereNotIn(column string, val []any) ksql.QueryInterface {
	o.where.NotIn(column, val)
	return o
}

func (o *Query) WhereInBy(column string, sub ksql.QueryInterface) ksql.QueryInterface {
	o.where.InBy(column, sub)
	return o
}

func (o *Query) WhereNotInBy(column string, sub ksql.QueryInterface) ksql.QueryInterface {
	o.where.NotInBy(column, sub)
	return o
}

func (o *Query) Where(column string, op string, val any) ksql.QueryInterface {
	o.where.Where(column, op, val)
	return o
}

func (o *Query) WhereExpress(expresses ...ksql.ExpressInterface) ksql.QueryInterface {
	for _, express := range expresses {
		o.where.Express(express)
	}

	return o
}

func (o *Query) OrWhere(call func(w ksql.WhereInterface)) ksql.QueryInterface {
	o.where.OrWhere(call)
	return o
}

func (o *Query) Between(column string, begin, end any) ksql.QueryInterface {
	o.where.Between(column, begin, end)
	return o
}

func (o *Query) Limit(limit int) ksql.QueryInterface {
	o.limit = limit
	return o
}

func (o *Query) Offset(offset int) ksql.QueryInterface {
	o.offset = offset
	return o
}

func (o *Query) Order(column string) ksql.QueryInterface {
	if o.order.Len() > 0 {
		o.order.WriteString(",")
	} else {
		o.order.WriteString(" ORDER BY ")
	}
	Column(column, &o.order)
	o.order.WriteString(" ASC")
	return o
}

func (o *Query) OrderDesc(column string) ksql.QueryInterface {
	if o.order.Len() > 0 {
		o.order.WriteString(",")
	} else {
		o.order.WriteString(" ORDER BY ")
	}
	Column(column, &o.order)
	o.order.WriteString(" DESC")
	return o
}

func (o *Query) Group(columns ...string) ksql.QueryInterface {
	if o.group.Len() > 0 {
		o.group.WriteString(",")
	} else {
		o.group.WriteString(" GROUP BY ")
	}
	Column(columns[0], &o.group)
	for _, column := range columns[1:] {
		o.group.WriteString(",")
		Column(column, &o.group)
	}

	return o
}

func (o *Query) _join(t string, table string) ksql.JoinInterface {
	join := NewJoin(t)
	join.Table(table)
	o.join = append(o.join, join)
	return join
}

func (o *Query) Join(table string) ksql.JoinInterface {
	return o._join(" JOIN ", table)
}

func (o *Query) _joinExpress(t string, express ksql.ExpressInterface) ksql.JoinInterface {
	join := NewJoin(t)
	join.Express(express)
	o.join = append(o.join, join)
	return join
}

func (o *Query) JoinExpress(express ksql.ExpressInterface) ksql.JoinInterface {
	return o._joinExpress("JOIN", express)
}

func (o *Query) LeftJoin(table string) ksql.JoinInterface {
	return o._join("LEFT JOIN", table)
}

func (o *Query) LeftJoinExpress(express ksql.ExpressInterface) ksql.JoinInterface {
	return o._joinExpress("LEFT JOIN", express)
}

func (o *Query) RightJoin(table string) ksql.JoinInterface {
	return o._join("RIGHT JOIN", table)
}

func (o *Query) RightJoinExpress(express ksql.ExpressInterface) ksql.JoinInterface {
	return o._joinExpress("RIGHT JOIN", express)
}

func (o *Query) Prepare() string {
	if o.hasPrepared {
		return o.base.Prepare()
	}

	o.hasPrepared = true
	if o.columns.Len() == 0 {
		o.builder.WriteString("*")
	} else {
		o.builder.WriteString(o.columns.String())
	}

	o.builder.WriteString(" FROM ")
	Column(o.table, &o.builder)
	if o.as != "" {
		o.builder.WriteString(" AS ")
		Backtick(o.as, &o.builder)
	}

	for _, join := range o.join {
		o.builder.WriteString(" ")
		o.builder.WriteString(join.Prepare())
		o._data(join.Binds()...)
	}
	if !o.where.Empty() {
		o.builder.WriteString(" WHERE ")
		o.builder.WriteString(o.where.Prepare())
		o._data(o.where.Binds()...)
	}

	if o.group.Len() > 0 {
		o.builder.WriteString(o.group.String())
	}

	if !o.having.Empty() {
		o.builder.WriteString(" HAVING ")
		o.builder.WriteString(o.having.Prepare())
		o._data(o.having.Binds()...)
	}

	if o.order.Len() > 0 {
		o.builder.WriteString(o.order.String())
	}

	if o.limit > 0 {
		o.builder.WriteString(" LIMIT ")
		o.builder.WriteString(RawValue(o.offset))
		o.builder.WriteString(",")
		o.builder.WriteString(RawValue(o.limit))
	}

	if o.forUpdate {
		o.builder.WriteString(" FOR UPDATE")
	}

	return o.base.Prepare()
}

func (o *Query) Having(column string, op string, val any) ksql.QueryInterface {
	o.having.Having(column, op, val)
	return o
}

func (o *Query) HavingExpress(expresses ...ksql.ExpressInterface) ksql.QueryInterface {
	for _, raw := range expresses {
		o.having.Express(raw)
	}
	return o
}

func (o *Query) OrHaving(call func(h ksql.HavingInterface)) ksql.QueryInterface {
	o.having.OrHaving(call)
	return o
}

func (o *Query) ForUpdate() ksql.QueryInterface {
	o.forUpdate = true
	return o
}
