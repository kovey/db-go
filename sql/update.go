package sql

import (
	"strconv"
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type Update struct {
	*base
	assignments *assignments
	where       ksql.WhereInterface
	table       string
	priority    string
	ignore      string
	order       *orderInfo
	limit       string
}

func NewUpdate() *Update {
	u := &Update{base: newBase(), assignments: &assignments{}, order: &orderInfo{}}
	u.opChain.Append(u._keyword, u._set, u._where, u._order, u._limit)
	return u
}

func (u *Update) _keyword(builder *strings.Builder) {
	builder.WriteString("UPDATE")
	operator.BuildPureString(u.priority, builder)
	operator.BuildPureString(u.ignore, builder)
	operator.BuildColumnString(u.table, builder)
}

func (u *Update) _set(builder *strings.Builder) {
	builder.WriteString(" SET")
	u.assignments.Build(builder)
	u.binds = append(u.binds, u.assignments.binds...)
}

func (u *Update) _where(builder *strings.Builder) {
	if u.where == nil || u.where.Empty() {
		return
	}

	builder.WriteString(" ")
	u.where.Build(builder)
	u.binds = append(u.binds, u.where.Binds()...)
}

func (u *Update) _order(builder *strings.Builder) {
	if u.order.Empty() {
		return
	}

	u.order.Build(builder)
}

func (u *Update) _limit(builder *strings.Builder) {
	if u.limit == "" {
		return
	}

	builder.WriteString(" LIMIT ")
	builder.WriteString(u.limit)
}

func (u *Update) Table(table string) ksql.UpdateInterface {
	u.table = table
	return u
}

func (u *Update) Set(column string, data any) ksql.UpdateInterface {
	u.assignments.Append(&assignment{column: column, data: data, isData: true})
	return u
}

func (u *Update) Where(where ksql.WhereInterface) ksql.UpdateInterface {
	u.where = where
	return u
}

func (u *Update) LowPriority() ksql.UpdateInterface {
	u.priority = "LOW_PRIORITY"
	return u
}

func (u *Update) Ignore() ksql.UpdateInterface {
	u.ignore = "IGNORE"
	return u
}

func (u *Update) OrderByAsc(columns ...string) ksql.UpdateInterface {
	for _, column := range columns {
		u.order.Append(&orderMeta{column: &columnInfo{column: column}, typ: "ASC"})
	}
	return u
}

func (u *Update) OrderByDesc(columns ...string) ksql.UpdateInterface {
	for _, column := range columns {
		u.order.Append(&orderMeta{column: &columnInfo{column: column}, typ: "DESC"})
	}
	return u
}

func (u *Update) SetExpress(expre ksql.ExpressInterface) ksql.UpdateInterface {
	u.assignments.Append(&assignment{expr: expre})
	return u
}

func (u *Update) SetColumn(column string, otherColumn string) ksql.UpdateInterface {
	u.assignments.Append(&assignment{column: column, value: otherColumn, isField: true})
	return u
}

func (u *Update) IncColumn(column string, data int) ksql.UpdateInterface {
	if data >= 0 {
		u.assignments.Append(&assignment{column: column, data: data, isData: true, op: "+"})
		return u
	}
	u.assignments.Append(&assignment{column: column, data: 0 - data, isData: true, op: "-"})
	return u
}

func (u *Update) Limit(limit int) ksql.UpdateInterface {
	u.limit = strconv.Itoa(limit)
	return u
}
