package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	q := NewQuery()
	q.Table("user").As("u").Between("u.age", 10, 20).Column("nickname", "name").Columns("avatar", "sex", "id_card").ColumnsExpress(Raw("username as account")).ForUpdate()
	q.NotBetween("u.phone", 1000, 2000)
	q.Func("sum", "balance", "balance").Group("user_id").Group("sex").Having("balance", ">", 100).HavingExpress(Raw("name is not null")).Join("ext").As("e").On("e.id", "=", "u.user_id")
	q.JoinExpress(Raw("info as i on i.id = u.user_id and i.status = ?", 1))
	q.LeftJoin("account").As("a").On("a.id", "=", "u.user_id").OnOr(func(joi ksql.JoinOnInterface) {
		joi.On("a.balance", ">", "0").On("a.freezen", "=", "0")
	})
	q.RightJoin("login_info").As("li").On("li.user_id", "=", "u.user_id").OnOr(func(joi ksql.JoinOnInterface) {
		joi.On("li.count", ">", "0").On("li.date", "=", "0")
	})
	q.Where("u.id", ">", 1000).WhereIn("u.status", []any{1, 2, 3}).WhereIsNotNull("u.avatar").WhereIsNull("u.id_card").WhereExpress(Raw("i.id > ? and i.status = ?", 100, 1)).WhereNotIn("i.status", []any{4, 5})
	sub := NewQuery()
	sub.Table("game").Columns("id").Where("status", "=", 5)
	room := NewQuery()
	room.Table("room").Columns("id").Where("status", "=", 5)
	q.WhereInBy("u.game_id", sub).WhereNotInBy("u.room_id", room).OrWhere(func(wi ksql.WhereInterface) {
		wi.Where("u.game_status", "=", 10)
		wi.IsNull("u.game_other")
	})
	q.OrHaving(func(h ksql.HavingInterface) {
		h.Between("li.count", 100, 200)
		h.Having("li.coin", ">", 1000)
		h.NotBetween("li.page", 3000, 5000)
	})
	q.Limit(10).Offset(0).Order("u.user_id").OrderDesc("li.balance")

	assert.Equal(t, "SELECT `nickname` AS `name`,`avatar`,`sex`,`id_card`,username as account,sum(`balance`) AS `balance` FROM `user` AS `u`  JOIN  `ext` AS `e` ON `e`.`id`=`u`.`user_id` JOIN info as i on i.id = u.user_id and i.status = ? LEFT JOIN `account` AS `a` ON `a`.`id`=`u`.`user_id` OR (`a`.`balance`>0 AND `a`.`freezen`=0) RIGHT JOIN `login_info` AS `li` ON `li`.`user_id`=`u`.`user_id` OR (`li`.`count`>0 AND `li`.`date`=0) WHERE `u`.`age` BETWEEN ?  AND ?  AND `u`.`phone` NOT BETWEEN ?  AND ?  AND `u`.`id` > ? AND `u`.`status` IN (?,?,?) AND `u`.`avatar` IS NOT NULL AND `u`.`id_card` IS NULL AND i.id > ? and i.status = ? AND `i`.`status` NOT IN (?,?) AND `u`.`game_id` IN (SELECT `id` FROM `game` WHERE `status` = ?) AND `u`.`room_id` NOT IN (SELECT `id` FROM `room` WHERE `status` = ?) OR (`u`.`game_status` = ? AND `u`.`game_other` IS NULL) GROUP BY `user_id`,`sex` HAVING `balance` > ? AND name is not null OR (`li`.`count` BETWEEN ?  AND ?  AND `li`.`coin` > ? AND `li`.`page` NOT BETWEEN ?  AND ? ) ORDER BY `u`.`user_id` ASC,`li`.`balance` DESC LIMIT 0,10 FOR UPDATE", q.Prepare())
	assert.Equal(t, []any{1, 10, 20, 1000, 2000, 1000, 1, 2, 3, 100, 1, 4, 5, 5, 5, 10, 100, 100, 200, 1000, 3000, 5000}, q.Binds())
	assert.Equal(t, "SELECT `nickname` AS `name`,`avatar`,`sex`,`id_card`,username as account,sum(`balance`) AS `balance` FROM `user` AS `u`  JOIN  `ext` AS `e` ON `e`.`id`=`u`.`user_id` JOIN info as i on i.id = u.user_id and i.status = ? LEFT JOIN `account` AS `a` ON `a`.`id`=`u`.`user_id` OR (`a`.`balance`>0 AND `a`.`freezen`=0) RIGHT JOIN `login_info` AS `li` ON `li`.`user_id`=`u`.`user_id` OR (`li`.`count`>0 AND `li`.`date`=0) WHERE `u`.`age` BETWEEN ?  AND ?  AND `u`.`phone` NOT BETWEEN ?  AND ?  AND `u`.`id` > ? AND `u`.`status` IN (?,?,?) AND `u`.`avatar` IS NOT NULL AND `u`.`id_card` IS NULL AND i.id > ? and i.status = ? AND `i`.`status` NOT IN (?,?) AND `u`.`game_id` IN (SELECT `id` FROM `game` WHERE `status` = ?) AND `u`.`room_id` NOT IN (SELECT `id` FROM `room` WHERE `status` = ?) OR (`u`.`game_status` = ? AND `u`.`game_other` IS NULL) GROUP BY `user_id`,`sex` HAVING `balance` > ? AND name is not null OR (`li`.`count` BETWEEN ?  AND ?  AND `li`.`coin` > ? AND `li`.`page` NOT BETWEEN ?  AND ? ) ORDER BY `u`.`user_id` ASC,`li`.`balance` DESC LIMIT 0,10 FOR UPDATE", q.Prepare())
}
