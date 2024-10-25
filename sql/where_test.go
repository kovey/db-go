package sql

import (
	"testing"

	"github.com/kovey/db-go/v3"
)

func TestWhere(t *testing.T) {
	h := NewWhere()
	sub := NewQuery()
	sub.Table("user").Columns("user_id").Where("status", "<>", 1)
	t.Logf("sub prepare: %s", sub.Prepare())
	t.Logf("sub binds: %v", sub.Binds())
	h.Between("a.id", 10, 20).Express(Raw("b.name = ?", "kovey")).Where("a.age", ">", 100).In("a.sex", []any{1, 2, 3}).InBy("b.status", sub).IsNotNull("a.mail").IsNull("b.avatar")
	h.NotIn("c.id", []any{123, 45}).NotInBy("c.test", sub).OrWhere(func(o ksql.WhereInterface) {
		o.Between("d.info", 100, 200).Where("d.other", "like", "%kk%")
	})

	t.Logf("prepare: %s", h.Prepare())
	t.Logf("binds: %v", h.Binds())
}
