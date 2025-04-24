package table

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestAddIndex(t *testing.T) {
	a := NewAddIndex("idx_name")
	a.Index().Algorithm(ksql.Index_Alg_BTree).BlockSize("10M").Columns("user_name", "nickname")
	var builder strings.Builder
	a.Build(&builder)
	assert.Equal(t, "ADD INDEX `idx_name` USING 'BTREE' (`user_name`, `nickname`) KEY_BLOCK_SIZE = 10M", builder.String())
}
