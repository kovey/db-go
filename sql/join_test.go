package sql

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestLeftJoin(t *testing.T) {
	j := NewJoin().Left()
	j.Table("user").As("u").On("u.id", "=", "c.id").OnOr(func(join ksql.JoinOnInterface) {
		join.On("c.name", "=", "u.name").OnVal("c.age", ">", 1)
	})

	var builder strings.Builder
	j.Build(&builder)
	assert.Equal(t, "LEFT JOIN `user` AS `u` ON (`u`.`id` = `c`.`id`) OR (`c`.`name` = `u`.`name` AND `c`.`age` > ?)", builder.String())
	assert.Equal(t, []any{1}, j.Binds())
	assert.Equal(t, "LEFT JOIN `user` AS `u` ON (`u`.`id` = `c`.`id`) OR (`c`.`name` = `u`.`name` AND `c`.`age` > ?)", builder.String())
}

func TestJoin(t *testing.T) {
	j := NewJoin().Inner()
	j.Table("user").As("u").On("u.id", "=", "c.id").OnOr(func(join ksql.JoinOnInterface) {
		join.On("c.name", "=", "u.name").OnVal("c.age", ">", 1)
	})

	var builder strings.Builder
	j.Build(&builder)
	assert.Equal(t, "INNER JOIN `user` AS `u` ON (`u`.`id` = `c`.`id`) OR (`c`.`name` = `u`.`name` AND `c`.`age` > ?)", builder.String())
	assert.Equal(t, []any{1}, j.Binds())
	assert.Equal(t, "INNER JOIN `user` AS `u` ON (`u`.`id` = `c`.`id`) OR (`c`.`name` = `u`.`name` AND `c`.`age` > ?)", builder.String())
}

func TestRightJoin(t *testing.T) {
	j := NewJoin().Right()
	j.Table("user").As("u").On("u.id", "=", "c.id").OnOr(func(join ksql.JoinOnInterface) {
		join.On("c.name", "=", "u.name").OnVal("c.age", ">", 1)
	})

	var builder strings.Builder
	j.Build(&builder)
	assert.Equal(t, "RIGHT JOIN `user` AS `u` ON (`u`.`id` = `c`.`id`) OR (`c`.`name` = `u`.`name` AND `c`.`age` > ?)", builder.String())
	assert.Equal(t, []any{1}, j.Binds())
	assert.Equal(t, "RIGHT JOIN `user` AS `u` ON (`u`.`id` = `c`.`id`) OR (`c`.`name` = `u`.`name` AND `c`.`age` > ?)", builder.String())
}
