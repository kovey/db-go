package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestAlter(t *testing.T) {
	a := NewAlter()
	a.Table("user").AddColumn("user_name", "VARCHAR", 62, 0).Nullable().Default("NULL").Comment("用户名")
	a.DropColumn("age").AddColumn("balance", "DECIMAL", 10, 2).Default("0").Comment("余额")
	a.AddIndex("user_name", ksql.Index_Type_Unique, "user_name")
	a.DropIndex("idx_name").Comment("用户表").AddPrimary("id")

	assert.Equal(t, "ALTER TABLE `user` DROP COLUMN `age`,DROP INDEX `idx_name`,ADD COLUMN `user_name` VARCHAR(62)  NULL  DEFAULT NULL COMMENT '用户名',ADD COLUMN `balance` DECIMAL(10,2)  NOT NULL  DEFAULT '0' COMMENT '余额',ADD UNIQUE INDEX user_name (`user_name`),ADD PRIMARY INDEX (`id`),COMMENT = '用户表'", a.Prepare())
	assert.Nil(t, a.Binds())
	assert.Equal(t, "ALTER TABLE `user` DROP COLUMN `age`,DROP INDEX `idx_name`,ADD COLUMN `user_name` VARCHAR(62)  NULL  DEFAULT NULL COMMENT '用户名',ADD COLUMN `balance` DECIMAL(10,2)  NOT NULL  DEFAULT '0' COMMENT '余额',ADD UNIQUE INDEX user_name (`user_name`),ADD PRIMARY INDEX (`id`),COMMENT = '用户表'", a.Prepare())
}
