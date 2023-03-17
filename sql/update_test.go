package sql

import "testing"

func TestUpdatePrepare(t *testing.T) {
	up := NewUpdate("product")
	up.Set("name", "kovey").Set("sex", 1).WhereByList([]string{"id = 1"})
	up.WhereByMap(map[string]any{"name": "golang", "value": 1})
	where := NewWhere()
	where.Eq("other", 1).Gt("last_id", 20).Between("time", 1000, 2000)
	up.Where(where)
	t.Logf("sql: %s", up)
}

func TestCkUpdatePrepare(t *testing.T) {
	up := NewCkUpdate("product")
	up.Set("name", "kovey").Set("sex", 1).WhereByList([]string{"id = 1"})
	up.WhereByMap(map[string]any{"name": "golang", "value": 1})
	where := NewWhere()
	where.Eq("other", 1).Gt("last_id", 20).Between("time", 1000, 2000)
	up.Where(where)
	t.Logf("sql: %s", up)
}
