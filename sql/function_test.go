package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestFunction(t *testing.T) {
	f := NewFunction().Function("test_func")
	f.Definer("kovey").Comment("test").IfNotExists().Param("arg1", "VARCHAR(20)").Param("arg2", "INT").Returns("VARCHAR")
	raw := Raw("SET @name = '';SELECT name INTO name FROM user WHERE id = 1;RETURN @name")
	f.RoutineBody(raw)
	assert.Equal(t, "CREATE DEFINER = kovey FUNCTION IF NOT EXISTS `test_func`(`arg1` VARCHAR(20), `arg2` INT) RETURNS VARCHAR COMMENT 'test' BEGIN SET @name = '';SELECT name INTO name FROM user WHERE id = 1;RETURN @name; END", f.Prepare())
	assert.Nil(t, f.Binds())
}

func TestFunctionNoArgs(t *testing.T) {
	f := NewFunction().Function("test_func")
	f.Returns("VARCHAR").Deterministic().Language().SqlType(ksql.Procedure_Sql_Type_Contains_Sql).SqlSecurity(ksql.Sql_Security_Definer)
	raw := Raw("SET @name = '';SELECT name INTO name FROM user WHERE id = 1;RETURN @name")
	f.RoutineBody(raw)
	assert.Equal(t, "CREATE FUNCTION `test_func`() RETURNS VARCHAR LANGUAGE SQL DETERMINISTIC CONTAINS SQL SQL SECURITY DEFINER BEGIN SET @name = '';SELECT name INTO name FROM user WHERE id = 1;RETURN @name; END", f.Prepare())
	assert.Nil(t, f.Binds())
}

func TestFunctionNoArgsNot(t *testing.T) {
	f := NewFunction().Function("test_func")
	f.Returns("VARCHAR").DeterministicNot().Language().SqlType(ksql.Procedure_Sql_Type_Contains_Sql).SqlSecurity(ksql.Sql_Security_Definer)
	raw := Raw("SET @name = '';SELECT name INTO name FROM user WHERE id = 1;RETURN @name")
	f.RoutineBody(raw)
	assert.Equal(t, "CREATE FUNCTION `test_func`() RETURNS VARCHAR LANGUAGE SQL NOT DETERMINISTIC CONTAINS SQL SQL SECURITY DEFINER BEGIN SET @name = '';SELECT name INTO name FROM user WHERE id = 1;RETURN @name; END", f.Prepare())
	assert.Nil(t, f.Binds())
}

func TestFunctionAlter(t *testing.T) {
	f := NewFunction().Function("test_func").Alter()
	f.Returns("VARCHAR").DeterministicNot().Language().SqlType(ksql.Procedure_Sql_Type_Contains_Sql).SqlSecurity(ksql.Sql_Security_Definer).Param("test", "INT").Definer("root").IfNotExists().Returns("BIGINT")
	raw := Raw("SET @name = '';SELECT name INTO name FROM user WHERE id = 1;RETURN @name")
	f.RoutineBody(raw)
	assert.Equal(t, "ALTER FUNCTION `test_func` LANGUAGE SQL NOT DETERMINISTIC CONTAINS SQL SQL SECURITY DEFINER", f.Prepare())
	assert.Nil(t, f.Binds())
}
