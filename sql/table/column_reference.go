package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type onItem struct {
	op     ksql.ReferenceOnOpt
	option ksql.ReferenceOption
}

func (o *onItem) Build(builder *strings.Builder) {
	builder.WriteString(" ON ")
	builder.WriteString(string(o.op))
	builder.WriteString(" ")
	builder.WriteString(string(o.option))
}

type ColumnReference struct {
	table   string
	columns *IndexColumns
	match   ksql.ReferenceMatch
	ons     []*onItem
	opChain *operator.Chain
}

func NewColumnReference(table string) *ColumnReference {
	c := &ColumnReference{table: table, opChain: operator.NewChain(), columns: &IndexColumns{}}
	c.opChain.Append(c._table, c._macth, c._on)
	return c
}

func (c *ColumnReference) _table(builder *strings.Builder) {
	builder.WriteString(" REFERENCES ")
	operator.Backtick(c.table, builder)
	c.columns.Build(builder)
}

func (c *ColumnReference) _macth(builder *strings.Builder) {
	if c.match == "" {
		return
	}

	builder.WriteString(" MATCH ")
	builder.WriteString(string(c.match))
}

func (c *ColumnReference) _on(builder *strings.Builder) {
	if c.ons == nil {
		return
	}

	for _, on := range c.ons {
		on.Build(builder)
	}
}

func (c *ColumnReference) Column(name string, length int, order ksql.Order) ksql.ColumnReferenceInterface {
	c.columns.Append(&IndexColumn{Name: name, Length: length, Order: order, Type: Index_Column_Type_Name})
	return c
}

func (c *ColumnReference) Express(express string, order ksql.Order) ksql.ColumnReferenceInterface {
	c.columns.Append(&IndexColumn{Name: express, Order: order, Type: Index_Column_Type_Expr})
	return c
}

func (c *ColumnReference) Match(match ksql.ReferenceMatch) ksql.ColumnReferenceInterface {
	c.match = match
	return c
}

func (c *ColumnReference) On(op ksql.ReferenceOnOpt, option ksql.ReferenceOption) ksql.ColumnReferenceInterface {
	c.ons = append(c.ons, &onItem{op: op, option: option})
	return c
}

func (c *ColumnReference) Build(builder *strings.Builder) {
	c.opChain.Call(builder)
}
