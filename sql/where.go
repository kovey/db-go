package sql

import (
	ksql "github.com/kovey/db-go/v3"
)

type Where struct {
	*base
}

func NewWhere() *Where {
	return &Where{base: &base{hasPrepared: false}}
}

func (w *Where) Where(column string, op string, data any) ksql.WhereInterface {
	if w.builder.Len() > 0 {
		w.builder.WriteString(" AND ")
	}

	Column(column, &w.builder)
	w.builder.WriteString(" ")
	w.builder.WriteString(op)
	w.builder.WriteString(" ?")
	w.binds = append(w.binds, data)
	return w
}

func (w *Where) In(column string, data []any) ksql.WhereInterface {
	if w.builder.Len() > 0 {
		w.builder.WriteString(" AND ")
	}

	Column(column, &w.builder)
	w.builder.WriteString(" IN (")
	for i := 0; i < len(data); i++ {
		if i > 0 {
			w.builder.WriteString(",")
		}
		w.builder.WriteString("?")
	}

	w.builder.WriteString(")")
	w.binds = append(w.binds, data...)
	return w
}

func (w *Where) NotIn(column string, data []any) ksql.WhereInterface {
	if w.builder.Len() > 0 {
		w.builder.WriteString(" AND ")
	}

	Column(column, &w.builder)
	w.builder.WriteString(" NOT IN (")
	for i := 0; i < len(data); i++ {
		if i > 0 {
			w.builder.WriteString(",")
		}
		w.builder.WriteString("?")
	}

	w.builder.WriteString(")")
	w.binds = append(w.binds, data...)
	return w
}

func (w *Where) _is(column string, op string) ksql.WhereInterface {
	if w.builder.Len() > 0 {
		w.builder.WriteString(" AND ")
	}
	Column(column, &w.builder)
	w.builder.WriteString(" IS ")
	w.builder.WriteString(op)
	return w
}

func (w *Where) IsNull(column string) ksql.WhereInterface {
	return w._is(column, "NULL")
}

func (w *Where) IsNotNull(column string) ksql.WhereInterface {
	return w._is(column, "NOT NULL")
}

func (w *Where) Express(raw ksql.ExpressInterface) ksql.WhereInterface {
	if w.builder.Len() > 0 {
		w.builder.WriteString(" AND ")
	}

	w.builder.WriteString(raw.Statement())
	w.binds = append(w.binds, raw.Binds()...)
	return w
}

func (w *Where) OrWhere(call func(o ksql.WhereInterface)) ksql.WhereInterface {
	return w._by("OR", call)
}

func (w *Where) _by(op string, call func(o ksql.WhereInterface)) ksql.WhereInterface {
	w.builder.WriteString(" ")
	w.builder.WriteString(op)
	w.builder.WriteString(" (")
	n := NewWhere()
	call(n)
	w.builder.WriteString(n.Prepare())
	w.builder.WriteString(")")
	w.binds = append(w.binds, n.binds...)
	return w
}

func (w *Where) AndWhere(call func(o ksql.WhereInterface)) ksql.WhereInterface {
	return w._by("AND", call)
}

func (w *Where) _inBy(column string, sub ksql.QueryInterface, op string) ksql.WhereInterface {
	if w.builder.Len() > 0 {
		w.builder.WriteString(" AND ")
	}
	Column(column, &w.builder)
	w.builder.WriteString(" ")
	w.builder.WriteString(op)
	w.builder.WriteString(" (")
	w.builder.WriteString(sub.Prepare())
	w.builder.WriteString(")")
	w.binds = append(w.binds, sub.Binds()...)
	return w
}

func (w *Where) InBy(column string, sub ksql.QueryInterface) ksql.WhereInterface {
	return w._inBy(column, sub, "IN")
}

func (w *Where) NotInBy(column string, sub ksql.QueryInterface) ksql.WhereInterface {
	return w._inBy(column, sub, "NOT IN")
}

func (w *Where) Between(column string, begin, end any) ksql.WhereInterface {
	if w.builder.Len() > 0 {
		w.builder.WriteString(" AND ")
	}

	Column(column, &w.builder)
	w.builder.WriteString(" BETWEEN ? ")
	w.builder.WriteString(" AND ? ")
	w.binds = append(w.binds, begin, end)
	return w
}

func (w *Where) Empty() bool {
	return w.builder.Len() == 0
}
