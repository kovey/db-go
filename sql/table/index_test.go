package table

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestIndexUnique(t *testing.T) {
	var builder strings.Builder
	i := NewIndex("idx_test")
	i.Unique().SubType(ksql.Index_Sub_Type_Key)
	i.Columns("user_id", "type")
	i.Build(&builder)
	assert.Equal(t, " UNIQUE KEY `idx_test` (`user_id`, `type`)", builder.String())
	assert.Equal(t, ksql.Index_Type_Unique, i.typ)
	builder.Reset()
	i = NewIndex("idx_test")
	i.Unique().SubType(ksql.Index_Sub_Type_Key)
	i.Columns("user_id", "type")
	i.SubType(ksql.Index_Sub_Type_Index).Build(&builder)
	assert.Equal(t, " UNIQUE INDEX `idx_test` (`user_id`, `type`)", builder.String())
	assert.Equal(t, ksql.Index_Type_Unique, i.typ)
}

func TestIndexNormal(t *testing.T) {
	var builder strings.Builder
	i := NewIndex("idx_test")
	i.Columns("user_id")
	i.Build(&builder)
	assert.Equal(t, " INDEX `idx_test` (`user_id`)", builder.String())
	assert.Equal(t, ksql.Index_Type_Normal, i.typ)
	builder.Reset()
	i = NewIndex("idx_test")
	i.Columns("user_id")
	i.SubType(ksql.Index_Sub_Type_Key).Build(&builder)
	assert.Equal(t, " KEY `idx_test` (`user_id`)", builder.String())
	assert.Equal(t, ksql.Index_Type_Normal, i.typ)
}

func TestIndexPrimary(t *testing.T) {
	var builder strings.Builder
	i := NewIndex("idx_test")
	i.Primary()
	i.Columns("user_id").Build(&builder)
	assert.Equal(t, " PRIMARY KEY (`user_id`)", builder.String())
	assert.Equal(t, ksql.Index_Type_Primary, i.typ)
}

func TestIndexSpatial(t *testing.T) {
	var builder strings.Builder
	i := NewIndex("idx_test")
	i.Type(ksql.Index_Type_Spatial).SubType(ksql.Index_Sub_Type_Key)
	i.Columns("user_id").Build(&builder)
	assert.Equal(t, " SPATIAL KEY `idx_test` (`user_id`)", builder.String())
	assert.Equal(t, ksql.Index_Type_Spatial, i.typ)
	builder.Reset()
	i = NewIndex("idx_test")
	i.Type(ksql.Index_Type_Spatial).SubType(ksql.Index_Sub_Type_Key)
	i.Columns("user_id")
	i.SubType(ksql.Index_Sub_Type_Index).Build(&builder)
	assert.Equal(t, " SPATIAL INDEX `idx_test` (`user_id`)", builder.String())
}

func TestIndexFullText(t *testing.T) {
	var builder strings.Builder
	i := NewIndex("idx_test")
	i.Type(ksql.Index_Type_FullText).SubType(ksql.Index_Sub_Type_Key)
	i.Columns("user_id").Build(&builder)
	assert.Equal(t, " FULLTEXT KEY `idx_test` (`user_id`)", builder.String())
	assert.Equal(t, ksql.Index_Type_FullText, i.typ)
	builder.Reset()
	i = NewIndex("idx_test")
	i.Type(ksql.Index_Type_FullText).SubType(ksql.Index_Sub_Type_Key)
	i.Columns("user_id")
	i.SubType(ksql.Index_Sub_Type_Index).Build(&builder)
	assert.Equal(t, " FULLTEXT INDEX `idx_test` (`user_id`)", builder.String())
}

func TestIndexForeign(t *testing.T) {
	var builder strings.Builder
	i := NewIndex("idx_test")
	i.Foreign().Constraint("test")
	i.Reference("user_table").Column("user_id", 10, ksql.Order_Asc).Match(ksql.Reference_Match_Full).On(ksql.Reference_On_Opt_DELETE, ksql.Reference_Option_Cascade).Express("round(amount * age, 2)", ksql.Order_Asc)
	i.Columns("user_id").Build(&builder)
	assert.Equal(t, " CONSTRAINT test FOREIGN KEY `idx_test` (`user_id`) REFERENCES `user_table` (`user_id`(10) ASC, (round(amount * age, 2)) ASC) MATCH FULL ON DELETE CASCADE", builder.String())
	assert.Equal(t, ksql.Index_Type_Foreign, i.typ)
}
