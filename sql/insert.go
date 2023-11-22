package sql

import (
	"fmt"
	"strings"

	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/pool"
	"github.com/kovey/pool/object"
)

const (
	format   = "INSERT INTO %s (%s) VALUES (%s)"
	ins_name = "Insert"
)

func init() {
	pool.DefaultNoCtx(namespace, ins_name, func() any {
		return &Insert{ObjNoCtx: object.NewObjNoCtx(namespace, ins_name), data: meta.NewData(), placeholder: make(map[string]string)}
	})
}

type Insert struct {
	*object.ObjNoCtx
	data        meta.Data
	table       string
	placeholder map[string]string
	args        []any
	fields      []string
}

func NewInsert(table string) *Insert {
	return &Insert{table: table, data: meta.NewData(), placeholder: make(map[string]string)}
}

func NewInsertBy(ctx object.CtxInterface, table string) *Insert {
	obj := ctx.GetNoCtx(namespace, ins_name).(*Insert)
	obj.table = table
	return obj
}

func (i *Insert) Reset() {
	i.data = meta.NewData()
	i.table = emptyStr
	i.placeholder = make(map[string]string)
	i.args = nil
	i.fields = nil
}

func (i *Insert) Set(field string, value any) *Insert {
	i.data.Add(field, value)
	i.placeholder[field] = question
	return i
}

func (i *Insert) Args() []any {
	return i.args
}

func (i *Insert) Prepare() string {
	return fmt.Sprintf(format, formatValue(i.table), strings.Join(i.getFields(), comma), strings.Join(i.getPlaceholder(), comma))
}

func (i *Insert) getPlaceholder() []string {
	placeholders := make([]string, len(i.placeholder))
	index := 0
	for _, v := range i.placeholder {
		placeholders[index] = v
		index++
	}

	return placeholders
}

func (i *Insert) getFields() []string {
	fields := make([]string, len(i.data))
	i.args = make([]any, len(i.data))
	i.fields = make([]string, len(i.data))
	index := 0
	for field, val := range i.data {
		i.fields[index] = field
		fields[index] = formatValue(field)
		i.args[index] = val
		index++
	}

	return fields
}

func (i *Insert) String() string {
	return String(i)
}

func (i *Insert) ParseValue(fields []string) {
	i.args = make([]any, len(fields))
	for index, field := range fields {
		i.args[index] = i.data[field]
	}
}

func (i *Insert) Fields() []string {
	return i.fields
}
