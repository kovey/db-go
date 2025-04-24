package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteMulti(t *testing.T) {
	m := NewDeleteMulti().Table("user").TableAs("user_ext", "ue").Ignore().Quick().LowPriority()
	m.Join("game").As("g").On("g.user_id", "=", "user.user_id")
	m.JoinExpress(Raw("LEFT JOIN hat as h ON h.user_id = ue.user_id"))
	m.LeftJoin("sex").As("s").On("s.user_id", "=", "ue.user_id")
	m.RightJoin("info").As("i").On("i.user_id", "=", "ue.user_id")
	w := NewWhere().Where("user.user_id", ">", 1)
	m.Where(w)

	assert.Equal(t, "DELETE LOW_PRIORITY QUICK IGNORE `user`, `user_ext` AS `ue` FROM INNER JOIN `game` AS `g` ON (`g`.`user_id` = `user`.`user_id`) LEFT JOIN hat as h ON h.user_id = ue.user_id LEFT JOIN `sex` AS `s` ON (`s`.`user_id` = `ue`.`user_id`) RIGHT JOIN `info` AS `i` ON (`i`.`user_id` = `ue`.`user_id`) WHERE `user`.`user_id` > ?", m.Prepare())
	assert.Equal(t, []any{1}, m.Binds())
}
