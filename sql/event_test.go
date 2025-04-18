package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestEventAt(t *testing.T) {
	e := NewEvent()
	e.Definer("kovey").IfNotExists().Event("k_event").OnCompletion().Status(ksql.Event_Status_Disable).Comment("test event")
	s := NewInsert().Table("user").Add("name", "kovey").Add("sex", 1)
	e.At("2025-04-08 11:11:11").AtInterval("1", ksql.Interval_Day).Do(s).AtInterval("2", ksql.Interval_Hour).Every("3", ksql.Interval_Minute)
	assert.Equal(t, "CREATE DEFINER = kovey EVENT IF NOT EXISTS `k_event` ON SCHEDULE AT '2025-04-08 11:11:11' + INTERVAL 1 DAY + INTERVAL 2 HOUR ON COMPLETION PRESERVE DISABLE COMMENT 'test event' DO BEGIN INSERT INTO `user` (`name`, `sex`) VALUES (?, ?); END", e.Prepare())
	assert.Equal(t, []any{"kovey", 1}, e.Binds())
	assert.Equal(t, "CREATE DEFINER = kovey EVENT IF NOT EXISTS `k_event` ON SCHEDULE AT '2025-04-08 11:11:11' + INTERVAL 1 DAY + INTERVAL 2 HOUR ON COMPLETION PRESERVE DISABLE COMMENT 'test event' DO BEGIN INSERT INTO `user` (`name`, `sex`) VALUES (?, ?); END", e.Prepare())
}

func TestEventEvery(t *testing.T) {
	e := NewEvent()
	e.Event("k_event")
	s := NewInsert().Table("user").Add("name", "kovey").Add("sex", 1)
	e.Every("1", ksql.Interval_Hour).Do(s).Starts("2025-04-08 11:11:11").StartsInterval("3:1", ksql.Interval_Day_Minute).StartsInterval("2:1", ksql.Interval_Minute_Second)
	e.Ends("2026-04-09 11:11:11").EndsInterval("1", ksql.Interval_Minute).At("2025-04-08 11:11:11").DoRaw(nil)
	assert.Equal(t, "CREATE EVENT `k_event` ON SCHEDULE EVERY + INTERVAL 1 HOUR STARTS '2025-04-08 11:11:11' + INTERVAL '3:1' DAY_MINUTE + INTERVAL '2:1' MINUTE_SECOND ENDS '2026-04-09 11:11:11' + INTERVAL 1 MINUTE DO BEGIN INSERT INTO `user` (`name`, `sex`) VALUES (?, ?); END", e.Prepare())
	assert.Equal(t, []any{"kovey", 1}, e.Binds())
	assert.Equal(t, "CREATE EVENT `k_event` ON SCHEDULE EVERY + INTERVAL 1 HOUR STARTS '2025-04-08 11:11:11' + INTERVAL '3:1' DAY_MINUTE + INTERVAL '2:1' MINUTE_SECOND ENDS '2026-04-09 11:11:11' + INTERVAL 1 MINUTE DO BEGIN INSERT INTO `user` (`name`, `sex`) VALUES (?, ?); END", e.Prepare())
}

func TestEventDoRaw(t *testing.T) {
	e := NewEvent()
	e.Event("k_event").OnCompletionNot().Rename("test")
	s := Raw("INSERT INTO `user` (`name`, `sex`) VALUES (?, ?)", "kovey", 1)
	e.Every("1:1", ksql.Interval_Hour_Minute).DoRaw(s).Starts("2025-04-08 11:11:11").StartsInterval("1:1", ksql.Interval_Day_Minute).StartsInterval("2:3", ksql.Interval_Day_Hour)
	e.Ends("2026-04-09 11:11:11").EndsInterval("1", ksql.Interval_Minute).At("2025-04-08 11:11:11").Do(nil)
	assert.Equal(t, "CREATE EVENT `k_event` ON SCHEDULE EVERY + INTERVAL '1:1' HOUR_MINUTE STARTS '2025-04-08 11:11:11' + INTERVAL '1:1' DAY_MINUTE + INTERVAL '2:3' DAY_HOUR ENDS '2026-04-09 11:11:11' + INTERVAL 1 MINUTE ON COMPLETION NOT PRESERVE DO BEGIN INSERT INTO `user` (`name`, `sex`) VALUES (?, ?); END", e.Prepare())
	assert.Equal(t, []any{"kovey", 1}, e.Binds())
	assert.Equal(t, "CREATE EVENT `k_event` ON SCHEDULE EVERY + INTERVAL '1:1' HOUR_MINUTE STARTS '2025-04-08 11:11:11' + INTERVAL '1:1' DAY_MINUTE + INTERVAL '2:3' DAY_HOUR ENDS '2026-04-09 11:11:11' + INTERVAL 1 MINUTE ON COMPLETION NOT PRESERVE DO BEGIN INSERT INTO `user` (`name`, `sex`) VALUES (?, ?); END", e.Prepare())
}

func TestEventAlter(t *testing.T) {
	e := NewEvent().Alter().Event("k_event").Rename("k_new_event").OnCompletionNot()
	assert.Equal(t, "ALTER EVENT `k_event` ON COMPLETION NOT PRESERVE RENAME TO `k_new_event`", e.Prepare())
	assert.Nil(t, e.Binds())
	assert.Equal(t, "ALTER EVENT `k_event` ON COMPLETION NOT PRESERVE RENAME TO `k_new_event`", e.Prepare())
}

func TestEventAlterHasSchedule(t *testing.T) {
	e := NewEvent().Alter().Event("k_event").Rename("k_new_event").OnCompletionNot().At("2025-04-08 11:11:11")
	assert.Equal(t, "ALTER EVENT `k_event` ON SCHEDULE AT '2025-04-08 11:11:11' ON COMPLETION NOT PRESERVE RENAME TO `k_new_event`", e.Prepare())
	assert.Nil(t, e.Binds())
	assert.Equal(t, "ALTER EVENT `k_event` ON SCHEDULE AT '2025-04-08 11:11:11' ON COMPLETION NOT PRESERVE RENAME TO `k_new_event`", e.Prepare())
}
