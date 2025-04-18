package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type Call struct {
	*base
	spName string
	params []string
}

func NewCall() *Call {
	c := &Call{base: newBase()}
	c.opChain.Append(c._build)
	return c
}

func (c *Call) _build(builder *strings.Builder) {
	builder.WriteString("CALL")
	operator.BuildColumnString(c.spName, builder)
	if len(c.params) > 0 {
		builder.WriteString(" (")
		for index, param := range c.params {
			if index > 0 {
				builder.WriteString(", ")
			}
			if !strings.HasPrefix(param, "@") {
				builder.WriteString("@")
			}
			builder.WriteString(param)
		}
		builder.WriteString(")")
	}
}

func (c *Call) Call(spName string) ksql.CallInterface {
	c.spName = spName
	return c
}

func (c *Call) Params(params ...string) ksql.CallInterface {
	c.params = append(c.params, params...)
	return c
}
