package sql

import (
	"fmt"
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type Having struct {
	ops       []*whereOp
	orWheres  []*Having
	andWheres []*Having
	opChain   *operator.Chain
	onlyBody  bool
	binds     []any
}

func NewHaving() *Having {
	h := &Having{opChain: operator.NewChain()}
	h.opChain.Append(h._keyword, h._ops, h._andWheres, h._orWheres)
	return h
}

func (w *Having) Clone() ksql.HavingInterface {
	o := NewHaving()
	o.ops = w.ops
	o.orWheres = w.orWheres
	o.andWheres = w.andWheres
	return o
}

func (w *Having) _keyword(builder *strings.Builder) {
	if w.onlyBody {
		return
	}

	builder.WriteString("HAVING")
}

func (w *Having) _ops(builder *strings.Builder) {
	if len(w.ops) == 0 {
		return
	}

	if !w.onlyBody {
		builder.WriteString(" ")
	}
	for index, op := range w.ops {
		if index > 0 {
			builder.WriteString(" AND ")
		}
		op.Build(builder)
		w.binds = append(w.binds, op.binds()...)
	}
}

func (w *Having) _andWheres(builder *strings.Builder) {
	if len(w.ops) == 0 {
		w._where("", "AND", w.andWheres, builder)
		return
	}

	w._where("AND", "AND", w.andWheres, builder)
}

func (w *Having) _orWheres(builder *strings.Builder) {
	if len(w.ops) == 0 && len(w.andWheres) == 0 {
		w._where("", "OR", w.orWheres, builder)
		return
	}

	w._where("OR", "OR", w.orWheres, builder)
}

func (w *Having) _where(prefix, op string, wheres []*Having, builder *strings.Builder) {
	if len(wheres) == 0 {
		return
	}

	builder.WriteString(" ")
	if prefix != "" {
		builder.WriteString(prefix)
	}

	for index, where := range wheres {
		if index > 0 {
			builder.WriteString(" ")
			builder.WriteString(op)
		}

		builder.WriteString(" (")
		where.Build(builder)
		builder.WriteString(")")
		w.binds = append(w.binds, where.Binds()...)
	}
}

func (w *Having) Binds() []any {
	return w.binds
}

func (w *Having) Build(builder *strings.Builder) {
	w.opChain.Call(builder)
}

func (w *Having) Having(column string, op ksql.Op, data any) ksql.HavingInterface {
	if !ksql.SupportOp(op) {
		panic(fmt.Sprintf("op %s not support", op))
	}

	w.ops = append(w.ops, &whereOp{column: column, op: op, value: data})
	return w
}

func (w *Having) In(column string, data []any) ksql.HavingInterface {
	w.ops = append(w.ops, &whereOp{column: column, values: data, isArr: true, op: "IN"})
	return w
}

func (w *Having) NotIn(column string, data []any) ksql.HavingInterface {
	w.ops = append(w.ops, &whereOp{column: column, values: data, isArr: true, op: "NOT IN"})
	return w
}

func (w *Having) _is(column string, op string) ksql.HavingInterface {
	w.ops = append(w.ops, &whereOp{column: column, op: "IS", constValue: op, isConst: true})
	return w
}

func (w *Having) IsNull(column string) ksql.HavingInterface {
	return w._is(column, "NULL")
}

func (w *Having) IsNotNull(column string) ksql.HavingInterface {
	return w._is(column, "NOT NULL")
}

func (w *Having) Express(raw ksql.ExpressInterface) ksql.HavingInterface {
	w.ops = append(w.ops, &whereOp{expr: raw})
	return w
}

func (h *Having) _by(op string, call func(o ksql.HavingInterface)) ksql.HavingInterface {
	n := NewHaving()
	n.onlyBody = true
	call(n)
	if op == "OR" {
		h.orWheres = append(h.orWheres, n)
		return h
	}

	h.andWheres = append(h.andWheres, n)
	return h
}

func (w *Having) OrHaving(call func(o ksql.HavingInterface)) ksql.HavingInterface {
	return w._by("OR", call)
}

func (w *Having) AndHaving(call func(o ksql.HavingInterface)) ksql.HavingInterface {
	return w._by("AND", call)
}

func (w *Having) _inBy(column string, sub ksql.QueryInterface, op string) ksql.HavingInterface {
	w.ops = append(w.ops, &whereOp{column: column, sub: sub, op: ksql.Op(op)})
	return w
}

func (w *Having) InBy(column string, sub ksql.QueryInterface) ksql.HavingInterface {
	return w._inBy(column, sub, "IN")
}

func (w *Having) NotInBy(column string, sub ksql.QueryInterface) ksql.HavingInterface {
	return w._inBy(column, sub, "NOT IN")
}

func (w *Having) Between(column string, begin, end any) ksql.HavingInterface {
	return w.between("BETWEEN", column, begin, end)
}

func (w *Having) between(op, column string, begin, end any) ksql.HavingInterface {
	w.ops = append(w.ops, &whereOp{column: column, op: ksql.Op(op), values: []any{begin, end}, isBetween: true})
	return w
}

func (w *Having) NotBetween(column string, begin, end any) ksql.HavingInterface {
	return w.between("NOT BETWEEN", column, begin, end)
}

func (w *Having) Empty() bool {
	return len(w.ops) == 0 && len(w.andWheres) == 0 && len(w.orWheres) == 0
}
