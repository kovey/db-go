package sql

import (
	"testing"

	"github.com/kovey/db-go/v3"
)

func TestQuery(t *testing.T) {
	q := NewQuery()
	q.Table("user").As("u").Between("u.age", 10, 20).Column("nickname", "name").Columns("avatar", "sex", "id_card").ColumnsExpress(Raw("username as account")).ForUpdate()
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
	})
	q.Limit(10).Offset(0).Order("u.user_id").OrderDesc("li.balance")

	t.Logf("prepare: %s", q.Prepare())
	t.Logf("binds: %v", q.Binds())
}
