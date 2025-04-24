package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlterOptions(t *testing.T) {
	options := &AlterOptions{}
	assert.True(t, options.Empty())
	options.Append(NewAddAlterOption(NewAlterOption("METHOD", "PREFIX", "FIELD").IsField()))
	options.Append(NewEditAlterOption(NewAlterOption("METHOD", "PREFIX", "STR").IsStr()))
	options.Append(NewEqAlterOption("KEY", "VALUE"))
	var builder strings.Builder
	options.Build(&builder)
	assert.Equal(t, "ADD METHOD PREFIX `FIELD`, ALTER METHOD PREFIX 'STR', KEY = VALUE", builder.String())
}

func TestDefaultAlterOption(t *testing.T) {
	d := &DefaultAlterOption{}
	d.Option("KEY", "VALUE")
	d.Option("KEY1", "VALUE1")
	var builder strings.Builder
	d.Build(&builder)
	assert.Equal(t, "DEFAULT KEY = VALUE KEY1 = VALUE1", builder.String())
}
