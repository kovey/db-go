package sql

import (
	"fmt"
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type whereOp struct {
	column     string
	op         ksql.Op
	value      any
	isArr      bool
	isConst    bool
	constValue string
	values     []any
	expr       ksql.ExpressInterface
	sub        ksql.QueryInterface
	isBetween  bool
}

func (w *whereOp) binds() []any {
	if w.expr != nil {
		return w.expr.Binds()
	}

	if w.sub != nil {
		return w.sub.Binds()
	}

	if w.isConst {
		return nil
	}

	if w.isArr || w.isBetween {
		return w.values
	}

	return []any{w.value}
}

func (w *whereOp) Build(builder *strings.Builder) {
	if w.expr != nil {
		builder.WriteString(w.expr.Statement())
		return
	}

	operator.Column(w.column, builder)
	operator.BuildPureString(string(w.op), builder)
	if w.sub != nil {
		builder.WriteString(" (")
		builder.WriteString(w.sub.Prepare())
		builder.WriteString(")")
		return
	}

	if w.isConst {
		operator.BuildPureString(w.constValue, builder)
		return
	}

	if w.isBetween {
		operator.BuildPureString("?", builder)
		operator.BuildPureString("AND", builder)
		operator.BuildPureString("?", builder)
		return
	}

	if !w.isArr {
		operator.BuildPureString("?", builder)
		return
	}

	builder.WriteString(" (")
	for index := range w.values {
		if index > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("?")
	}
	builder.WriteString(")")
}

type Where struct {
	ops       []*whereOp
	orWheres  []*Where
	andWheres []*Where
	opChain   *operator.Chain
	onlyBody  bool
	binds     []any
}

func NewWhere() *Where {
	w := &Where{opChain: operator.NewChain()}
	w.opChain.Append(w._keyword, w._ops, w._andWheres, w._orWheres)
	return w
}

func (w *Where) _keyword(builder *strings.Builder) {
	if w.onlyBody {
		return
	}

	builder.WriteString("WHERE")
}

func (w *Where) _ops(builder *strings.Builder) {
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

func (w *Where) _andWheres(builder *strings.Builder) {
	if len(w.ops) == 0 {
		w._where("", "AND", w.andWheres, builder)
		return
	}

	w._where("AND", "AND", w.andWheres, builder)
}

func (w *Where) _orWheres(builder *strings.Builder) {
	if len(w.ops) == 0 && len(w.andWheres) == 0 {
		w._where("", "OR", w.orWheres, builder)
		return
	}

	w._where("OR", "OR", w.orWheres, builder)
}

func (w *Where) _where(prefix, op string, wheres []*Where, builder *strings.Builder) {
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

func (w *Where) Build(builder *strings.Builder) {
	w.opChain.Call(builder)
}

func (w *Where) Where(column string, op ksql.Op, data any) ksql.WhereInterface {
	if !ksql.SupportOp(op) {
		panic(fmt.Sprintf("op %s not support", op))
	}
	w.ops = append(w.ops, &whereOp{column: column, op: op, value: data})
	return w
}

func (w *Where) In(column string, data []any) ksql.WhereInterface {
	w.ops = append(w.ops, &whereOp{column: column, op: "IN", values: data, isArr: true})
	return w
}

func (w *Where) NotIn(column string, data []any) ksql.WhereInterface {
	w.ops = append(w.ops, &whereOp{column: column, op: "NOT IN", values: data, isArr: true})
	return w
}

func (w *Where) _is(column string, op string) ksql.WhereInterface {
	w.ops = append(w.ops, &whereOp{column: column, op: "IS", constValue: op, isConst: true})
	return w
}

func (w *Where) IsNull(column string) ksql.WhereInterface {
	return w._is(column, "NULL")
}

func (w *Where) IsNotNull(column string) ksql.WhereInterface {
	return w._is(column, "NOT NULL")
}

func (w *Where) Express(raw ksql.ExpressInterface) ksql.WhereInterface {
	w.ops = append(w.ops, &whereOp{expr: raw})
	return w
}

func (w *Where) OrWhere(call func(o ksql.WhereInterface)) ksql.WhereInterface {
	return w._by("OR", call)
}

func (w *Where) _by(op string, call func(o ksql.WhereInterface)) ksql.WhereInterface {
	n := NewWhere()
	n.onlyBody = true
	call(n)
	if op == "OR" {
		w.orWheres = append(w.orWheres, n)
		return w
	}

	w.andWheres = append(w.andWheres, n)
	return w
}

func (w *Where) AndWhere(call func(o ksql.WhereInterface)) ksql.WhereInterface {
	return w._by("AND", call)
}

func (w *Where) _inBy(column string, sub ksql.QueryInterface, op string) ksql.WhereInterface {
	w.ops = append(w.ops, &whereOp{column: column, sub: sub, op: ksql.Op(op)})
	return w
}

func (w *Where) InBy(column string, sub ksql.QueryInterface) ksql.WhereInterface {
	return w._inBy(column, sub, "IN")
}

func (w *Where) NotInBy(column string, sub ksql.QueryInterface) ksql.WhereInterface {
	return w._inBy(column, sub, "NOT IN")
}

func (w *Where) Between(column string, begin, end any) ksql.WhereInterface {
	return w.between("BETWEEN", column, begin, end)
}

func (w *Where) NotBetween(column string, begin, end any) ksql.WhereInterface {
	return w.between("NOT BETWEEN", column, begin, end)
}

func (w *Where) between(op, column string, begin, end any) ksql.WhereInterface {
	w.ops = append(w.ops, &whereOp{column: column, op: ksql.Op(op), values: []any{begin, end}, isBetween: true})
	return w
}

func (w *Where) Empty() bool {
	return len(w.ops) == 0 && len(w.andWheres) == 0 && len(w.orWheres) == 0
}

func (w *Where) Binds() []any {
	return w.binds
}
