package sql

import "testing"

func TestDeletePrepare(t *testing.T) {
	del := NewDelete("product")
	where := NewWhere()
	where.Like("nickname", "%hello%")
	del.Where(where)
	del.WhereByMap(map[string]any{"id": 1, "name": "kovey"}).WhereByList([]string{"last_id > 0"})
	t.Logf("sql: %s", del)
}

func TestCkDeletePrepare(t *testing.T) {
	del := NewCkDelete("product")
	where := NewWhere()
	where.Like("nickname", "%hello%")
	del.Where(where)
	del.WhereByMap(map[string]any{"id": 1, "name": "kovey"}).WhereByList([]string{"last_id > 0"})
	t.Logf("sql: %s", del)
}
