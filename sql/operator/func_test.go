package operator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunc(t *testing.T) {
	var builder strings.Builder
	Quote("kovey", &builder)
	assert.Equal(t, "'kovey'", builder.String())
	builder.Reset()
	Backtick("kovey", &builder)
	assert.Equal(t, "`kovey`", builder.String())
	builder.Reset()
	Backtick("(kovey)", &builder)
	assert.Equal(t, "(kovey)", builder.String())
	builder.Reset()
	Column("dt.column", &builder)
	assert.Equal(t, "`dt`.`column`", builder.String())
	builder.Reset()
	Column("(dt.column)", &builder)
	assert.Equal(t, "(dt.column)", builder.String())
	builder.Reset()
	BuildPureString("kovey", &builder)
	assert.Equal(t, " kovey", builder.String())
	builder.Reset()
	BuildQuoteString("kovey", &builder)
	assert.Equal(t, " 'kovey'", builder.String())
	builder.Reset()
	BuildBacktickString("kovey", &builder)
	assert.Equal(t, " `kovey`", builder.String())
	builder.Reset()
	BuildColumnString("dt.kovey", &builder)
	assert.Equal(t, " `dt`.`kovey`", builder.String())
	BuildColumnString("", &builder)
	assert.Equal(t, " `dt`.`kovey`", builder.String())
	BuildPureString("", &builder)
	assert.Equal(t, " `dt`.`kovey`", builder.String())
	BuildBacktickString("", &builder)
	assert.Equal(t, " `dt`.`kovey`", builder.String())
	BuildQuoteString("", &builder)
	assert.Equal(t, " `dt`.`kovey`", builder.String())
}
