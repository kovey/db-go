package table

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestPartitionOperateAdd(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Add().Option(ksql.Partition_Definition_Opt_Key_Comment, "comment").Option(ksql.Partition_Definition_Opt_Key_Data_Directory, "data.dt").Option(ksql.Partition_Definition_Opt_Key_Engine, "memery").Values().LessThan().MaxValue()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "ADD (`part1` VALUES  LESS THAN MAXVALUE COMMENT = 'comment' DATA DIRECTORY = 'data.dt' STORAGE ENGINE = memery)", builder.String())
}

func TestPartitionOperateAnalyze(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Analyze()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "ANALYZE PARTITION `part1`, `part2`", builder.String())
}

func TestPartitionOperateCheck(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Check()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "CHECK PARTITION `part1`, `part2`", builder.String())
}

func TestPartitionOperateCoalesce(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Coalesce(20)
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "COALESCE PARTITION `part1`, `part2` 20", builder.String())
}

func TestPartitionOperateDiscard(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Discard()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "DISCARD PARTITION `part1`, `part2` TABLESPACE", builder.String())
}

func TestPartitionOperateDrop(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Drop()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "DROP `part1`, `part2`", builder.String())
}

func TestPartitionOperateImport(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Import()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "IMPORT PARTITION `part1`, `part2` TABLESPACE", builder.String())
}

func TestPartitionOperateExchange(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Exchange("table_name", "VAL")
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "EXCHANGE PARTITION `part1`, `part2` WITH TABLE `table_name` VAL VALIDATION", builder.String())
}

func TestPartitionOperateOptimize(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Optimize()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "OPTIMIZE PARTITION `part1`, `part2`", builder.String())
}

func TestPartitionOperateRebuild(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Rebuild()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "REBUILD PARTITION `part1`, `part2`", builder.String())
}

func TestPartitionOperateRemove(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Remove()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "REMOVE PARTITIONING `part1`, `part2`", builder.String())
}

func TestPartitionOperateReorganize(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Reorganize()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "REORGANIZE PARTITION INTO (`part1`)", builder.String())
}

func TestPartitionOperateRepair(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Repair()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "REPAIR PARTITION `part1`, `part2`", builder.String())
}

func TestPartitionOperateTruncate(t *testing.T) {
	o := NewPartitionOperate("part1", "part2")
	o.Truncate()
	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "TRUNCATE PARTITION `part1`, `part2`", builder.String())
}
