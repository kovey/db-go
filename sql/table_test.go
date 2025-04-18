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
	ta.AddPrimary("id")
	ta.AddIndex("idx_username").Type(ksql.Index_Type_Unique).Columns("username")
	ta.AddIndex("idx_name_age").Type(ksql.Index_Type_Normal).Columns("username", "age")
	ta.Engine("InnoDB").Charset("utf8").Collate("test").Comment("用户表")

	assert.Equal(t, "CREATE TABLE `user` (`id` BIGINT(20) UNSIGNED AUTO_INCREMENT COMMENT '主键', `username` VARCHAR(31) NULL DEFAULT NULL COMMENT '用户名', `password` VARCHAR(64) DEFAULT '' COMMENT '密码', `age` INT(11) DEFAULT '0' COMMENT '密码', `create_time` TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间', `update_time` TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',  PRIMARY KEY (`id`),  INDEX `idx_username` (`username`),  INDEX `idx_name_age` (`username`, `age`)) ENGINE = InnoDB, CHARACTER SET = utf8, COLLATE = test, COMMENT = '用户表'", ta.Prepare())
	assert.Nil(t, ta.Binds())
	assert.Equal(t, "CREATE TABLE `user` (`id` BIGINT(20) UNSIGNED AUTO_INCREMENT COMMENT '主键', `username` VARCHAR(31) NULL DEFAULT NULL COMMENT '用户名', `password` VARCHAR(64) DEFAULT '' COMMENT '密码', `age` INT(11) DEFAULT '0' COMMENT '密码', `create_time` TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间', `update_time` TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',  PRIMARY KEY (`id`),  INDEX `idx_username` (`username`),  INDEX `idx_name_age` (`username`, `age`)) ENGINE = InnoDB, CHARACTER SET = utf8, COLLATE = test, COMMENT = '用户表'", ta.Prepare())
}

func TestTableLike(t *testing.T) {
	ta := NewTable().Table("user").Like("users")
	assert.Equal(t, "CREATE TABLE `user` (LIKE `users`)", ta.Prepare())
	assert.Nil(t, ta.Binds())
	assert.Equal(t, "CREATE TABLE `user` (LIKE `users`)", ta.Prepare())
}

func TestTableFrom(t *testing.T) {
	query := NewQuery()
	query.Table("user_back").Columns("u.name", "u.kovey", "e.date").As("u")
	query.Join("email").As("e").On("e.id", "=", "u.id")
	query.Where("u.name", "LIKE", "%test%")
	ta := NewTable().Table("user").As(query)
	assert.Equal(t, "CREATE TABLE `user` AS SELECT `u`.`name`, `u`.`kovey`, `e`.`date` FROM `user_back` AS `u` INNER JOIN `email` AS `e` ON (`e`.`id` = `u`.`id`) WHERE `u`.`name` LIKE ?", ta.Prepare())
	assert.Equal(t, []any{"%test%"}, ta.Binds())
	assert.Equal(t, "CREATE TABLE `user` AS SELECT `u`.`name`, `u`.`kovey`, `e`.`date` FROM `user_back` AS `u` INNER JOIN `email` AS `e` ON (`e`.`id` = `u`.`id`) WHERE `u`.`name` LIKE ?", ta.Prepare())
}

func TestTableAddColumn(t *testing.T) {
	a := NewTable()
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
	a.AddUnique("uni_name_email", "bin", "geo_metry")
	assert.Equal(t, "CREATE TABLE `user` (`balance` DECIMAL(20,2) UNSIGNED DEFAULT '0' COMMENT '余额', `rate` DOUBLE(20,2) UNSIGNED COMMENT '比率', `radio` FLOAT(20,2) UNSIGNED COMMENT 'radio', `bin` BINARY(20) NULL COMMENT '字节', `geo_metry` GEOMETRY NULL COMMENT 'geo metry', `p_polygon` POLYGON NULL COMMENT 'polygon', `p_point` POINT NULL COMMENT 'point', `line_string` LINESTRING NULL COMMENT 'line string', `b_blob` BLOB NULL COMMENT 'blob', `content` TEXT NULL COMMENT 'text content', `s_set` SET('1','3','5') NULL COMMENT 'set', `e_enum` ENUM('A','B','C') NULL COMMENT 'enum', `birthday` DATE NULL COMMENT 'birthday', `birthday_time` DATETIME(0) NULL COMMENT 'birthday time', `create_time` TIMESTAMP(0) NULL COMMENT 'create time', `status` SMALLINT(3) DEFAULT '1' COMMENT 'status', `state` TINYINT(1) DEFAULT '1' COMMENT 'state', `money` BIGINT(20) UNSIGNED DEFAULT '0' COMMENT 'money', `day_count` INT(11) UNSIGNED DEFAULT '0' COMMENT 'day count', `name` VARCHAR(225) DEFAULT '' COMMENT 'name', `result` CHAR(10) DEFAULT '' COMMENT 'char',  UNIQUE INDEX `uni_name_email` (`bin`, `geo_metry`))", a.Prepare())
	assert.Nil(t, a.Binds())
}

func TestTableAddInvalidColumn(t *testing.T) {
	ta := NewTable()
	assert.Nil(t, ta.Table("user").AddColumn("id", "testtype", 10, 20))
	assert.Nil(t, ta.Binds())
	assert.Equal(t, "CREATE TABLE `user`", ta.Prepare())
}
