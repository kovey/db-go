package table

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestIndexOptionFull(t *testing.T) {
	i := NewIndexOption().Invisible().Algorithm(ksql.Index_Alg_BTree).BlockSize("10M").Comment("index comment").EngineAttribute("engine attr").SecondaryEngineAttribute("sec engine attr").WithParser("parser")
	var builder strings.Builder
	i.Build(&builder)
	assert.Equal(t, " KEY_BLOCK_SIZE = 10M USING BTREE WITH PARSER parser COMMENT 'index comment' INVISIBLE ENGINE_ATTRIBUTE = 'engine attr' SECONDARY_ENGINE_ATTRIBUTE = 'sec engine attr'", builder.String())
}

func TestIndexOptionVisible(t *testing.T) {
	i := NewIndexOption().Visible().Algorithm(ksql.Index_Alg_BTree).BlockSize("10M")
	var builder strings.Builder
	i.Build(&builder)
	assert.Equal(t, " KEY_BLOCK_SIZE = 10M USING BTREE VISIBLE", builder.String())
}

func TestIndexOptionComment(t *testing.T) {
	i := NewIndexOption().Comment("index comment").EngineAttribute("engine attr").SecondaryEngineAttribute("sec engine attr").WithParser("parser")
	var builder strings.Builder
	i.Build(&builder)
	assert.Equal(t, " WITH PARSER parser COMMENT 'index comment' ENGINE_ATTRIBUTE = 'engine attr' SECONDARY_ENGINE_ATTRIBUTE = 'sec engine attr'", builder.String())
}
