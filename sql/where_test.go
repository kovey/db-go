package sql

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestWhere(t *testing.T) {
	h := NewWhere()
	sub := NewQuery()
	sub.Table("user").Columns("user_id").Where("status", "<>", 1)
	t.Logf("sub prepare: %s", sub.Prepare())
	t.Logf("sub binds: %v", sub.Binds())
	h.Between("a.id", 10, 20).Express(Raw("b.name = ?", "kovey")).Where("a.age", ">", 100).In("a.sex", []any{1, 2, 3}).InBy("b.status", sub).IsNotNull("a.mail").IsNull("b.avatar")
	h.NotBetween("a.num", 1000, 2000)
	h.NotIn("c.id", []any{123, 45}).NotInBy("c.test", sub).OrWhere(func(o ksql.WhereInterface) {
		o.Between("d.info", 100, 200).Where("d.other", "like", "%kk%")
	})
	h.AndWhere(func(o ksql.WhereInterface) {
		o.IsNotNull("d.test").Where("d.ss", ksql.Neq, 1)
	})

	var builder strings.Builder
	h.Build(&builder)
	assert.Equal(t, "WHERE `a`.`id` BETWEEN ? AND ? AND b.name = ? AND `a`.`age` > ? AND `a`.`sex` IN (?, ?, ?) AND `b`.`status` IN (SELECT `user_id` FROM `user` WHERE `status` <> ?) AND `a`.`mail` IS NOT NULL AND `b`.`avatar` IS NULL AND `a`.`num` NOT BETWEEN ? AND ? AND `c`.`id` NOT IN (?, ?) AND `c`.`test` NOT IN (SELECT `user_id` FROM `user` WHERE `status` <> ?) AND (`d`.`test` IS NOT NULL AND `d`.`ss` <> ?) OR (`d`.`info` BETWEEN ? AND ? AND `d`.`other` like ?)", builder.String())
	assert.Equal(t, []any{10, 20, "kovey", 100, 1, 2, 3, 1, 1000, 2000, 123, 45, 1, 1, 100, 200, "%kk%"}, h.Binds())
	assert.Equal(t, "WHERE `a`.`id` BETWEEN ? AND ? AND b.name = ? AND `a`.`age` > ? AND `a`.`sex` IN (?, ?, ?) AND `b`.`status` IN (SELECT `user_id` FROM `user` WHERE `status` <> ?) AND `a`.`mail` IS NOT NULL AND `b`.`avatar` IS NULL AND `a`.`num` NOT BETWEEN ? AND ? AND `c`.`id` NOT IN (?, ?) AND `c`.`test` NOT IN (SELECT `user_id` FROM `user` WHERE `status` <> ?) AND (`d`.`test` IS NOT NULL AND `d`.`ss` <> ?) OR (`d`.`info` BETWEEN ? AND ? AND `d`.`other` like ?)", builder.String())
}
