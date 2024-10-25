package sql

import (
	"strconv"
	"strings"

	"github.com/kovey/db-go/v3"
)

type Join struct {
	*base
	table     string
	as        string
	t         string
	isExpress bool
}

func NewJoin(t string) *Join {
	return &Join{base: &base{hasPrepared: false}, t: t}
}

func (j *Join) Type() string {
	return j.t
}

func (j *Join) Table(table string) ksql.JoinInterface {
	j.table = table
	return j
}

func (j *Join) As(as string) ksql.JoinInterface {
	j.as = as
	return j
}

func (j *Join) _on(column, op, val, condition string) *Join {
	if j.isExpress {
		return j
	}

	if j.builder.Len() == 0 {
		if j.t != "" {
			j.builder.WriteString(j.t)
			j.builder.WriteString(" ")
			Column(j.table, &j.builder)
			if j.as != "" {
				j.builder.WriteString(" AS ")
				Column(j.as, &j.builder)
			}
			j.builder.WriteString(" ON ")
		}
	} else {
		j.builder.WriteString(" ")
		j.builder.WriteString(condition)
		j.builder.WriteString(" ")
	}

	Column(column, &j.builder)
	j.builder.WriteString(op)
	if strings.Contains(val, ".") {
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			Column(val, &j.builder)
			return j
		}
	}

	j.builder.WriteString(val)
	return j
}

func (j *Join) On(column, op, val string) ksql.JoinInterface {
	return j._on(column, op, val, "AND")
}

func (j *Join) OnOr(call func(join ksql.JoinOnInterface)) ksql.JoinInterface {
	if j.builder.Len() == 0 {
		return j
	}

	join := NewJoin("")
	call(join)
	j.builder.WriteString(" OR (")
	j.builder.WriteString(join.Prepare())
	j.builder.WriteString(")")
	return j
}

func (j *Join) Express(ex ksql.ExpressInterface) ksql.JoinInterface {
	j.isExpress = true
	j.builder.Reset()
	j.builder.WriteString(j.t)
	j.builder.WriteString(" ")
	j.builder.WriteString(ex.Statement())
	j.binds = append(j.binds, ex.Binds()...)
	return j
}
