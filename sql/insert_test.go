package sql

import "testing"

func TestInsert(t *testing.T) {
	in := NewInsert("product")
	in.Set("name", "kovey").Set("sex", 1).Set("time", "2021-01-01 11:11:11")
	sql := in.String()
	t.Logf("sql: %s", sql)
}
