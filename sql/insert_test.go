package sql

import "testing"

func TestInsert(t *testing.T) {
	in := NewInsert()
	in.Table("user")
	in.Add("name", "kovey").Add("time", "2024-03-05")
	t.Logf("prepare: %s", in.Prepare())
	t.Logf("binds: %s", in.Binds())
}

func TestInsertFrom(t *testing.T) {
	query := NewQuery()
	query.Table("user_back").Columns("u.name", "u.kovey", "e.date").As("u")
	query.Join("email").As("e").On("e.id", "=", "u.id")
	query.Where("u.name", "LIKE", "%test%")
	in := NewInsert()
	in.Table("user").Columns("name", "kovey", "date").From(query)
	t.Logf("prepare: %s", in.Prepare())
	t.Logf("binds: %s", in.Binds())
}
