package express

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestSelectExpress(t *testing.T) {
	s := NewStatement("SELECT * FROM user WHERE id > ?", []any{1})
	assert.Equal(t, ksql.Sql_Type_Query, s.Type())
	assert.Equal(t, "SELECT * FROM user WHERE id > ?", s.Statement())
	assert.Equal(t, []any{1}, s.Binds())
	assert.False(t, s.IsExec())
}

func TestInsertExpress(t *testing.T) {
	s := NewStatement("INSERT INTO user (id, name, age) VALUES (?, ?, ?)", []any{1, "kovey", 18})
	assert.Equal(t, ksql.Sql_Type_Insert, s.Type())
	assert.Equal(t, "INSERT INTO user (id, name, age) VALUES (?, ?, ?)", s.Statement())
	assert.Equal(t, []any{1, "kovey", 18}, s.Binds())
	assert.True(t, s.IsExec())
}

func TestUpdateExpress(t *testing.T) {
	s := NewStatement("UPDATE user SET name = ? WHERE id = ?", []any{"kovey", 1})
	assert.Equal(t, ksql.Sql_Type_Update, s.Type())
	assert.Equal(t, "UPDATE user SET name = ? WHERE id = ?", s.Statement())
	assert.Equal(t, []any{"kovey", 1}, s.Binds())
	assert.True(t, s.IsExec())
}

func TestDeteleExpress(t *testing.T) {
	s := NewStatement("DELETE FROM user WHERE id = ?", []any{1})
	assert.Equal(t, ksql.Sql_Type_Delete, s.Type())
	assert.Equal(t, "DELETE FROM user WHERE id = ?", s.Statement())
	assert.Equal(t, []any{1}, s.Binds())
	assert.True(t, s.IsExec())
}

func TestShowExpress(t *testing.T) {
	s := NewStatement("SHOW CREATE TABLE user", nil)
	assert.Equal(t, ksql.Sql_Type_Query, s.Type())
	assert.Equal(t, "SHOW CREATE TABLE user", s.Statement())
	assert.Nil(t, s.Binds())
	assert.False(t, s.IsExec())
}

func TestCreateTableExpress(t *testing.T) {
	s := NewStatement("CREATE TABLE user LIKE other", nil)
	assert.Equal(t, ksql.Sql_Type_Create, s.Type())
	assert.Equal(t, "CREATE TABLE user LIKE other", s.Statement())
	assert.Nil(t, s.Binds())
	assert.True(t, s.IsExec())
}

func TestDropTableExpress(t *testing.T) {
	s := NewStatement("DROP TABLE IF EXISTS user", nil)
	assert.Equal(t, ksql.Sql_Type_Drop, s.Type())
	assert.Equal(t, "DROP TABLE IF EXISTS user", s.Statement())
	assert.Nil(t, s.Binds())
	assert.True(t, s.IsExec())
}

func TestAlterTableExpress(t *testing.T) {
	s := NewStatement("ALTER TABLE `user` DROP COLUMN `age`,DROP INDEX `idx_name`,ADD COLUMN `user_name` VARCHAR(62)  NULL  DEFAULT NULL COMMENT '用户名',ADD COLUMN `balance` DECIMAL(10,2)  NOT NULL  DEFAULT '0' COMMENT '余额',ADD UNIQUE INDEX user_name (`user_name`),ADD PRIMARY INDEX (`id`),COMMENT = '用户表'", nil)
	assert.Equal(t, ksql.Sql_Type_Alter, s.Type())
	assert.Equal(t, "ALTER TABLE `user` DROP COLUMN `age`,DROP INDEX `idx_name`,ADD COLUMN `user_name` VARCHAR(62)  NULL  DEFAULT NULL COMMENT '用户名',ADD COLUMN `balance` DECIMAL(10,2)  NOT NULL  DEFAULT '0' COMMENT '余额',ADD UNIQUE INDEX user_name (`user_name`),ADD PRIMARY INDEX (`id`),COMMENT = '用户表'", s.Statement())
	assert.Nil(t, s.Binds())
	assert.True(t, s.IsExec())
}
