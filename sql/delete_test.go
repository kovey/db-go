package sql

import "testing"

func TestDetele(t *testing.T) {
	d := NewDelete()
	w := NewWhere()
	w.Where("id", "=", 1)
	d.Table("user").Where(w)

	t.Logf("prepare: %s", d.Prepare())
	t.Logf("binds: %v", d.Binds())
}
