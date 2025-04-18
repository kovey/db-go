package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartitionOptionsHash(t *testing.T) {
	p := &PartitionOptionsHash{expr: "ROUND(amount, 2)"}
	var builder strings.Builder
	p.Build(&builder)
	assert.Equal(t, " HASH(ROUND(amount, 2))", builder.String())
	builder.Reset()
	p.Linear()
	p.Build(&builder)
	assert.Equal(t, " LINEAR HASH(ROUND(amount, 2))", builder.String())
}
