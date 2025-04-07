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
	a.AddUnique("uni_name_age", "user_name", "age")

	assert.Nil(t, a.AddColumn("test_name", "testname", 10, 0))
	assert.Equal(t, "ALTER TABLE `user` DROP COLUMN `age`,DROP INDEX `idx_name`,ADD COLUMN `user_name` VARCHAR(62) NULL DEFAULT NULL COMMENT '用户名',ADD COLUMN `balance` DECIMAL(10,2) NOT NULL DEFAULT '0' COMMENT '余额',ADD UNIQUE INDEX `user_name` (`user_name`),ADD PRIMARY INDEX (`id`),ADD UNIQUE INDEX `uni_name_age` (`user_name`,`age`),COMMENT = '用户表'", a.Prepare())
	assert.Nil(t, a.Binds())
	assert.Equal(t, "ALTER TABLE `user` DROP COLUMN `age`,DROP INDEX `idx_name`,ADD COLUMN `user_name` VARCHAR(62) NULL DEFAULT NULL COMMENT '用户名',ADD COLUMN `balance` DECIMAL(10,2) NOT NULL DEFAULT '0' COMMENT '余额',ADD UNIQUE INDEX `user_name` (`user_name`),ADD PRIMARY INDEX (`id`),ADD UNIQUE INDEX `uni_name_age` (`user_name`,`age`),COMMENT = '用户表'", a.Prepare())
}

func TestAlterChangeColumn(t *testing.T) {
	a := NewAlter()
	a.Table("user").ChangeColumn("user_name", "nickname", "varchar", 20, 0).Default("").Comment("昵称")
	a.Table("user").ChangeColumn("user_age", "age", "int", 10, 0).Default("").Comment("年龄").AutoIncrement().Unsigned()
	a.Table("user").ModifyColumn("other", "int", 10, 0)
	assert.Equal(t, "ALTER TABLE `user` CHANGE COLUMN `user_name` `nickname` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '昵称',CHANGE COLUMN `user_age` `age` INT(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '年龄',MODIFY COLUMN `other` INT(10) NOT NULL", a.Prepare())
	assert.Nil(t, a.Binds())
}

func TestAlterDropColumnIfExists(t *testing.T) {
	a := NewAlter()
	a.Table("user").DropColumnIfExists("user_name")
	assert.Equal(t, "ALTER TABLE `user` DROP COLUMN IF EXISTS `user_name`", a.Prepare())
	assert.Nil(t, a.Binds())
}

func TestAlterAddColumn(t *testing.T) {
	a := NewAlter()
	a.Table("user").AddDecimal("balance", 20, 2).Unsigned().Default("0").Comment("余额")
	a.AddDouble("rate", 20, 2).Unsigned().Comment("比率")
	a.AddFloat("radio", 20, 2).Unsigned().Comment("radio")
	a.AddBinary("bin", 20).Nullable().Comment("字节")
	a.AddGeoMetry("geo_metry").Nullable().Comment("geo metry")
	a.AddPolygon("p_polygon").Nullable().Comment("polygon")
	a.AddPoint("p_point").Nullable().Comment("point")
	a.AddLineString("line_string").Nullable().Comment("line string")
	a.AddBlob("b_blob").Nullable().Comment("blob")
	a.AddText("content").Nullable().Comment("text content")
	a.AddSet("s_set", []string{"1", "3", "5"}).Nullable().Comment("set")
	a.AddEnum("e_enum", []string{"A", "B", "C"}).Nullable().Comment("enum")
	a.AddDate("birthday").Nullable().Comment("birthday")
	a.AddDateTime("birthday_time").Nullable().Comment("birthday time")
	a.AddTimestamp("create_time").Nullable().Comment("create time")
	a.AddSmallInt("status").Default("1").Comment("status")
	a.AddTinyInt("state").Default("1").Comment("state")
	a.AddBigInt("money").Default("0").Unsigned().Comment("money")
	a.AddInt("day_count").Default("0").Unsigned().Comment("day count")
	a.AddString("name", 225).Default("").Comment("name")
	a.AddChar("result", 10).Default("").Comment("char")
	assert.Equal(t, "ALTER TABLE `user` ADD COLUMN `balance` DECIMAL(20,2) unsigned NOT NULL DEFAULT '0' COMMENT '余额',ADD COLUMN `rate` DOUBLE(20,2) unsigned NOT NULL COMMENT '比率',ADD COLUMN `radio` FLOAT(20,2) unsigned NOT NULL COMMENT 'radio',ADD COLUMN `bin` BINARY(20) NULL COMMENT '字节',ADD COLUMN `geo_metry` GEOMETRY NULL COMMENT 'geo metry',ADD COLUMN `p_polygon` POLYGON NULL COMMENT 'polygon',ADD COLUMN `p_point` POINT NULL COMMENT 'point',ADD COLUMN `line_string` LINESTRING NULL COMMENT 'line string',ADD COLUMN `b_blob` BLOB NULL COMMENT 'blob',ADD COLUMN `content` TEXT NULL COMMENT 'text content',ADD COLUMN `s_set` SET('1','3','5') NULL COMMENT 'set',ADD COLUMN `e_enum` ENUM('A','B','C') NULL COMMENT 'enum',ADD COLUMN `birthday` DATE NULL COMMENT 'birthday',ADD COLUMN `birthday_time` DATETIME(0) NULL COMMENT 'birthday time',ADD COLUMN `create_time` TIMESTAMP(0) NULL COMMENT 'create time',ADD COLUMN `status` SMALLINT(3) NOT NULL DEFAULT '1' COMMENT 'status',ADD COLUMN `state` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'state',ADD COLUMN `money` BIGINT(20) unsigned NOT NULL DEFAULT '0' COMMENT 'money',ADD COLUMN `day_count` INT(11) unsigned NOT NULL DEFAULT '0' COMMENT 'day count',ADD COLUMN `name` VARCHAR(225) NOT NULL DEFAULT '' COMMENT 'name',ADD COLUMN `result` CHAR(10) NOT NULL DEFAULT '' COMMENT 'char'", a.Prepare())
	assert.Nil(t, a.Binds())
}

func TestAlterChangeTable(t *testing.T) {
	a := NewAlter()
	a.Table("user").Charset("utf8").Collate("utf8_general_ci").Engine("InnoDB").Comment("user table")
	assert.Equal(t, "ALTER TABLE `user` CHARACTER SET = 'utf8',COLLATE = 'utf8_general_ci',ENGINE = 'InnoDB',COMMENT = 'user table'", a.Prepare())
	assert.Nil(t, a.Binds())
}

func TestAlterRename(t *testing.T) {
	a := NewAlter()
	a.Table("user").Rename("users")
	assert.Equal(t, "ALTER TABLE `user` RENAME AS `users`", a.Prepare())
	assert.Nil(t, a.Binds())
}
