package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestTrigger(t *testing.T) {
	tg := NewTrigger().Trigger("up_trigger").After().BodyRaw(Raw("SELECT * FROM user WHERE user_id > ?", 1)).Update()
	in := NewInsert().Add("id", 1).Add("name", "kovey").Table("user_ext")
	tg.Body(in).Definer("root@local").IfNotExists().On("user").Order(ksql.Trigger_Order_Type_Follows, "other_trigger")

	assert.Equal(t, "CREATE DEFINER = root@local TRIGGER IF NOT EXISTS `up_trigger` AFTER UPDATE ON `user` FOR EACH ROW FOLLOWS `other_trigger` BEGIN SELECT * FROM user WHERE user_id > ?; INSERT INTO `user_ext` (`id`, `name`) VALUES (?, ?); END", tg.Prepare())
	assert.Equal(t, []any{1, 1, "kovey"}, tg.Binds())
}
