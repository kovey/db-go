package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type Trigger struct {
	*base
	name        string
	definer     string
	ifNotExists bool
	time        string
	event       string
	table       string
	orderType   ksql.TriggerOrderType
	other       string
	body        []string
}

func NewTrigger() *Trigger {
	t := &Trigger{base: newBase()}
	t.opChain.Append(t._keyword, t._timeEnvent, t._on, t._order, t._body)
	return t
}

func (t *Trigger) _keyword(builder *strings.Builder) {
	builder.WriteString("CREATE")
	if t.definer != "" {
		builder.WriteString(" DEFINER = ")
		builder.WriteString(t.definer)
	}

	builder.WriteString(" TRIGGER")
	if t.ifNotExists {
		builder.WriteString(" IF NOT EXISTS")
	}

	operator.BuildColumnString(t.name, builder)
}

func (t *Trigger) _timeEnvent(builder *strings.Builder) {
	operator.BuildPureString(t.time, builder)
	operator.BuildPureString(t.event, builder)
}

func (t *Trigger) _on(builder *strings.Builder) {
	builder.WriteString(" ON")
	operator.BuildColumnString(t.table, builder)
	builder.WriteString(" FOR EACH ROW")
}

func (t *Trigger) _order(builder *strings.Builder) {
	if t.orderType != "" {
		builder.WriteString(" ")
		builder.WriteString(string(t.orderType))
		operator.BuildColumnString(t.other, builder)
	}
}

func (t *Trigger) _body(builder *strings.Builder) {
	builder.WriteString(" BEGIN ")
	for _, sql := range t.body {
		builder.WriteString(sql)
		builder.WriteString("; ")
	}
	builder.WriteString("END")
}

func (t *Trigger) Trigger(trigger string) ksql.TriggerInterface {
	t.name = trigger
	return t
}

func (t *Trigger) Definer(definer string) ksql.TriggerInterface {
	t.definer = definer
	return t
}

func (t *Trigger) IfNotExists() ksql.TriggerInterface {
	t.ifNotExists = true
	return t
}

func (t *Trigger) Before() ksql.TriggerInterface {
	t.time = "BEFORE"
	return t
}

func (t *Trigger) After() ksql.TriggerInterface {
	t.time = "AFTER"
	return t
}

func (t *Trigger) Insert() ksql.TriggerInterface {
	t.event = "INSERT"
	return t
}

func (t *Trigger) Update() ksql.TriggerInterface {
	t.event = "UPDATE"
	return t
}

func (t *Trigger) Delete() ksql.TriggerInterface {
	t.event = "DELETE"
	return t
}

func (t *Trigger) On(table string) ksql.TriggerInterface {
	t.table = table
	return t
}

func (t *Trigger) Order(typ ksql.TriggerOrderType, otherTrigger string) ksql.TriggerInterface {
	t.orderType = typ
	t.other = otherTrigger
	return t
}

func (t *Trigger) Body(sql ksql.SqlInterface) ksql.TriggerInterface {
	t.body = append(t.body, sql.Prepare())
	t.binds = append(t.binds, sql.Binds()...)
	return t
}

func (t *Trigger) BodyRaw(sql ksql.ExpressInterface) ksql.TriggerInterface {
	t.body = append(t.body, sql.Statement())
	t.binds = append(t.binds, sql.Binds()...)
	return t
}
