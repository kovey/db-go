package sql

import (
	"testing"

	"github.com/kovey/db-go/v3"
)

func TestAlter(t *testing.T) {
	a := NewAlter()
	a.Table("user").AddColumn("user_name", "VARCHAR", 62, 0).Nullable().Default("NULL", true).Comment("用户名")
	a.DropColumn("age").AddColumn("balance", "DECIMAL", 10, 2).Default("0", false).Comment("余额")
	a.AddIndex("user_name", ksql.Index_Type_Unique, "user_name")
	a.DropIndex("idx_name").Comment("用户表").AddPrimary("id")

	t.Logf("prepare: %s", a.Prepare())
	t.Logf("binds: %v", a.Binds())
}
