package sql

import "testing"

func TestUpdate(t *testing.T) {
	u := NewUpdate()
	w := NewWhere()
	w.Where("id", "=", 1)
	u.Table("user").Set("name", "kovey").Set("age", 10).Where(w)
	t.Logf("prepare: %s", u.Prepare())
	t.Logf("binds: %v", u.Binds())
}
