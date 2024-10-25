package sql

import (
	"testing"

	"github.com/kovey/db-go/v3"
)

func TestJoin(t *testing.T) {
	j := NewJoin("LEFT JOIN")
	j.Table("user").As("u").On("u.id", "=", "c.id").OnOr(func(join ksql.JoinOnInterface) {
		join.On("c.name", "=", "u.name").On("c.age", ">", "1")
	})

	t.Logf("prepare: %s", j.Prepare())
	t.Logf("binds: %v", j.Binds())
}
