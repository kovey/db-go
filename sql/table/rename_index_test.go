package table

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestRenameIndex(t *testing.T) {
	r := &RenameIndex{}
	r.Old("idx_name").New("n_index_name").Type(ksql.Index_Sub_Type_Index)
	var builder strings.Builder
	r.Build(&builder)
	assert.Equal(t, "RENAME INDEX `idx_name` TO `n_index_name`", builder.String())
}
