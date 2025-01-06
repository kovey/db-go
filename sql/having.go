package sql

import ksql "github.com/kovey/db-go/v3"

type Having struct {
	*base
}

func NewHaving() *Having {
	return &Having{base: &base{hasPrepared: false}}
}

func (w *Having) Having(column string, op string, data any) ksql.HavingInterface {
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

func (w *Having) In(column string, data []any) ksql.HavingInterface {
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

func (w *Having) NotIn(column string, data []any) ksql.HavingInterface {
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

func (w *Having) _is(column string, op string) ksql.HavingInterface {
	if w.builder.Len() > 0 {
		w.builder.WriteString(" AND ")
	}
	Column(column, &w.builder)
	w.builder.WriteString(" IS ")
	w.builder.WriteString(op)
	return w
}

func (w *Having) IsNull(column string) ksql.HavingInterface {
	return w._is(column, "NULL")
}

func (w *Having) IsNotNull(column string) ksql.HavingInterface {
	return w._is(column, "NOT NULL")
}

func (w *Having) Express(raw ksql.ExpressInterface) ksql.HavingInterface {
	if w.builder.Len() > 0 {
		w.builder.WriteString(" AND ")
	}

	w.builder.WriteString(raw.Statement())
	w.binds = append(w.binds, raw.Binds()...)
	return w
}

func (h *Having) _by(op string, call func(o ksql.HavingInterface)) ksql.HavingInterface {
	h.builder.WriteString(" ")
	h.builder.WriteString(op)
	h.builder.WriteString(" (")
	n := NewHaving()
	call(n)
	h.builder.WriteString(n.Prepare())
	h.builder.WriteString(")")
	h.binds = append(h.binds, n.binds...)

	return h
}

func (w *Having) OrHaving(call func(o ksql.HavingInterface)) ksql.HavingInterface {
	return w._by("OR", call)
}

func (w *Having) AndHaving(call func(o ksql.HavingInterface)) ksql.HavingInterface {
	return w._by("AND", call)
}

func (w *Having) _inBy(column string, sub ksql.QueryInterface, op string) ksql.HavingInterface {
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

func (w *Having) InBy(column string, sub ksql.QueryInterface) ksql.HavingInterface {
	return w._inBy(column, sub, "IN")
}

func (w *Having) NotInBy(column string, sub ksql.QueryInterface) ksql.HavingInterface {
	return w._inBy(column, sub, "NOT IN")
}

func (w *Having) Between(column string, begin, end any) ksql.HavingInterface {
	if w.builder.Len() > 0 {
		w.builder.WriteString(" AND ")
	}

	Column(column, &w.builder)
	w.builder.WriteString(" BETWEEN ? ")
	w.builder.WriteString(" AND ? ")
	w.binds = append(w.binds, begin, end)
	return w
}

func (w *Having) Empty() bool {
	return w.builder.Len() == 0
}
