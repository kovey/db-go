package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type Default struct {
	IsKeyword bool
	Value     string
	IsByte    bool
	IsExpr    bool
}

func (d *Default) Build(builder *strings.Builder) {
	builder.WriteString(" DEFAULT ")
	if d.IsKeyword {
		builder.WriteString(d.Value)
		return
	}

	if d.IsExpr {
		builder.WriteString("(")
		builder.WriteString(d.Value)
		builder.WriteString(")")
		return
	}

	if d.IsByte {
		builder.WriteString("b")
		operator.Quote(d.Value, builder)
		return
	}

	operator.Quote(d.Value, builder)
}

type Column struct {
	name          string
	t             *ColumnType
	def           *Default
	comment       string
	isUnsigned    bool
	opChain       *operator.Chain
	null          string
	index         string
	unique        string
	autoInc       string
	visible       string
	collate       string
	format        ksql.ColumnFormat
	engineAttr    string
	secEngineAttr string
	storage       ksql.ColumnStorage
	reference     *ColumnReference
	check         *ColumnCheckConstraint
}

func NewColumn(name string, t *ColumnType) *Column {
	c := &Column{name: name, t: t, opChain: operator.NewChain()}
	c.opChain.Append(c._nameInfo, c._null, c._default, c._pureString, c._format, c._attr, c._reference)
	return c
}

func (c *Column) _nameInfo(builder *strings.Builder) {
	operator.Backtick(c.name, builder)
	c.t.Build(builder)
	if c.isUnsigned && c.t.IsNumeric() {
		builder.WriteString(" UNSIGNED")
	}
}

func (c *Column) _null(builder *strings.Builder) {
	if c.autoInc != "" && c.t.IsInteger() {
		return
	}

	operator.BuildPureString(c.null, builder)
}

func (c *Column) _default(builder *strings.Builder) {
	if c.def == nil {
		return
	}

	c.def.Build(builder)
}

func (c *Column) _pureString(builder *strings.Builder) {
	operator.BuildPureString(c.visible, builder)
	if c.t.IsInteger() {
		operator.BuildPureString(c.autoInc, builder)
	}
	operator.BuildPureString(c.unique, builder)
	operator.BuildPureString(c.index, builder)
	if c.comment != "" {
		builder.WriteString(" COMMENT")
		operator.BuildQuoteString(c.comment, builder)
	}
	operator.BuildPureString(c.collate, builder)
}

func (c *Column) _format(builder *strings.Builder) {
	if c.format == "" {
		return
	}

	builder.WriteString(" COLUMN_FORMAT")
	operator.BuildPureString(string(c.format), builder)
}

func (c *Column) _attr(builder *strings.Builder) {
	operator.BuildQuoteString(c.engineAttr, builder)
	operator.BuildQuoteString(c.secEngineAttr, builder)
	operator.BuildPureString(string(c.storage), builder)
}

func (c *Column) _reference(builder *strings.Builder) {
	if c.reference != nil {
		c.reference.Build(builder)
	}

	if c.check != nil {
		c.check.Build(builder)
	}
}

func (c *Column) Build(builder *strings.Builder) {
	c.opChain.Call(builder)
}

func (c *Column) NotNullable() ksql.ColumnInterface {
	c.null = "NOT NULL"
	return c
}

func (c *Column) DefaultExpress(expr string) ksql.ColumnInterface {
	c.def = &Default{Value: expr, IsExpr: true}
	return c
}

func (c *Column) DefaultBit(bit string) ksql.ColumnInterface {
	c.def = &Default{Value: bit, IsByte: true}
	return c
}

func (c *Column) Visible() ksql.ColumnInterface {
	c.visible = "VISIBLE"
	return c
}

func (c *Column) Invisible() ksql.ColumnInterface {
	c.visible = "INVISIBLE"
	return c
}

func (c *Column) Unique() ksql.ColumnInterface {
	c.unique = "UNIQUE KEY"
	return c
}

func (c *Column) Index() ksql.ColumnInterface {
	c.index = "KEY"
	return c
}

func (c *Column) Primary() ksql.ColumnInterface {
	c.index = "PRIMARY KEY"
	return c
}

func (c *Column) Collate(collate string) ksql.ColumnInterface {
	c.collate = collate
	return c
}

func (c *Column) Format(format ksql.ColumnFormat) ksql.ColumnInterface {
	c.format = format
	return c
}

func (c *Column) EngineAttribute(engineAttr string) ksql.ColumnInterface {
	c.engineAttr = engineAttr
	return c
}

func (c *Column) SecondaryEngineAttribute(secEngineAttr string) ksql.ColumnInterface {
	c.secEngineAttr = secEngineAttr
	return c
}

func (c *Column) Storage(storage ksql.ColumnStorage) ksql.ColumnInterface {
	c.storage = storage
	return c
}

func (c *Column) Reference(table string) ksql.ColumnReferenceInterface {
	c.reference = NewColumnReference(table)
	return c.reference
}

func (c *Column) CheckConstraint() ksql.ColumnCheckConstraintInterface {
	c.check = NewColumnCheckConstraint()
	return c.check
}

func (c *Column) Nullable() ksql.ColumnInterface {
	c.null = "NULL"
	return c
}

func (c *Column) AutoIncrement() ksql.ColumnInterface {
	c.autoInc = "AUTO_INCREMENT"
	return c
}

func (c *Column) Unsigned() ksql.ColumnInterface {
	c.isUnsigned = true
	return c
}

func (c *Column) UseCurrent() ksql.ColumnInterface {
	c.Default(ksql.CURRENT_TIMESTAMP)
	return c
}

func (c *Column) UseCurrentOnUpdate() ksql.ColumnInterface {
	c.Default(ksql.CURRENT_TIMESTAMP_ON_UPDATE_CURRENT_TIMESTAMP)
	return c
}

func (c *Column) Default(value string) ksql.ColumnInterface {
	if c.def != nil {
		return c
	}

	c.def = &Default{Value: value, IsKeyword: ksql.IsDefaultKeyword(value)}
	return c
}

func (c *Column) Comment(comment string) ksql.ColumnInterface {
	c.comment = comment
	return c
}
