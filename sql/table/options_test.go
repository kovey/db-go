package table

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestOptionsFull(t *testing.T) {
	o := NewOptions().AppendWithPrefix(ksql.Table_Opt_Prefix_Data, ksql.Table_Opt_Key_Directory, "./data").Append(ksql.Table_Opt_Key_Auto_Increment, "1000")
	o.AppendWith(ksql.Table_Opt_Key_Pack_Keys, ksql.Table_Opt_Val_0).RowFormat(ksql.Row_Format_Compact).InsertMethod(ksql.Insert_Method_First).StartTransation()
	o.Union("table1", "table2").SpaceOption("space_name", ksql.Column_Storage_Disk).SpaceOption("space_name_1", ksql.Column_Storage_Memory)

	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "DATA DIRECTORY = './data', AUTO_INCREMENT = 1000, PACK_KEYS = 0, ROW_FORMAT = COMPACT, INSERT_METHOD = FIRST, START TRANSACTION, UNION = (`table1`, `table2`), TABLESPACE `space_name` STORAGE DISK, TABLESPACE `space_name_1` STORAGE MEMORY", builder.String())
}

func TestOptionsNoUnionAndSpace(t *testing.T) {
	o := NewOptions().AppendWithPrefix(ksql.Table_Opt_Prefix_Data, ksql.Table_Opt_Key_Directory, "./data").Append(ksql.Table_Opt_Key_Auto_Increment, "1000")
	o.AppendWith(ksql.Table_Opt_Key_Pack_Keys, ksql.Table_Opt_Val_0).RowFormat(ksql.Row_Format_Compact).InsertMethod(ksql.Insert_Method_First).StartTransation()

	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "DATA DIRECTORY = './data', AUTO_INCREMENT = 1000, PACK_KEYS = 0, ROW_FORMAT = COMPACT, INSERT_METHOD = FIRST, START TRANSACTION", builder.String())
}

func TestOptionsOnlyUnionAndSpace(t *testing.T) {
	o := NewOptions()
	o.Union("table1", "table2").SpaceOption("space_name", ksql.Column_Storage_Disk).SpaceOption("space_name_1", ksql.Column_Storage_Memory)

	var builder strings.Builder
	o.Build(&builder)
	assert.Equal(t, "UNION = (`table1`, `table2`), TABLESPACE `space_name` STORAGE DISK, TABLESPACE `space_name_1` STORAGE MEMORY", builder.String())
}
