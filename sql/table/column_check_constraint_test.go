package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnCheckConstraint(t *testing.T) {
	c := NewColumnCheckConstraint().Check("SUM(amount)").Constraint("test").Enforced()
	var builder strings.Builder
	c.Build(&builder)
	assert.Equal(t, "CONSTRAINT test CHECK (SUM(amount)) ENFORCED", builder.String())
}

func TestColumnCheckConstraintNotEnforced(t *testing.T) {
	c := NewColumnCheckConstraint().Check("SUM(amount)").NotEnforced()
	var builder strings.Builder
	c.Build(&builder)
	assert.Equal(t, "CHECK (SUM(amount)) NOT ENFORCED", builder.String())
}
