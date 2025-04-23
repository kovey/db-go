package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRawInsert(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	n, err := InsertRaw(context.Background(), Raw("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)", "kovey", 18))
	assert.Equal(t, int64(1), n)
	assert.Nil(t, err)
}

func TestRawInsertErr(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	insErr := errors.New("insert error")
	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewErrorResult(insErr))
	n, err := InsertRaw(context.Background(), Raw("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)", "kovey", 18))
	assert.Equal(t, int64(0), n)
	assert.Equal(t, insErr, err.(*SqlErr).Err)
}

func TestRawUpdate(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	now := time.Now().Unix()
	mock.ExpectPrepare("UPDATE `user_ext` SET `last_time` = ? WHERE `id` = ?").ExpectExec().WithArgs(now, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	n, err := UpdateRaw(context.Background(), Raw("UPDATE `user_ext` SET `last_time` = ? WHERE `id` = ?", now, 1))
	assert.Equal(t, int64(1), n)
	assert.Nil(t, err)
}

func TestRawUpdateErr(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	upErr := errors.New("update error")
	now := time.Now().Unix()
	mock.ExpectPrepare("UPDATE `user_ext` SET `last_time` = ? WHERE `id` = ?").ExpectExec().WithArgs(now, 1).WillReturnResult(sqlmock.NewErrorResult(upErr))
	n, err := UpdateRaw(context.Background(), Raw("UPDATE `user_ext` SET `last_time` = ? WHERE `id` = ?", now, 1))
	assert.Equal(t, int64(0), n)
	assert.Equal(t, upErr, err.(*SqlErr).Err)
}

func TestRawDelete(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	mock.ExpectPrepare("DELETE FROM `user_ext` WHERE `id` = ?").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	n, err := DeleteRaw(context.Background(), Raw("DELETE FROM `user_ext` WHERE `id` = ?", 1))
	assert.Equal(t, int64(1), n)
	assert.Nil(t, err)
}

func TestRawDeleteErr(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	delErr := errors.New("del error")
	mock.ExpectPrepare("DELETE FROM `user_ext` WHERE `id` = ?").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(delErr))
	n, err := DeleteRaw(context.Background(), Raw("DELETE FROM `user_ext` WHERE `id` = ?", 1))
	assert.Equal(t, int64(0), n)
	assert.Equal(t, delErr, err.(*SqlErr).Err)
}

func TestRawQuery(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	var rows []*test_user
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` >= ? ORDER BY `id` ASC, `age` DESC").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23).AddRow(2, 15, "test", now, 34.13))
	err = QueryRaw(context.Background(), Raw("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` >= ? ORDER BY `id` ASC, `age` DESC", 1), &rows)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rows))
	tu := rows[0]
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
	tu = rows[1]
	assert.Equal(t, int64(2), tu.Id)
	assert.Equal(t, 15, tu.Age)
	assert.Equal(t, "test", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(34.13), tu.Balance)
}

func TestRawQueryRow(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = QueryRowRaw(context.Background(), Raw("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ?", 1), tu)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestRawHasTable(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	mock.ExpectPrepare("SHOW TABLES LIKE 'user'").ExpectQuery().WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{"table"}).AddRow("user"))
	has, err := HasTable(context.Background(), "user")
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.True(t, has)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRawHasColumn(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	mock.ExpectPrepare("SHOW COLUMNS FROM `user` LIKE 'id'").ExpectQuery().WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{"field", "type", "null", "key", "default", "extra"}).AddRow("id", "int", "NO", "PRI", nil, "auto_increment"))
	has, err := HasColumn(context.Background(), "user", "id")
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.True(t, has)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRawHasIndex(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	mock.ExpectPrepare("SHOW INDEX FROM `user` WHERE Key_name = ?").ExpectQuery().WithArgs("PRIMARY").WillReturnRows(sqlmock.NewRows([]string{"Table", "Non_unique", "Key_name", "Seq_in_index", "Column_name", "Collation", "Cardinality", "Sub_part", "Packed", "Null", "Index_type", "Comment", "Index_comment", "Visible", "Expression"}).AddRow("user", 0, "PRIMARY", 1, "id", "A", 2, nil, nil, "", "BTREE", "", "", "YES", nil))
	has, err := HasIndex(context.Background(), "user", "PRIMARY")
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.True(t, has)
	assert.Nil(t, mock.ExpectationsWereMet())
}
