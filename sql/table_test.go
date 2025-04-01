package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
)

func TestTable(t *testing.T) {
	ta := NewTable()
	ta.Table("user").AddColumn("id", "bigint", 20, 0).AutoIncrement().Unsigned().Comment("主键")
	ta.AddColumn("username", "VARCHAR", 31, 0).Nullable().Default("NULL").Comment("用户名")
	ta.AddColumn("password", "VARCHAR", 64, 0).Default("").Comment("密码")
	ta.AddColumn("age", "int", 11, 0).Default("0").Comment("密码")
	ta.AddColumn("create_time", "TIMESTAMP", 0, 0).Default(ksql.CURRENT_TIMESTAMP).Comment("创建时间")
	ta.AddColumn("update_time", "TIMESTAMP", 0, 0).Default(ksql.CURRENT_TIMESTAMP_ON_UPDATE_CURRENT_TIMESTAMP).Comment("更新时间")
	ta.AddPrimary("id").AddIndex("idx_username", ksql.Index_Type_Unique, "username").AddIndex("idx_name_age", ksql.Index_Type_Normal, "username", "age")
	ta.Engine("InnoDB").Charset("utf8").Collate("test").Comment("用户表")

	t.Logf("prepare: %s", ta.Prepare())
	t.Logf("binds: %v", ta.Binds())
}

func TestTableLike(t *testing.T) {
	ta := NewTable().Table("user").Like("users")
	t.Logf("prepare: %s", ta.Prepare())
	t.Logf("binds: %v", ta.Binds())
}

func TestTableFrom(t *testing.T) {
	query := NewQuery()
	query.Table("user_back").Columns("u.name", "u.kovey", "e.date").As("u")
	query.Join("email").As("e").On("e.id", "=", "u.id")
	query.Where("u.name", "LIKE", "%test%")
	ta := NewTable().Table("user").From(query)
	t.Logf("prepare: %s", ta.Prepare())
	t.Logf("binds: %v", ta.Binds())
}
