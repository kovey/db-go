package sql

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestHaving(t *testing.T) {
	h := NewHaving()
	sub := NewQuery()
	sub.Table("user").Columns("user_id").Where("status", "<>", 1)
	h.Between("a.id", 10, 20).Express(Raw("b.name = ?", "kovey")).Having("a.age", ">", 100).In("a.sex", []any{1, 2, 3}).InBy("b.status", sub).IsNotNull("a.mail").IsNull("b.avatar")
	h.NotBetween("b.num", 100, 200)
	h.NotIn("c.id", []any{123, 45}).NotInBy("c.test", sub).OrHaving(func(o ksql.HavingInterface) {
		o.Between("d.info", 100, 200).Having("d.other", "like", "%kk%")
	})

	var builder strings.Builder
	h.Build(&builder)
	assert.Equal(t, "HAVING `a`.`id` BETWEEN ? AND ? AND b.name = ? AND `a`.`age` > ? AND `a`.`sex` IN (?, ?, ?) AND `b`.`status` IN (SELECT `user_id` FROM `user` WHERE `status` <> ?) AND `a`.`mail` IS NOT NULL AND `b`.`avatar` IS NULL AND `b`.`num` NOT BETWEEN ? AND ? AND `c`.`id` NOT IN (?, ?) AND `c`.`test` NOT IN (SELECT `user_id` FROM `user` WHERE `status` <> ?) OR (`d`.`info` BETWEEN ? AND ? AND `d`.`other` like ?)", builder.String())
	assert.Equal(t, []any{10, 20, "kovey", 100, 1, 2, 3, 1, 100, 200, 123, 45, 1, 100, 200, "%kk%"}, h.Binds())
	assert.Equal(t, "HAVING `a`.`id` BETWEEN ? AND ? AND b.name = ? AND `a`.`age` > ? AND `a`.`sex` IN (?, ?, ?) AND `b`.`status` IN (SELECT `user_id` FROM `user` WHERE `status` <> ?) AND `a`.`mail` IS NOT NULL AND `b`.`avatar` IS NULL AND `b`.`num` NOT BETWEEN ? AND ? AND `c`.`id` NOT IN (?, ?) AND `c`.`test` NOT IN (SELECT `user_id` FROM `user` WHERE `status` <> ?) OR (`d`.`info` BETWEEN ? AND ? AND `d`.`other` like ?)", builder.String())
}
