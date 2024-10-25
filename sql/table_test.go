package sql

import (
	"testing"

	"github.com/kovey/db-go/v3"
)

func TestTable(t *testing.T) {
	ta := NewTable()
	ta.Table("user").AddColumn("id", "bigint", 20, 0).AutoIncrement().Unsigned().Comment("主键")
	ta.AddColumn("username", "VARCHAR", 31, 0).Nullable().Default("NULL", true).Comment("用户名")
	ta.AddColumn("password", "VARCHAR", 64, 0).Default("", false).Comment("密码")
	ta.AddColumn("age", "int", 11, 0).Default("0", false).Comment("密码")
	ta.AddColumn("create_time", "TIMESTAMP", 0, 0).Default(ksql.CURRENT_TIMESTAMP, true).Comment("创建时间")
	ta.AddColumn("update_time", "TIMESTAMP", 0, 0).Default(ksql.CURRENT_TIMESTAMP_ON_UPDATE_CURRENT_TIMESTAMP, true).Comment("更新时间")
	ta.AddPrimary("id").AddIndex("idx_username", ksql.Index_Type_Unique, "username").AddIndex("idx_name_age", ksql.Index_Type_Normal, "username", "age")
	ta.Engine("InnoDB").Charset("utf8").Collate("test").Comment("用户表")

	t.Logf("prepare: %s", ta.Prepare())
	t.Logf("binds: %v", ta.Binds())
}
