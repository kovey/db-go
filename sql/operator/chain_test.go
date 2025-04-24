package operator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testSelectCall(builder *strings.Builder) {
	builder.WriteString("SELECT")
}

func testSelectColumn(builder *strings.Builder) {
	builder.WriteString(" *")
}

func testSelectFrom(builder *strings.Builder) {
	builder.WriteString(" FROM")
}

func TestChain(t *testing.T) {
	chain := NewChain().Append(testSelectCall, testSelectColumn, testSelectFrom)
	var builder strings.Builder
	chain.Call(&builder)
	assert.Equal(t, "SELECT * FROM", builder.String())
	chain.Call(&builder)
	assert.Equal(t, "SELECT * FROM", builder.String())
	chain.Reset()
	chain.Call(&builder)
	assert.Equal(t, "SELECT * FROMSELECT * FROM", builder.String())
}
