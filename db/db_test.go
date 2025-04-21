package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql"
	"github.com/stretchr/testify/assert"
)

func TestDbInsert(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	mock.ExpectPrepare("INSERT INTO `user` \\(`name`, `age`\\) VALUES \\(\\?, \\?\\)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	dt := NewData()
	dt.Set("name", "kovey")
	dt.Set("age", 18)
	n, err := Insert(context.Background(), "user", dt)
	assert.Equal(t, int64(1), n)
	assert.Nil(t, err)
}

func TestDbInsertErr(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	insErr := errors.New("insert error")
	mock.ExpectPrepare("INSERT INTO `user` \\(`name`, `age`\\) VALUES \\(\\?, \\?\\)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewErrorResult(insErr))
	dt := NewData()
	dt.Set("name", "kovey")
	dt.Set("age", 18)
	n, err := Insert(context.Background(), "user", dt)
	assert.Equal(t, int64(0), n)
	assert.Equal(t, insErr, err.(*SqlErr).Err)
}

func TestDbUpdate(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	now := time.Now().Unix()
	mock.ExpectPrepare("UPDATE `user_ext` SET `last_time` = \\? WHERE `id` = \\?").ExpectExec().WithArgs(now, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	dt := NewData()
	dt.Set("last_time", now)
	where := sql.NewWhere()
	where.Where("id", ksql.Eq, 1)
	n, err := Update(context.Background(), "user_ext", dt, where)
	assert.Equal(t, int64(1), n)
	assert.Nil(t, err)
}

func TestDbUpdateErr(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	upErr := errors.New("update error")
	now := time.Now().Unix()
	mock.ExpectPrepare("UPDATE `user_ext` SET `last_time` = \\? WHERE `id` = \\?").ExpectExec().WithArgs(now, 1).WillReturnResult(sqlmock.NewErrorResult(upErr))
	dt := NewData()
	dt.Set("last_time", now)
	where := sql.NewWhere()
	where.Where("id", ksql.Eq, 1)
	n, err := Update(context.Background(), "user_ext", dt, where)
	assert.Equal(t, int64(0), n)
	assert.Equal(t, upErr, err.(*SqlErr).Err)
}

func TestDbDelete(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	mock.ExpectPrepare("DELETE FROM `user_ext` WHERE `id` = \\?").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	where := sql.NewWhere()
	where.Where("id", ksql.Eq, 1)
	n, err := Delete(context.Background(), "user_ext", where)
	assert.Equal(t, int64(1), n)
	assert.Nil(t, err)
}

func TestDbDeleteErr(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	delErr := errors.New("del error")
	mock.ExpectPrepare("DELETE FROM `user_ext` WHERE `id` = \\?").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(delErr))
	where := sql.NewWhere()
	where.Where("id", ksql.Eq, 1)
	n, err := Delete(context.Background(), "user_ext", where)
	assert.Equal(t, int64(0), n)
	assert.Equal(t, delErr, err.(*SqlErr).Err)
}
