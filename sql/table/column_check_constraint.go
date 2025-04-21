package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type ColumnCheckConstraint struct {
	symbol     string
	expr       string
	enforced   string
	opChain    *operator.Chain
	constraint string
}

func NewColumnCheckConstraint() *ColumnCheckConstraint {
	c := &ColumnCheckConstraint{opChain: operator.NewChain()}
	c.opChain.Append(c._constraint, c._check, c._enforced)
	return c
}

func (c *ColumnCheckConstraint) _constraint(builder *strings.Builder) {
	if c.constraint == "" {
		return
	}

	builder.WriteString(c.constraint)
	builder.WriteString(" ")
	builder.WriteString(c.symbol)
}

func (c *ColumnCheckConstraint) _check(builder *strings.Builder) {
	if c.expr == "" {
		return
	}

	if c.constraint == "" {
		builder.WriteString("CHECK (")
	} else {

		builder.WriteString(" CHECK (")
	}
	builder.WriteString(c.expr)
	builder.WriteString(")")
}

func (c *ColumnCheckConstraint) _enforced(builder *strings.Builder) {
	if c.enforced == "" {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(c.enforced)
}

func (c *ColumnCheckConstraint) Constraint(symbol string) ksql.ColumnCheckConstraintInterface {
	c.constraint = "CONSTRAINT"
	c.symbol = symbol
	return c
}

func (c *ColumnCheckConstraint) Check(expr string) ksql.ColumnCheckConstraintInterface {
	c.expr = expr
	return c
}

func (c *ColumnCheckConstraint) Enforced() ksql.ColumnCheckConstraintInterface {
	c.enforced = "ENFORCED"
	return c
}

func (c *ColumnCheckConstraint) NotEnforced() ksql.ColumnCheckConstraintInterface {
	c.enforced = "NOT ENFORCED"
	return c
}

func (c *ColumnCheckConstraint) Build(builder *strings.Builder) {
	c.opChain.Call(builder)
}
