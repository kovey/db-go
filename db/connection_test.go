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
