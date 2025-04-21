package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestConnectionAttr(t *testing.T) {
	testDb, _, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	assert.Equal(t, "mysql", conn.DriverName())
	assert.Equal(t, "mysql", conn.Clone().DriverName())
	assert.NotNil(t, conn.Database())
}

func TestConnectionCommit(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	in := NewInsert()
	in.Table("user").Add("name", "kovey").Add("age", 18)
	now := time.Now().Unix()
	up := NewUpdate()
	up.Table("user_ext").Set("last_time", now).Where(NewWhere().Where("id", ksql.Eq, 1))
	del := NewDelete()
	del.Table("email").Where(NewWhere().Where("id", ksql.Eq, 1))
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO `user` \\(`name`, `age`\\) VALUES \\(\\?, \\?\\)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare(del.Prepare()).ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("UPDATE `user_ext` SET `last_time` = \\? WHERE `id` = \\?").ExpectExec().WithArgs(now, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.Background()
	err = conn.Transaction(ctx, func(ctx context.Context, conn ksql.ConnectionInterface) error {
		if _, err := conn.Insert(ctx, in); err != nil {
			return err
		}

		if _, err := conn.Delete(ctx, del); err != nil {
			return err
		}

		_, err := conn.Update(ctx, up)
		return err
	})

	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
}

func TestConnectionRollback(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	in := NewInsert()
	in.Table("user").Add("name", "kovey").Add("age", 18)
	now := time.Now().Unix()
	up := NewUpdate()
	up.Table("user_ext").Set("last_time", now).Where(NewWhere().Where("id", ksql.Eq, 1))
	expErr := errors.New("update error")
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO `user` \\(`name`, `age`\\) VALUES \\(\\?, \\?\\)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("UPDATE `user_ext` SET `last_time` = \\? WHERE `id` = \\?").ExpectExec().WithArgs(now, 1).WillReturnResult(sqlmock.NewErrorResult(expErr))
	mock.ExpectRollback()

	ctx := context.Background()
	err = conn.Transaction(ctx, func(ctx context.Context, conn ksql.ConnectionInterface) error {
		if _, err := conn.Insert(ctx, in); err != nil {
			return err
		}

		_, err := conn.Update(ctx, up)
		return err
	})
	assert.Equal(t, expErr, err.(*TxErr).callErr.(*SqlErr).Err)
}

func TestConnectionCommitMulti(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	in := NewInsert()
	in.Table("user").Add("name", "kovey").Add("age", 18)
	now := time.Now().Unix()
	up := NewUpdate()
	up.Table("user_ext").Set("last_time", now).Where(NewWhere().Where("id", ksql.Eq, 1))
	del := NewDelete()
	del.Table("email").Where(NewWhere().Where("id", ksql.Eq, 1))
	expErr := errors.New("update error")
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO `user` \\(`name`, `age`\\) VALUES \\(\\?, \\?\\)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare(del.Prepare()).ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("SAVEPOINT ?").ExpectExec().WithArgs("trans_1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("UPDATE `user_ext` SET `last_time` = \\? WHERE `id` = \\?").ExpectExec().WithArgs(now, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("RELEASE SAVEPOINT ?").ExpectExec().WithArgs("trans_1").WillReturnResult(sqlmock.NewErrorResult(expErr))
	mock.ExpectCommit()

	ctx := context.Background()
	err = conn.Transaction(ctx, func(ctx context.Context, conn ksql.ConnectionInterface) error {
		if _, err := conn.Insert(ctx, in); err != nil {
			return err
		}

		if _, err := conn.Delete(ctx, del); err != nil {
			return err
		}

		return conn.Transaction(ctx, func(ctx context.Context, conn ksql.ConnectionInterface) error {
			_, err := conn.Update(ctx, up)
			return err
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
}

func TestConnectionRollbackMulti(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	in := NewInsert()
	in.Table("user").Add("name", "kovey").Add("age", 18)
	now := time.Now().Unix()
	up := NewUpdate()
	up.Table("user_ext").Set("last_time", now).Where(NewWhere().Where("id", ksql.Eq, 1))
	del := NewDelete()
	del.Table("email").Where(NewWhere().Where("id", ksql.Eq, 1))
	expErr := errors.New("update error")
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO `user` \\(`name`, `age`\\) VALUES \\(\\?, \\?\\)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare(del.Prepare()).ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("SAVEPOINT \\?").ExpectExec().WithArgs("trans_1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("UPDATE `user_ext` SET `last_time` = \\? WHERE `id` = \\?").ExpectExec().WithArgs(now, 1).WillReturnResult(sqlmock.NewErrorResult(expErr))
	mock.ExpectPrepare("ROLLBACK TO SAVEPOINT \\?").ExpectExec().WithArgs("trans_1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()

	ctx := context.Background()
	err = conn.Transaction(ctx, func(ctx context.Context, conn ksql.ConnectionInterface) error {
		if _, err := conn.Insert(ctx, in); err != nil {
			return err
		}

		if _, err := conn.Delete(ctx, del); err != nil {
			return err
		}

		return conn.Transaction(ctx, func(ctx context.Context, conn ksql.ConnectionInterface) error {
			_, err := conn.Update(ctx, up)
			return err
		})
	})

	assert.Equal(t, expErr, err.(*TxErr).callErr.(*TxErr).callErr.(*SqlErr).Err)
}

func TestConnectionQuery(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	query := NewQuery()
	query.Table("user").Columns(columns...).Where("id", ksql.Eq, 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = \\?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = conn.QueryRow(context.Background(), query, tu)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestConnectionQueryRaw(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	query := Raw("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ?", 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = \\?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = conn.QueryRowRaw(context.Background(), query, tu)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestConnectionScan(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	query := NewQuery()
	query.Table("user").Columns(columns...).Where("id", ksql.Eq, 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = \\?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = conn.Scan(context.Background(), query, tu.Values()...)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestConnectionScanRaw(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	query := Raw("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ?", 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = \\?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = conn.ScanRaw(context.Background(), query, tu.Values()...)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}
