package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestProcedure(t *testing.T) {
	p := NewProcedure().Procedure("test")
	p.Definer("'user'@'local'").IfNotExists().In("input1", "VARCHAR(20)").In("input2", "INT").Out("out1", "INT").Comment("test procedure")
	query := Raw("SELECT `id`,`name`,`age` INTO `out1` FROM `user` WHERE `id` >= input2 AND `name` = ?", "input2")
	p.RoutineBody(query)
	assert.Equal(t, "CREATE DEFINER = 'user'@'local' PROCEDURE IF NOT EXISTS `test`(IN `input1` VARCHAR(20), IN `input2` INT, OUT `out1` INT) COMMENT 'test procedure' BEGIN SELECT `id`,`name`,`age` INTO `out1` FROM `user` WHERE `id` >= input2 AND `name` = input2; END", p.Prepare())
	assert.Nil(t, p.Binds())
}

func TestProcedureNoArgs(t *testing.T) {
	query := Raw("SELECT `id`,`name`,`age` FROM `user` WHERE `id` >= 1")
	p := NewProcedure().Procedure("test").RoutineBody(query)
	assert.Equal(t, "CREATE PROCEDURE `test`() BEGIN SELECT `id`,`name`,`age` FROM `user` WHERE `id` >= 1; END", p.Prepare())
	assert.Nil(t, p.Binds())
	assert.Equal(t, "CREATE PROCEDURE `test`() BEGIN SELECT `id`,`name`,`age` FROM `user` WHERE `id` >= 1; END", p.Prepare())
}

func TestProcedureRaw(t *testing.T) {
	p := NewProcedure().Procedure("test")
	p.InOut("inout1", "BIGINT").InOut("inout2", "TEXT").Deterministic().Language().SqlType(ksql.Procedure_Sql_Type_Contains_Sql).SqlSecurity(ksql.Sql_Security_Definer)
	raw := Raw("SELECT id INTO inout1, name INTO inout2 WHERE user_id > inout1 AND name LIKE ?", "inout2")
	p.RoutineBody(raw)
	assert.Equal(t, "CREATE PROCEDURE `test`(INOUT `inout1` BIGINT, INOUT `inout2` TEXT) LANGUAGE SQL DETERMINISTIC CONTAINS SQL SQL SECURITY DEFINER BEGIN SELECT id INTO inout1, name INTO inout2 WHERE user_id > inout1 AND name LIKE inout2; END", p.Prepare())
	assert.Nil(t, p.Binds())
}

func TestProcedureRawNoArgs(t *testing.T) {
	raw := Raw("SELECT * WHERE user_id > 1")
	ins := Raw("DELETE FROM user WHERE id = 1")
	p := NewProcedure().Procedure("test").RoutineBody(raw).DeterministicNot().RoutineBody(ins)
	assert.Equal(t, "CREATE PROCEDURE `test`() NOT DETERMINISTIC BEGIN SELECT * WHERE user_id > 1; END", p.Prepare())
	assert.Nil(t, p.Binds())
	assert.Equal(t, "CREATE PROCEDURE `test`() NOT DETERMINISTIC BEGIN SELECT * WHERE user_id > 1; END", p.Prepare())
}

func TestProcedureAlter(t *testing.T) {
	raw := Raw("SELECT * WHERE user_id > 1")
	ins := Raw("DELETE FROM user WHERE id = 1")
	p := NewProcedure().Alter().Procedure("test").RoutineBody(raw).DeterministicNot().RoutineBody(ins).Deterministic().Comment("alter").Definer("root").IfNotExists()
	p.In("input1", "INT").Out("out1", "INT").InOut("inout1", "BIGINT").Language().SqlSecurity(ksql.Sql_Security_Definer).SqlType(ksql.Procedure_Sql_Type_Contains_Sql)
	assert.Equal(t, "ALTER PROCEDURE `test` COMMENT 'alter' LANGUAGE SQL CONTAINS SQL SQL SECURITY DEFINER", p.Prepare())
	assert.Nil(t, p.Binds())
}
