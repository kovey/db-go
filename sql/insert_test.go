package sql

import "testing"

func TestInsert(t *testing.T) {
	in := NewInsert()
	in.Table("user")
	in.Add("name", "kovey").Add("time", "2024-03-05")
	t.Logf("prepare: %s", in.Prepare())
	t.Logf("binds: %s", in.Binds())
}
