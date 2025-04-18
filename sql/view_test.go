package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestView(t *testing.T) {
	q := NewQuery().Table("user").Columns("id", "name", "age").Where("id", ">=", 1)
	v := NewView().View("tmp_view").Algorithm(ksql.View_Alg_Merge).As(q)
	v.Definer("root@local").Columns("id", "name", "age").Replace().SqlSecurity(ksql.Sql_Security_Definer).WithCascaded().CheckOption()
	assert.Equal(t, "CREATE OR REPLACE ALGORITHM = MERGE DEFINER = root@local SQL SECURITY DEFINER VIEW `tmp_view` (`id`, `name`, `age`) AS SELECT `id`, `name`, `age` FROM `user` WHERE `id` >= ? WITH CASCADED CHECK OPTION", v.Prepare())
	assert.Equal(t, []any{1}, v.Binds())
}

func TestViewAlter(t *testing.T) {
	q := NewQuery().Table("user").Columns("id", "name", "age").Where("id", ">=", 1)
	v := NewView().View("tmp_view").Algorithm(ksql.View_Alg_Merge).As(q).Alter()
	v.Definer("root@local").Columns("id", "name", "age").Replace().SqlSecurity(ksql.Sql_Security_Definer).WithCascaded().CheckOption()
	assert.Equal(t, "ALTER ALGORITHM = MERGE DEFINER = root@local SQL SECURITY DEFINER VIEW `tmp_view` (`id`, `name`, `age`) AS SELECT `id`, `name`, `age` FROM `user` WHERE `id` >= ? WITH CASCADED CHECK OPTION", v.Prepare())
	assert.Equal(t, []any{1}, v.Binds())
}
