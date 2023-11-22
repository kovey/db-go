package sql

import (
	"fmt"
	"strings"

	"github.com/kovey/pool"
	"github.com/kovey/pool/object"
)

const (
	whereFormat     string = "WHERE (%s)"
	betweenFormat   string = "%s BETWEEN %s AND %s"
	inFormat        string = "%s IN(%s)"
	notInFormat     string = "%s NOT IN(%s)"
	isNullFormat    string = "%s IS NULL"
	isNotNullFormat string = "%s IS NOT NULL"
	where_name             = "Where"
)

func init() {
	pool.DefaultNoCtx(namespace, where_name, func() any {
		return &Where{ObjNoCtx: object.NewObjNoCtx(namespace, where_name)}
	})
}

type Where struct {
	*object.ObjNoCtx
	Fields []string `json:"fields"`
	Binds  []any    `json:"binds"`
}

func NewWhere() *Where {
	return &Where{Fields: make([]string, 0), Binds: make([]any, 0)}
}

func NewWhereBy(ctx object.CtxInterface) *Where {
	return ctx.GetNoCtx(namespace, where_name).(*Where)
}

func (w *Where) Reset() {
	w.Fields = nil
	w.Binds = nil
}

func (w *Where) Eq(field string, value any) {
	w.set(eq, field, value)
}

func (w *Where) set(op string, field string, value any) {
	w.Fields = append(w.Fields, fmt.Sprintf(whereFields, formatValue(field), op))
	w.Binds = append(w.Binds, value)
}

func (w *Where) Neq(field string, value any) {
	w.set(neq, field, value)
}

func (w *Where) Like(field string, value any) {
	w.set(like, field, value)
}

func (w *Where) Between(field string, from any, to any) {
	w.Fields = append(w.Fields, fmt.Sprintf(betweenFormat, formatValue(field), question, question))
	w.Binds = append(w.Binds, from, to)
}

func (w *Where) Gt(field string, value any) {
	w.set(gt, field, value)
}

func (w *Where) Ge(field string, value any) {
	w.set(ge, field, value)
}

func (w *Where) Lt(field string, value any) {
	w.set(lt, field, value)
}

func (w *Where) Le(field string, value any) {
	w.set(le, field, value)
}

func (w *Where) setIn(format string, field string, value []any) {
	placeholders := make([]string, len(value))
	for i := 0; i < len(value); i++ {
		placeholders[i] = question
	}

	w.Fields = append(w.Fields, fmt.Sprintf(format, formatValue(field), strings.Join(placeholders, comma)))
	w.Binds = append(w.Binds, value...)
}

func (w *Where) In(field string, value []any) {
	w.setIn(inFormat, field, value)
}

func (w *Where) NotIn(field string, value []any) {
	w.setIn(notInFormat, field, value)
}

func (w *Where) setNull(format string, field string) {
	w.Fields = append(w.Fields, fmt.Sprintf(format, formatValue(field)))
}

func (w *Where) IsNull(field string) {
	w.setNull(isNullFormat, field)
}

func (w *Where) IsNotNull(field string) {
	w.setNull(isNotNullFormat, field)
}

func (w *Where) Statement(statement string) {
	w.Fields = append(w.Fields, statement)
}

func (w *Where) Args() []any {
	return w.Binds
}

func (w *Where) prepare(op string) string {
	if len(w.Fields) == 0 {
		return emptyStr
	}

	return fmt.Sprintf(whereFormat, strings.Join(w.Fields, op))
}

func (w *Where) Prepare() string {
	return w.prepare(and)
}

func (w *Where) OrPrepare() string {
	return w.prepare(or)
}

func (w *Where) String() string {
	return String(w)
}
