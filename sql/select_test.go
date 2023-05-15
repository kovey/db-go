package sql

import (
	"testing"

	"github.com/kovey/db-go/v2/sql/meta"
)

func TestSelectPrepare(t *testing.T) {
	s := NewSelect("user", "u")
	s.Columns(meta.NewColumn("id", "uid"), meta.NewColumn("name", "username"), meta.NewColumn("sex", "ss")).Limit(20).Offset(10).Order("create_time DESC", "id ASC").Group("id", "name")
	caseWhen1 := meta.NewCaseWhen("uid")
	caseWhen1.AddWhenThen("u.user_id > 0", "u.user_id")
	caseWhen1.AddWhenThen("u.user_id <= 0", "0")
	caseWhen1.Else("-1")
	s.CaseWhen(caseWhen1)
	field := meta.NewField("1", "", true)
	s.Columns(meta.NewFuncComlumn(field, "count", meta.Func_COUNT, nil))
	field1 := meta.NewField("amount", "u", false)
	s.Columns(meta.NewFuncComlumn(field1, "count", meta.Func_SUM, nil))
	j := NewJoin("demo", "d", "d.user_id=u.id")
	j.Columns(meta.NewColumn("name", "dname"), meta.NewColumn("password", "password"))
	s.LeftJoin(j)
	j1 := NewJoin("cron", "c", "c.user_id=u.id")
	j1.Columns(meta.NewColumn("time", "tt"), meta.NewColumn("run", "runn"))
	s.RightJoin(j1)
	j2 := NewJoin("ext", "e", "e.user_id=u.id")
	j2.Columns(meta.NewColumn("nickname", "nnn"), meta.NewColumn("addr", "addr"))
	caseWhen := meta.NewCaseWhen("uid")
	caseWhen.AddWhenThen("e.user_id > 0", "e.user_id")
	caseWhen.AddWhenThen("e.user_id <= 0", "0")
	caseWhen.Else("-1")
	j2.CaseWhen(caseWhen)
	s.InnerJoin(j2)

	where := NewWhere()
	orWhere := NewWhere()
	having := NewHaving()

	where.Eq("u.id", 1)
	where.Ge("u.sex", 10)
	orWhere.Neq("u.name", "kovey")
	orWhere.Like("c.name", "cron")
	having.In("u.sex", []any{1, 2, 3})

	s.Where(where).OrWhere(orWhere).Having(having)

	t.Logf("sql: %s", s)
}
