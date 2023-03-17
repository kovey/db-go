package sql

import (
	"fmt"
	"strings"
)

const (
	whereFormat     string = "WHERE (%s)"
	betweenFormat   string = "%s BETWEEN %s AND %s"
	inFormat        string = "%s IN(%s)"
	notInFormat     string = "%s NOT IN(%s)"
	isNullFormat    string = "%s IS NULL"
	isNotNullFormat string = "%s IS NOT NULL"
)

type Where struct {
	Fields []string `json:"fields"`
	Binds  []any    `json:"binds"`
}

func NewWhere() *Where {
	return &Where{Fields: make([]string, 0), Binds: make([]any, 0)}
}

func (w *Where) Eq(field string, value any) WhereInterface {
	return w.set("=", field, value)
}

func (w *Where) set(op string, field string, value any) WhereInterface {
	w.Fields = append(w.Fields, fmt.Sprintf("%s %s ?", formatValue(field), op))
	w.Binds = append(w.Binds, value)
	return w
}

func (w *Where) Neq(field string, value any) WhereInterface {
	return w.set("<>", field, value)
}

func (w *Where) Like(field string, value any) WhereInterface {
	return w.set("LIKE", field, value)
}

func (w *Where) Between(field string, from any, to any) WhereInterface {
	w.Fields = append(w.Fields, fmt.Sprintf(betweenFormat, formatValue(field), "?", "?"))
	w.Binds = append(w.Binds, from, to)
	return w
}

func (w *Where) Gt(field string, value any) WhereInterface {
	return w.set(">", field, value)
}

func (w *Where) Ge(field string, value any) WhereInterface {
	return w.set(">=", field, value)
}

func (w *Where) Lt(field string, value any) WhereInterface {
	return w.set("<", field, value)
}

func (w *Where) Le(field string, value any) WhereInterface {
	return w.set("<=", field, value)
}

func (w *Where) setIn(format string, field string, value []any) WhereInterface {
	placeholders := make([]string, len(value))
	for i := 0; i < len(value); i++ {
		placeholders[i] = "?"
	}

	w.Fields = append(w.Fields, fmt.Sprintf(format, formatValue(field), strings.Join(placeholders, ",")))
	w.Binds = append(w.Binds, value...)
	return w
}

func (w *Where) In(field string, value []any) WhereInterface {
	return w.setIn(inFormat, field, value)
}

func (w *Where) NotIn(field string, value []any) WhereInterface {
	return w.setIn(notInFormat, field, value)
}

func (w *Where) setNull(format string, field string) WhereInterface {
	w.Fields = append(w.Fields, fmt.Sprintf(format, formatValue(field)))
	return w
}

func (w *Where) IsNull(field string) WhereInterface {
	return w.setNull(isNullFormat, field)
}

func (w *Where) IsNotNull(field string) WhereInterface {
	return w.setNull(isNotNullFormat, field)
}

func (w *Where) Statement(statement string) WhereInterface {
	w.Fields = append(w.Fields, statement)
	return w
}

func (w *Where) Args() []any {
	return w.Binds
}

func (w *Where) prepare(op string) string {
	if len(w.Fields) == 0 {
		return ""
	}

	return fmt.Sprintf(whereFormat, strings.Join(w.Fields, op))
}

func (w *Where) Prepare() string {
	return w.prepare(" AND ")
}

func (w *Where) OrPrepare() string {
	return w.prepare(" OR ")
}

func (w *Where) String() string {
	return String(w)
}
