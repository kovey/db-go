package sql

import (
	"fmt"
	"strings"

	"github.com/kovey/pool"
	"github.com/kovey/pool/object"
)

const (
	havingFormat = "HAVING (%s)"
	whereFields  = "%s %s ?"
	eq           = "="
	neq          = "<>"
	like         = "LIKE"
	gt           = ">"
	ge           = ">="
	lt           = "<"
	le           = "<="
	question     = "?"
	or           = " OR "
	having_name  = "Having"
)

func init() {
	pool.DefaultNoCtx(namespace, having_name, func() any {
		return &Having{ObjNoCtx: object.NewObjNoCtx(namespace, having_name)}
	})
}

type Having struct {
	*object.ObjNoCtx
	fields []string
	args   []any
}

func NewHaving() *Having {
	return &Having{fields: make([]string, 0), args: make([]any, 0)}
}

func NewHavingBy(ctx object.CtxInterface) *Having {
	return ctx.GetNoCtx(namespace, having_name).(*Having)
}

func (h *Having) Reset() {
	h.fields = nil
	h.args = nil
}

func (w *Having) Eq(field string, value any) *Having {
	return w.set(eq, field, value)
}

func (w *Having) set(op string, field string, value any) *Having {
	w.fields = append(w.fields, fmt.Sprintf(whereFields, formatValue(field), op))
	w.args = append(w.args, value)
	return w
}

func (w *Having) Neq(field string, value any) *Having {
	return w.set(neq, field, value)
}

func (w *Having) Like(field string, value any) *Having {
	return w.set(like, field, value)
}

func (w *Having) Between(field string, from any, to any) *Having {
	w.fields = append(w.fields, fmt.Sprintf(betweenFormat, formatValue(field), question, question))
	w.args = append(w.args, from, to)
	return w
}

func (w *Having) Gt(field string, value any) *Having {
	return w.set(gt, field, value)
}

func (w *Having) Ge(field string, value any) *Having {
	return w.set(ge, field, value)
}

func (w *Having) Lt(field string, value any) *Having {
	return w.set(lt, field, value)
}

func (w *Having) Le(field string, value any) *Having {
	return w.set(le, field, value)
}

func (w *Having) setIn(format string, field string, value []any) *Having {
	placeholders := make([]string, len(value))
	for i := 0; i < len(value); i++ {
		placeholders[i] = question
	}

	w.fields = append(w.fields, fmt.Sprintf(format, formatValue(field), strings.Join(placeholders, comma)))
	w.args = append(w.args, value...)
	return w
}

func (w *Having) In(field string, value []any) *Having {
	return w.setIn(inFormat, field, value)
}

func (w *Having) NotIn(field string, value []any) *Having {
	return w.setIn(notInFormat, field, value)
}

func (w *Having) setNull(format string, field string) *Having {
	w.fields = append(w.fields, fmt.Sprintf(format, formatValue(field)))
	return w
}

func (w *Having) IsNull(field string) *Having {
	return w.setNull(isNullFormat, field)
}

func (w *Having) IsNotNull(field string) *Having {
	return w.setNull(isNotNullFormat, field)
}

func (w *Having) Statement(statement string) *Having {
	w.fields = append(w.fields, statement)
	return w
}

func (w *Having) Args() []any {
	return w.args
}

func (w *Having) prepare(op string) string {
	return fmt.Sprintf(havingFormat, strings.Join(w.fields, op))
}

func (w *Having) Prepare() string {
	return w.prepare(and)
}

func (w *Having) OrPrepare() string {
	return w.prepare(or)
}

func (w *Having) String() string {
	return String(w)
}
