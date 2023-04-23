package sql

import (
	"testing"

	"github.com/kovey/db-go/v2/sql/meta"
)

func TestDeletePrepare(t *testing.T) {
	del := NewDelete("product")
	where := NewWhere()
	where.Like("nickname", "%hello%")
	del.Where(where)
	del.WhereByMap(meta.Where{"id": 1, "name": "kovey"}).WhereByList([]string{"last_id > 0"})
	t.Logf("sql: %s", del)
}

func TestCkDeletePrepare(t *testing.T) {
	del := NewCkDelete("product")
	where := NewWhere()
	where.Like("nickname", "%hello%")
	del.Where(where)
	del.WhereByMap(meta.Where{"id": 1, "name": "kovey"}).WhereByList([]string{"last_id > 0"})
	t.Logf("sql: %s", del)
}
