package sharding

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
	"github.com/stretchr/testify/assert"
)

func TestDbInsert(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()
	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)

	mock1.ExpectPrepare("INSERT INTO `user_0` (`user_id`, `name`, `age`) VALUES (?, ?, ?)").ExpectExec().WithArgs(2, "kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	mock2.ExpectPrepare("INSERT INTO `user_1` (`user_id`, `name`, `age`) VALUES (?, ?, ?)").ExpectExec().WithArgs(1, "kovey1", 19).WillReturnResult(sqlmock.NewResult(1, 1))
	in1 := db.NewData()
	in1.Set("user_id", 2).Set("name", "kovey").Set("age", 18)
	in2 := db.NewData()
	in2.Set("user_id", 1).Set("name", "kovey1").Set("age", 19)
	id1, err := Insert(2, context.Background(), "user_0", in1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id1)
	id2, err := Insert(1, context.Background(), "user_1", in2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestDbInsertErr(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()
	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)

	sqlErr := errors.New("error")
	mock1.ExpectPrepare("INSERT INTO `user_0` (`user_id`, `name`, `age`) VALUES (?, ?, ?)").ExpectExec().WithArgs(2, "kovey", 18).WillReturnResult(sqlmock.NewErrorResult(sqlErr))
	mock2.ExpectPrepare("INSERT INTO `user_1` (`user_id`, `name`, `age`) VALUES (?, ?, ?)").ExpectExec().WithArgs(1, "kovey1", 19).WillReturnResult(sqlmock.NewErrorResult(sqlErr))
	in1 := db.NewData()
	in1.Set("user_id", 2).Set("name", "kovey").Set("age", 18)
	in2 := db.NewData()
	in2.Set("user_id", 1).Set("name", "kovey1").Set("age", 19)
	id1, err := Insert(2, context.Background(), "user_0", in1)
	assert.Equal(t, sqlErr, err.(*db.SqlErr).Err)
	assert.Equal(t, int64(0), id1)
	id2, err := Insert(1, context.Background(), "user_1", in2)
	assert.Equal(t, sqlErr, err.(*db.SqlErr).Err)
	assert.Equal(t, int64(0), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestDbUpdate(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()

	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)
	mock1.ExpectPrepare("UPDATE `user_0` SET `name` = ?, `age` = ? WHERE `user_id` = ?").ExpectExec().WithArgs("kovey", 18, 2).WillReturnResult(sqlmock.NewResult(0, 1))
	mock2.ExpectPrepare("UPDATE `user_1` SET `name` = ?, `age` = ? WHERE `user_id` = ?").ExpectExec().WithArgs("kovey1", 19, 1).WillReturnResult(sqlmock.NewResult(0, 1))
	in1 := db.NewData()
	w1 := db.NewWhere().Where("user_id", "=", 2)
	in1.Set("name", "kovey").Set("age", 18)
	in2 := db.NewData()
	w2 := db.NewWhere().Where("user_id", "=", 1)
	in2.Set("name", "kovey1").Set("age", 19)
	id1, err := Update(2, context.Background(), "user_0", in1, w1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id1)
	id2, err := Update(1, context.Background(), "user_1", in2, w2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestDbUpdateErr(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()

	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)

	sqlErr := errors.New("error")
	mock1.ExpectPrepare("UPDATE `user_0` SET `name` = ?, `age` = ? WHERE `user_id` = ?").ExpectExec().WithArgs("kovey", 18, 2).WillReturnResult(sqlmock.NewErrorResult(sqlErr))
	mock2.ExpectPrepare("UPDATE `user_1` SET `name` = ?, `age` = ? WHERE `user_id` = ?").ExpectExec().WithArgs("kovey1", 19, 1).WillReturnResult(sqlmock.NewErrorResult(sqlErr))
	in1 := db.NewData()
	w1 := db.NewWhere().Where("user_id", "=", 2)
	in1.Set("name", "kovey").Set("age", 18)
	in2 := db.NewData()
	w2 := db.NewWhere().Where("user_id", "=", 1)
	in2.Set("name", "kovey1").Set("age", 19)
	id1, err := Update(2, context.Background(), "user_0", in1, w1)
	assert.Equal(t, sqlErr, err.(*db.SqlErr).Err)
	assert.Equal(t, int64(0), id1)
	id2, err := Update(1, context.Background(), "user_1", in2, w2)
	assert.Equal(t, sqlErr, err.(*db.SqlErr).Err)
	assert.Equal(t, int64(0), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestDbDelete(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()

	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)
	mock1.ExpectPrepare("DELETE FROM `user_0` WHERE `user_id` = ?").ExpectExec().WithArgs(2).WillReturnResult(sqlmock.NewResult(0, 1))
	mock2.ExpectPrepare("DELETE FROM `user_1` WHERE `user_id` = ?").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	w1 := db.NewWhere().Where("user_id", "=", 2)
	w2 := db.NewWhere().Where("user_id", "=", 1)
	id1, err := Delete(2, context.Background(), "user_0", w1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id1)
	id2, err := Delete(1, context.Background(), "user_1", w2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestDbDeleteErr(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()

	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)
	sqlErr := errors.New("error")
	mock1.ExpectPrepare("DELETE FROM `user_0` WHERE `user_id` = ?").ExpectExec().WithArgs(2).WillReturnResult(sqlmock.NewErrorResult(sqlErr))
	mock2.ExpectPrepare("DELETE FROM `user_1` WHERE `user_id` = ?").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(sqlErr))
	w1 := db.NewWhere().Where("user_id", "=", 2)
	w2 := db.NewWhere().Where("user_id", "=", 1)
	id1, err := Delete(2, context.Background(), "user_0", w1)
	assert.Equal(t, sqlErr, err.(*db.SqlErr).Err)
	assert.Equal(t, int64(0), id1)
	id2, err := Delete(1, context.Background(), "user_1", w2)
	assert.Equal(t, sqlErr, err.(*db.SqlErr).Err)
	assert.Equal(t, int64(0), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestDbQuery(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()

	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)
	var tms []*test_model
	var tm1s []*test_model
	mock1.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_0` WHERE `user_id` = ?").ExpectQuery().WithArgs(2).WillReturnRows(sqlmock.NewRows(newTestModel().Columns()).AddRow(1, 2, 18, "kovey"))
	mock2.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_1` WHERE `user_id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(newTestModel().Columns()).AddRow(1, 1, 19, "kovey1"))
	in1 := db.NewQuery()
	in1.Table("user_0").Where("user_id", "=", 2).Columns(newTestModel().Columns()...)
	err = Query(2, context.Background(), in1, &tms)
	if err != nil {
		t.Fatal(err)
	}
	tm := tms[0]
	assert.Nil(t, err)
	assert.Equal(t, 1, tm.Id)
	assert.Equal(t, int64(2), tm.UserId)
	assert.Equal(t, 18, tm.Age)
	assert.Equal(t, "kovey", tm.Name)
	err = QueryRaw(1, context.Background(), db.Raw("SELECT `id`, `user_id`, `age`, `name` FROM `user_1` WHERE `user_id` = ?", 1), &tm1s)
	if err != nil {
		t.Fatal(err)
	}
	tm1 := tm1s[0]
	assert.Nil(t, err)
	assert.Equal(t, 1, tm1.Id)
	assert.Equal(t, int64(1), tm1.UserId)
	assert.Equal(t, 19, tm1.Age)
	assert.Equal(t, "kovey1", tm1.Name)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestDbQueryRow(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()

	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)
	tm := newTestModel()
	tm1 := newTestModel()
	mock1.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_0` WHERE `user_id` = ?").ExpectQuery().WithArgs(2).WillReturnRows(sqlmock.NewRows(tm.Columns()).AddRow(1, 2, 18, "kovey"))
	mock2.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_1` WHERE `user_id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(tm1.Columns()).AddRow(1, 1, 19, "kovey1"))
	in1 := db.NewQuery()
	in1.Table("user_0").Where("user_id", "=", 2).Columns(tm.Columns()...)
	tm.WithKey(2)
	err = QueryRow(2, context.Background(), in1, tm)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, tm.Id)
	assert.Equal(t, int64(2), tm.UserId)
	assert.Equal(t, 18, tm.Age)
	assert.Equal(t, "kovey", tm.Name)
	tm1.WithKey(1)
	err = QueryRowRaw(1, context.Background(), db.Raw("SELECT `id`, `user_id`, `age`, `name` FROM `user_1` WHERE `user_id` = ?", 1), tm1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, tm1.Id)
	assert.Equal(t, int64(1), tm1.UserId)
	assert.Equal(t, 19, tm1.Age)
	assert.Equal(t, "kovey1", tm1.Name)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestDbFind(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()

	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)
	tm := newTestModel()
	tm1 := newTestModel()
	mock1.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_0` WHERE `id` = ?").ExpectQuery().WithArgs(2).WillReturnRows(sqlmock.NewRows(tm.Columns()).AddRow(1, 2, 18, "kovey"))
	mock2.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_1` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(tm1.Columns()).AddRow(1, 1, 19, "kovey1"))
	tm.WithKey(2)
	err = Find(context.Background(), tm, 2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, tm.Id)
	assert.Equal(t, int64(2), tm.UserId)
	assert.Equal(t, 18, tm.Age)
	assert.Equal(t, "kovey", tm.Name)
	tm1.WithKey(1)
	err = Find(context.Background(), tm1, 1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, tm1.Id)
	assert.Equal(t, int64(1), tm1.UserId)
	assert.Equal(t, 19, tm1.Age)
	assert.Equal(t, "kovey1", tm1.Name)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestDbLock(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()

	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)
	tm := newTestModel()
	tm1 := newTestModel()
	mock1.ExpectBegin()
	mock1.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_0` WHERE `id` = ? FOR UPDATE").ExpectQuery().WithArgs(2).WillReturnRows(sqlmock.NewRows(tm.Columns()).AddRow(1, 2, 18, "kovey"))
	mock1.ExpectCommit()
	mock2.ExpectBegin()
	mock2.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_1` WHERE `user_id` = ? FOR UPDATE").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(tm1.Columns()).AddRow(1, 1, 19, "kovey1"))
	mock2.ExpectCommit()
	err = Transaction(context.Background(), []any{1, 2}, func(ctx context.Context, conn ConnectionInterface) error {
		tm.WithKey(2)
		err = Lock(context.Background(), conn, tm, 2)
		if err != nil {
			return err
		}

		tm1.WithKey(1)
		return LockBy(context.Background(), conn, tm1, func(query ksql.QueryInterface) {
			query.Where("user_id", "=", 1)
		})
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, tm.Id)
	assert.Equal(t, int64(2), tm.UserId)
	assert.Equal(t, 18, tm.Age)
	assert.Equal(t, "kovey", tm.Name)
	assert.Equal(t, 1, tm1.Id)
	assert.Equal(t, int64(1), tm1.UserId)
	assert.Equal(t, 19, tm1.Age)
	assert.Equal(t, "kovey1", tm1.Name)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}
