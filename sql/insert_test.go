package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	in := NewInsert()
	in.Table("user")
	in.Add("name", "kovey").Add("time", "2024-03-05")
	assert.Equal(t, "INSERT INTO `user` (`name`,`time`) VALUES (?,?)", in.Prepare())
	assert.Equal(t, []any{"kovey", "2024-03-05"}, in.Binds())
	assert.Equal(t, "INSERT INTO `user` (`name`,`time`) VALUES (?,?)", in.Prepare())
}

func TestInsertFrom(t *testing.T) {
	query := NewQuery()
	query.Table("user_back").Columns("u.name", "u.kovey", "e.date").As("u")
	query.Join("email").As("e").On("e.id", "=", "u.id")
	query.Where("u.name", "LIKE", "%test%")
	in := NewInsert()
	in.Table("user").Columns("name", "kovey", "date").From(query)
	assert.Equal(t, "INSERT INTO `user` (`name`,`kovey`,`date`) SELECT `u`.`name`,`u`.`kovey`,`e`.`date` FROM `user_back` AS `u`  JOIN  `email` AS `e` ON `e`.`id`=`u`.`id` WHERE `u`.`name` LIKE ?", in.Prepare())
	assert.Equal(t, []any{"%test%"}, in.Binds())
	assert.Equal(t, "INSERT INTO `user` (`name`,`kovey`,`date`) SELECT `u`.`name`,`u`.`kovey`,`e`.`date` FROM `user_back` AS `u`  JOIN  `email` AS `e` ON `e`.`id`=`u`.`id` WHERE `u`.`name` LIKE ?", in.Prepare())
}
