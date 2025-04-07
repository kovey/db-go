package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	ta := NewTable()
	ta.Table("user").AddColumn("id", "bigint", 20, 0).AutoIncrement().Unsigned().Comment("主键")
	ta.AddColumn("username", "VARCHAR", 31, 0).Nullable().Default("NULL").Comment("用户名")
	ta.AddColumn("password", "VARCHAR", 64, 0).Default("").Comment("密码")
	ta.AddColumn("age", "int", 11, 0).Default("0").Comment("密码")
	ta.AddColumn("create_time", "TIMESTAMP", 0, 0).Default(ksql.CURRENT_TIMESTAMP).Comment("创建时间")
	ta.AddColumn("update_time", "TIMESTAMP", 0, 0).Default(ksql.CURRENT_TIMESTAMP_ON_UPDATE_CURRENT_TIMESTAMP).Comment("更新时间")
	ta.AddPrimary("id").AddIndex("idx_username", ksql.Index_Type_Unique, "username").AddIndex("idx_name_age", ksql.Index_Type_Normal, "username", "age")
	ta.Engine("InnoDB").Charset("utf8").Collate("test").Comment("用户表")

	assert.Equal(t, "CREATE TABLE `user` (`id` BIGINT(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',`username` VARCHAR(31) NULL DEFAULT NULL COMMENT '用户名',`password` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '密码',`age` INT(11) NOT NULL DEFAULT '0' COMMENT '密码',`create_time` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',`update_time` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',PRIMARY KEY (`id`),UNIQUE KEY `idx_username` (`username`),KEY `idx_name_age` (`username`,`age`)) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=test COMMENT='用户表'", ta.Prepare())
	assert.Nil(t, ta.Binds())
	assert.Equal(t, "CREATE TABLE `user` (`id` BIGINT(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',`username` VARCHAR(31) NULL DEFAULT NULL COMMENT '用户名',`password` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '密码',`age` INT(11) NOT NULL DEFAULT '0' COMMENT '密码',`create_time` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',`update_time` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',PRIMARY KEY (`id`),UNIQUE KEY `idx_username` (`username`),KEY `idx_name_age` (`username`,`age`)) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=test COMMENT='用户表'", ta.Prepare())
}

func TestTableLike(t *testing.T) {
	ta := NewTable().Table("user").Like("users")
	assert.Equal(t, "CREATE TABLE `user` LIKE `users`", ta.Prepare())
	assert.Nil(t, ta.Binds())
	assert.Equal(t, "CREATE TABLE `user` LIKE `users`", ta.Prepare())
}

func TestTableFrom(t *testing.T) {
	query := NewQuery()
	query.Table("user_back").Columns("u.name", "u.kovey", "e.date").As("u")
	query.Join("email").As("e").On("e.id", "=", "u.id")
	query.Where("u.name", "LIKE", "%test%")
	ta := NewTable().Table("user").From(query)
	assert.Equal(t, "CREATE TABLE `user` AS SELECT `u`.`name`,`u`.`kovey`,`e`.`date` FROM `user_back` AS `u` JOIN `email` AS `e` ON `e`.`id`=`u`.`id` WHERE `u`.`name` LIKE ?", ta.Prepare())
	assert.Equal(t, []any{"%test%"}, ta.Binds())
	assert.Equal(t, "CREATE TABLE `user` AS SELECT `u`.`name`,`u`.`kovey`,`e`.`date` FROM `user_back` AS `u` JOIN `email` AS `e` ON `e`.`id`=`u`.`id` WHERE `u`.`name` LIKE ?", ta.Prepare())
}
