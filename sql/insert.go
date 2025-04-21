package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type assignment struct {
	column  string
	value   string
	isField bool
	isData  bool
	data    any
	expr    ksql.ExpressInterface
	op      string
}

func (a *assignment) binds() []any {
	if a.expr == nil {
		if a.isData {
			return []any{a.data}
		}
		return nil
	}

	return a.expr.Binds()
}

func (a *assignment) Build(builder *strings.Builder) {
	if a.expr != nil {
		operator.BuildPureString(a.expr.Statement(), builder)
		return
	}

	operator.BuildColumnString(a.column, builder)
	builder.WriteString(" =")
	if a.isField {
		operator.BuildColumnString(a.value, builder)
		return
	}

	if a.isData {
		if a.op != "" {
			operator.BuildColumnString(a.column, builder)
			operator.BuildPureString(a.op, builder)
		}

		builder.WriteString(" ?")
		return
	}

	operator.BuildPureString(a.value, builder)
}

type assignments struct {
	asses []*assignment
	binds []any
}

func (a *assignments) Append(ass *assignment) {
	a.asses = append(a.asses, ass)
}

func (a *assignments) Build(builder *strings.Builder) {
	for index, ass := range a.asses {
		if index > 0 {
			builder.WriteString(",")
		}

		ass.Build(builder)
		a.binds = append(a.binds, ass.binds()...)
	}
}

func (a *assignments) Empty() bool {
	return len(a.asses) == 0
}

type Insert struct {
	*base
	columns     []string
	from        ksql.QueryInterface
	fromTable   string
	rowAs       string
	columnsAs   []string
	sets        *assignments
	onUpdates   *assignments
	table       string
	lowPriority string
	ignore      string
	partitions  []string
}

func NewInsert() *Insert {
	i := &Insert{base: newBase(), sets: &assignments{}, onUpdates: &assignments{}}
	i.opChain.Append(i._keyword, i._name, i._set, i._columns, i._from, i._values, i._as, i._on)
	return i
}

func (i *Insert) _keyword(builder *strings.Builder) {
	builder.WriteString("INSERT")
	operator.BuildPureString(i.lowPriority, builder)
	operator.BuildPureString(i.ignore, builder)
	builder.WriteString(" INTO")
}

func (i *Insert) _name(builder *strings.Builder) {
	operator.BuildColumnString(i.table, builder)
	for index, partition := range i.partitions {
		if index > 0 {
			builder.WriteString(",")
		}
		operator.BuildColumnString(partition, builder)
	}
}

func (i *Insert) _set(builder *strings.Builder) {
	if i.sets.Empty() {
		return
	}

	builder.WriteString(" SET ")
	i.sets.Build(builder)
	i.binds = append(i.binds, i.sets.binds...)
}

func (i *Insert) _columns(builder *strings.Builder) {
	if len(i.columns) == 0 || !i.sets.Empty() {
		return
	}

	builder.WriteString(" (")
	for index, column := range i.columns {
		if index > 0 {
			builder.WriteString(", ")
		}
		operator.Column(column, builder)
	}
	builder.WriteString(")")
}

func (i *Insert) _from(builder *strings.Builder) {
	if i.from != nil {
		operator.BuildPureString(i.from.Prepare(), builder)
		i.binds = append(i.binds, i.from.Binds()...)
		return
	}

	if i.fromTable != "" {
		builder.WriteString(" TABLE")
		operator.BuildColumnString(i.fromTable, builder)
	}
}

func (i *Insert) _values(builder *strings.Builder) {
	count := len(i.columns)
	if count == 0 || !i.sets.Empty() || i.from != nil || i.fromTable != "" {
		return
	}

	builder.WriteString(" VALUES")
	index := 0
	total := len(i.binds)
	for {
		if index >= total {
			break
		}

		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(" (")
		for i := 0; i < count; i++ {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString("?")
			index++
		}
		builder.WriteString(")")
	}
}

func (i *Insert) _as(builder *strings.Builder) {
	if i.rowAs == "" {
		return
	}

	builder.WriteString(" AS")
	if len(i.columnsAs) > 0 {
		builder.WriteString(" ")
	}

	for index, as := range i.columnsAs {
		if index > 0 {
			builder.WriteString(",")
		}

		operator.BuildColumnString(as, builder)
	}
}

func (i *Insert) _on(builder *strings.Builder) {
	if i.onUpdates.Empty() {
		return
	}

	builder.WriteString(" ON DUPLICATE KEY UPDATE ")
	i.onUpdates.Build(builder)
	i.binds = append(i.binds, i.onUpdates.binds...)
}

func (i *Insert) Table(table string) ksql.InsertInterface {
	i.table = table
	return i
}

func (i *Insert) Add(column string, data any) ksql.InsertInterface {
	i.columns = append(i.columns, column)
	i.binds = append(i.binds, data)
	return i
}

func (i *Insert) Columns(columns ...string) ksql.InsertInterface {
	i.columns = columns
	return i
}

func (i *Insert) Values(datas ...any) ksql.InsertInterface {
	i.binds = append(i.binds, datas...)
	return i
}

func (i *Insert) From(query ksql.QueryInterface) ksql.InsertInterface {
	i.from = query
	i.fromTable = ""
	return i
}

func (i *Insert) Set(column, value string) ksql.InsertInterface {
	i.sets.Append(&assignment{column: column, value: value})
	return i
}

func (i *Insert) SetColumn(column, otherColumn string) ksql.InsertInterface {
	i.sets.Append(&assignment{column: column, value: otherColumn, isField: true})
	return i
}

func (i *Insert) SetExpress(expr ksql.ExpressInterface) ksql.InsertInterface {
	i.sets.Append(&assignment{expr: expr})
	return i
}

func (i *Insert) LowPriority() ksql.InsertInterface {
	i.lowPriority = "LOW_PRIORITY"
	return i
}

func (i *Insert) Delayed() ksql.InsertInterface {
	i.lowPriority = "DELAYED"
	return i
}

func (i *Insert) HighPriority() ksql.InsertInterface {
	i.lowPriority = "HIGH_PRIORITY"
	return i
}

func (i *Insert) Ignore() ksql.InsertInterface {
	i.ignore = "IGNORE"
	return i
}

func (i *Insert) Partitions(names ...string) ksql.InsertInterface {
	i.partitions = append(i.partitions, names...)
	return i
}

func (i *Insert) As(rowAlias string, colAlias ...string) ksql.InsertInterface {
	i.rowAs = rowAlias
	i.columnsAs = append(i.columnsAs, colAlias...)
	return i
}

func (i *Insert) OnDuplicateKeyUpdate(column, valueOrColumn string) ksql.InsertInterface {
	i.onUpdates.Append(&assignment{column: column, value: valueOrColumn})
	return i
}

func (i *Insert) OnDuplicateKeyUpdateExpress(expr ksql.ExpressInterface) ksql.InsertInterface {
	i.onUpdates.Append(&assignment{expr: expr})
	return i
}

func (i *Insert) OnDuplicateKeyUpdateColumn(column, otherColumn string) ksql.InsertInterface {
	i.onUpdates.Append(&assignment{column: column, value: otherColumn, isField: true})
	return i
}

func (i *Insert) FromTable(table string) ksql.InsertInterface {
	i.fromTable = table
	i.from = nil
	return i
}
