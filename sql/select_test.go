package sql

import (
	"testing"

	"github.com/kovey/db-go/v2/sql/meta"
)

func TestSelectPrepare(t *testing.T) {
	sSub := NewSelect("user", "um")
	sSub.Columns("id", "name", "sex", "create_time", "amount", "user_id")
	sSub.WhereByMap(meta.Where{"sex": 1})
	s := NewSelectSub(sSub, "u")
	s.ColMeta(meta.NewColumnAlias("id", "uid"), meta.NewColumnAlias("name", "username"), meta.NewColumnAlias("sex", "ss")).Limit(20).Offset(10).Order("create_time DESC", "id ASC").Group("id", "name")
	caseWhen1 := meta.NewCaseWhen("uid")
	caseWhen1.AddWhenThen("u.user_id > 0", "u.user_id")
	caseWhen1.AddWhenThen("u.user_id <= 0", "0")
	caseWhen1.Else("-1")
	s.CaseWhen(caseWhen1)
	field := meta.NewField("1", "", true)
	s.ColMeta(meta.NewColumnFunc(field, "count", meta.Func_COUNT, nil))
	field1 := meta.NewField("amount", "u", false)
	s.ColMeta(meta.NewColFuncWithNull(field1, "count", "0", meta.Func_SUM, nil))
	j := NewJoin("demo", "d", "d.user_id=u.id", "password")
	j.ColMeta(meta.NewColumnAlias("name", "dname"))
	s.LeftJoinWith(j)
	s.RightJoin("cron", "c", "c.user_id=u.id", "time", "run")
	sub := NewSelect("base_info", "bi")
	sub.Columns("age", "device", "sign").WhereByMap(meta.Where{"id": 100, "status": 1, "open": "on"})
	j2 := NewJoinSub(sub, "e", "bi.user_id=u.id", "age", "device", "sign")
	caseWhen := meta.NewCaseWhen("uid")
	caseWhen.AddWhenThen("e.user_id > 0", "e.user_id")
	caseWhen.AddWhenThen("e.user_id <= 0", "0")
	caseWhen.Else("-1")
	j2.CaseWhen(caseWhen)
	s.InnerJoinWith(j2)

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
