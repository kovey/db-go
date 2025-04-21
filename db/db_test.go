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

func TestDbInsertFrom(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	mock.ExpectPrepare("INSERT INTO `user` \\(`name`, `age`\\) SELECT `name`, `age` FROM `users` WHERE `id` > \\?").ExpectExec().WithArgs(100).WillReturnResult(sqlmock.NewResult(1, 1))
	query := NewQuery().Table("users").Columns("name", "age").Where("id", ">", 100)
	n, err := InsertFrom(context.Background(), "user", []string{"name", "age"}, query)
	assert.Equal(t, int64(1), n)
	assert.Nil(t, err)
}

func TestDbInsertFromErr(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	delErr := errors.New("del error")
	mock.ExpectPrepare("INSERT INTO `user` \\(`name`, `age`\\) SELECT `name`, `age` FROM `users` WHERE `id` > \\?").ExpectExec().WithArgs(100).WillReturnResult(sqlmock.NewErrorResult(delErr))
	query := NewQuery().Table("users").Columns("name", "age").Where("id", ">", 100)
	n, err := InsertFrom(context.Background(), "user", []string{"name", "age"}, query)
	assert.Equal(t, int64(0), n)
	assert.Equal(t, delErr, err.(*SqlErr).Err)
}

func TestDbQuery(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	var rows []*test_user
	query := NewQuery()
	query.Table("user").Columns(columns...).Order("id").OrderDesc("age").Where("id", ksql.Ge, 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` >= \\? ORDER BY `id` ASC, `age` DESC").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23).AddRow(2, 15, "test", now, 34.13))
	err = Query(context.Background(), query, &rows)
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

func TestDbQueryRow(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	query := NewQuery()
	query.Table("user").Columns(columns...).Where("id", ksql.Eq, 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = \\?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = QueryRow(context.Background(), query, tu)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestDbFind(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = \\?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = Find(context.Background(), tu, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestDbLock(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = \\? FOR UPDATE").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	mock.ExpectPrepare("INSERT INTO `user` \\(`name`, `age`\\) VALUES \\(\\?, \\?\\)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("UPDATE `user` SET `balance` = `balance` \\+ \\? WHERE `id` = \\?").ExpectExec().WithArgs(100, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	up := NewUpdate()
	up.Table("user").IncColumn("balance", 100).Where(NewWhere().Where("id", ksql.Eq, 1))
	in := NewInsert()
	in.Table("user").Add("name", "kovey").Add("age", 18)
	err = Transaction(context.Background(), func(ctx context.Context, db ksql.ConnectionInterface) error {
		tu := newTestUser()
		if err := Lock(ctx, db, tu, 1); err != nil {
			return err
		}

		if _, err := db.Insert(ctx, in); err != nil {
			return err
		}

		_, err := db.Update(ctx, up)
		return err
	})

	if err != nil {
		t.Fatal(err.(*TxErr).callErr)
	}

	assert.Nil(t, err)
}
