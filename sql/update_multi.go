package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type UpdateMulti struct {
	*base
	assignments *assignments
	where       ksql.WhereInterface
	table       string
	joins       []ksql.JoinInterface
	ignore      string
	priority    string
}

func NewUpdateMulti() *UpdateMulti {
	u := &UpdateMulti{base: newBase(), assignments: &assignments{}}
	u.opChain.Append(u._keyword, u._table, u._set, u._where)
	return u
}

func (u *UpdateMulti) _keyword(builder *strings.Builder) {
	builder.WriteString("UPDATE")
	operator.BuildPureString(u.priority, builder)
	operator.BuildPureString(u.ignore, builder)
}

func (u *UpdateMulti) _table(builder *strings.Builder) {
	operator.BuildColumnString(u.table, builder)
	for _, join := range u.joins {
		builder.WriteString(" ")
		join.Build(builder)
		u.binds = append(u.binds, join.Binds()...)
	}
}

func (u *UpdateMulti) _set(builder *strings.Builder) {
	builder.WriteString(" SET")
	u.assignments.Build(builder)
	u.binds = append(u.binds, u.assignments.binds...)
}

func (u *UpdateMulti) _where(builder *strings.Builder) {
	if u.where == nil || u.where.Empty() {
		return
	}

	builder.WriteString(" ")
	u.where.Build(builder)
	u.binds = append(u.binds, u.where.Binds()...)
}

func (u *UpdateMulti) Set(column string, data any) ksql.UpdateMultiInterface {
	u.assignments.Append(&assignment{column: column, data: data, isData: true})
	return u
}

func (u *UpdateMulti) Where(where ksql.WhereInterface) ksql.UpdateMultiInterface {
	u.where = where
	return u
}

func (u *UpdateMulti) Table(table string) ksql.UpdateMultiInterface {
	u.table = table
	return u
}

func (u *UpdateMulti) SetExpress(expre ksql.ExpressInterface) ksql.UpdateMultiInterface {
	u.assignments.Append(&assignment{expr: expre})
	return u
}

func (u *UpdateMulti) SetColumn(column string, otherColumn string) ksql.UpdateMultiInterface {
	u.assignments.Append(&assignment{column: column, value: otherColumn, isField: true})
	return u
}

func (u *UpdateMulti) Join(table string) ksql.JoinInterface {
	join := NewJoin().Inner()
	u.joins = append(u.joins, join)
	return join
}

func (u *UpdateMulti) JoinExpress(express ksql.ExpressInterface) ksql.JoinInterface {
	join := NewJoin().Express(express)
	u.joins = append(u.joins, join)
	return join
}

func (u *UpdateMulti) LeftJoin(table string) ksql.JoinInterface {
	join := NewJoin().Left()
	u.joins = append(u.joins, join)
	return join
}

func (u *UpdateMulti) RightJoin(table string) ksql.JoinInterface {
	join := NewJoin().Right()
	u.joins = append(u.joins, join)
	return join
}

func (u *UpdateMulti) LowPriority() ksql.UpdateMultiInterface {
	u.priority = "LOW_PRIORITY"
	return u
}

func (u *UpdateMulti) Ignore() ksql.UpdateMultiInterface {
	u.ignore = "IGNORE"
	return u
}
