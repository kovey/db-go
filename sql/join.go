package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type joinOn struct {
	left    string
	op      string
	right   string
	value   any
	isField bool
}

type JoinOn struct {
	ons   []*joinOn
	binds []any
}

func (j *JoinOn) Empty() bool {
	return len(j.ons) == 0
}

func (j *JoinOn) Build(builder *strings.Builder) {
	if len(j.ons) == 0 {
		return
	}

	builder.WriteString("(")
	for index, on := range j.ons {
		if index > 0 {
			builder.WriteString(" AND ")
		}

		operator.Column(on.left, builder)
		operator.BuildPureString(on.op, builder)
		if on.isField {
			operator.BuildColumnString(on.right, builder)
		} else {
			operator.BuildPureString("?", builder)
			j.binds = append(j.binds, on.value)
		}
	}
	builder.WriteString(")")
}

func (j *JoinOn) On(left, op, right string) ksql.JoinOnInterface {
	j.ons = append(j.ons, &joinOn{left: left, op: op, right: right, isField: true})
	return j
}

func (j *JoinOn) OnVal(left, op string, right any) ksql.JoinOnInterface {
	j.ons = append(j.ons, &joinOn{left: left, op: op, value: right})
	return j
}

type Join struct {
	table     string
	as        string
	t         string
	ons       *JoinOn
	orOns     []*JoinOn
	isExpress bool
	expr      ksql.ExpressInterface
	opChain   *operator.Chain
	binds     []any
}

func NewJoin() *Join {
	j := &Join{opChain: operator.NewChain(), ons: &JoinOn{}}
	j.opChain.Append(j._express, j._join, j._on)
	return j
}

func (j *Join) _express(builder *strings.Builder) {
	if !j.isExpress {
		return
	}

	builder.WriteString(j.expr.Statement())
	j.binds = append(j.binds, j.expr.Binds()...)
}

func (j *Join) _join(builder *strings.Builder) {
	if j.isExpress {
		return
	}

	builder.WriteString(j.t)
	operator.BuildColumnString(j.table, builder)
}

func (j *Join) _on(builder *strings.Builder) {
	if j.as != "" {
		builder.WriteString(" AS")
		operator.BuildColumnString(j.as, builder)
	}

	if !j.ons.Empty() || len(j.orOns) > 0 {
		builder.WriteString(" ON ")
		canOr := false
		if !j.ons.Empty() {
			canOr = true
			j.ons.Build(builder)
			j.binds = append(j.binds, j.ons.binds...)
		}

		if len(j.orOns) > 0 {
			if canOr {
				builder.WriteString(" OR ")
			}

			for index, on := range j.orOns {
				if index > 0 {
					builder.WriteString(" OR ")
				}

				on.Build(builder)
				j.binds = append(j.binds, on.binds...)
			}
		}
	}
}

func (j *Join) Build(builder *strings.Builder) {
	j.opChain.Call(builder)
}

func (j *Join) Left() ksql.JoinInterface {
	if j.isExpress {
		return j
	}

	j.t = "LEFT JOIN"
	return j
}

func (j *Join) Right() ksql.JoinInterface {
	if j.isExpress {
		return j
	}

	j.t = "RIGHT JOIN"
	return j
}

func (j *Join) Inner() ksql.JoinInterface {
	if j.isExpress {
		return j
	}

	j.t = "INNER JOIN"
	return j
}

func (j *Join) Table(table string) ksql.JoinInterface {
	if j.isExpress {
		return j
	}

	j.table = table
	return j
}

func (j *Join) As(as string) ksql.JoinInterface {
	if j.isExpress {
		return j
	}

	j.as = as
	return j
}

func (j *Join) On(column, op, val string) ksql.JoinInterface {
	if j.isExpress {
		return j
	}

	j.ons.On(column, op, val)
	return j
}

func (j *Join) OnOr(call func(join ksql.JoinOnInterface)) ksql.JoinInterface {
	if j.isExpress {
		return j
	}

	join := &JoinOn{}
	call(join)
	j.orOns = append(j.orOns, join)
	return j
}

func (j *Join) Express(ex ksql.ExpressInterface) ksql.JoinInterface {
	j.isExpress = true
	j.expr = ex
	return j
}

func (j *Join) Binds() []any {
	return j.binds
}
