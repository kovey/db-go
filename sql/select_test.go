package sql

import "testing"

func TestSelectPrepare(t *testing.T) {
	s := NewSelect("user", "u")
	s.Columns("id", "name", "sex").Limit(20).Offset(10).Order("create_time DESC", "id ASC").Group("id", "name")
	s.LeftJoin("demo", "d", "d.user_id=u.id", "name", "password")
	s.RightJoin("cron", "c", "c.user_id=u.id", "time", "run")
	s.InnerJoin("ext", "e", "e.user_id=u.id", "nickname", "addr")

	where := NewWhere()
	orWhere := NewWhere()
	having := NewHaving()

	where.Eq("u.id", 1).Ge("u.sex", 10)
	orWhere.Neq("u.name", "kovey").Like("c.name", "cron")
	having.In("u.sex", []any{1, 2, 3})

	s.Where(where).OrWhere(orWhere).Having(having)

	t.Logf("sql: %s", s)
}
