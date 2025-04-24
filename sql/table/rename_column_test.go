package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenameColumn(t *testing.T) {
	r := &RenameColumn{}
	r.Old("old_name").New("new_name")
	var builder strings.Builder
	r.Build(&builder)
	assert.Equal(t, "RENAME COLUMN `old_name` TO `new_name`", builder.String())
}
