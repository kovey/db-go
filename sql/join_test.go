package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestLeftJoin(t *testing.T) {
	j := NewJoin("LEFT JOIN")
	j.Table("user").As("u").On("u.id", "=", "c.id").OnOr(func(join ksql.JoinOnInterface) {
		join.On("c.name", "=", "u.name").On("c.age", ">", "1")
	})

	assert.Equal(t, "LEFT JOIN `user` AS `u` ON `u`.`id`=`c`.`id` OR (`c`.`name`=`u`.`name` AND `c`.`age`>1)", j.Prepare())
	assert.Nil(t, j.Binds())
	assert.Equal(t, "LEFT JOIN `user` AS `u` ON `u`.`id`=`c`.`id` OR (`c`.`name`=`u`.`name` AND `c`.`age`>1)", j.Prepare())
}

func TestJoin(t *testing.T) {
	j := NewJoin("JOIN")
	j.Table("user").As("u").On("u.id", "=", "c.id").OnOr(func(join ksql.JoinOnInterface) {
		join.On("c.name", "=", "u.name").On("c.age", ">", "1")
	})

	assert.Equal(t, "JOIN `user` AS `u` ON `u`.`id`=`c`.`id` OR (`c`.`name`=`u`.`name` AND `c`.`age`>1)", j.Prepare())
	assert.Nil(t, j.Binds())
	assert.Equal(t, "JOIN `user` AS `u` ON `u`.`id`=`c`.`id` OR (`c`.`name`=`u`.`name` AND `c`.`age`>1)", j.Prepare())
}

func TestRightJoin(t *testing.T) {
	j := NewJoin("RIGHT JOIN")
	j.Table("user").As("u").On("u.id", "=", "c.id").OnOr(func(join ksql.JoinOnInterface) {
		join.On("c.name", "=", "u.name").On("c.age", ">", "1")
	})

	assert.Equal(t, "RIGHT JOIN `user` AS `u` ON `u`.`id`=`c`.`id` OR (`c`.`name`=`u`.`name` AND `c`.`age`>1)", j.Prepare())
	assert.Nil(t, j.Binds())
	assert.Equal(t, "RIGHT JOIN `user` AS `u` ON `u`.`id`=`c`.`id` OR (`c`.`name`=`u`.`name` AND `c`.`age`>1)", j.Prepare())
}
