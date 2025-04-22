package sharding

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kovey/db-go/v3/db"
	"github.com/stretchr/testify/assert"
)

func TestConnectionInsert(t *testing.T) {
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
	in1 := db.NewInsert()
	in1.Table("user_0").Add("user_id", 2).Add("name", "kovey").Add("age", 18)
	in2 := db.NewInsert()
	in2.Table("user_1").Add("user_id", 1).Add("name", "kovey1").Add("age", 19)
	id1, err := database.Insert(2, context.Background(), in1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id1)
	id2, err := database.Insert(1, context.Background(), in2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestConnectionUpdate(t *testing.T) {
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
	in1 := db.NewUpdate()
	w1 := db.NewWhere().Where("user_id", "=", 2)
	in1.Table("user_0").Set("name", "kovey").Set("age", 18).Where(w1)
	in2 := db.NewUpdate()
	w2 := db.NewWhere().Where("user_id", "=", 1)
	in2.Table("user_1").Set("name", "kovey1").Set("age", 19).Where(w2)
	id1, err := database.Update(2, context.Background(), in1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id1)
	id2, err := database.Update(1, context.Background(), in2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestConnectionDelete(t *testing.T) {
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
	in1 := db.NewDelete()
	w1 := db.NewWhere().Where("user_id", "=", 2)
	in1.Table("user_0").Where(w1)
	in2 := db.NewDelete()
	w2 := db.NewWhere().Where("user_id", "=", 1)
	in2.Table("user_1").Where(w2)
	id1, err := database.Delete(2, context.Background(), in1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id1)
	id2, err := database.Delete(1, context.Background(), in2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestConnectionQuery(t *testing.T) {
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
	err = database.QueryRow(2, context.Background(), in1, tm)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, tm.Id)
	assert.Equal(t, int64(2), tm.UserId)
	assert.Equal(t, 18, tm.Age)
	assert.Equal(t, "kovey", tm.Name)
	tm1.WithKey(1)
	err = database.QueryRowRaw(1, context.Background(), db.Raw("SELECT `id`, `user_id`, `age`, `name` FROM `user_1` WHERE `user_id` = ?", 1), tm1)
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

func TestConnectionTransaction(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()
	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)

	mock1.ExpectBegin()
	mock1.ExpectPrepare("INSERT INTO `user_0` (`user_id`, `name`, `age`) VALUES (?, ?, ?)").ExpectExec().WithArgs(2, "kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	mock1.ExpectPrepare("UPDATE `user_0` SET `name` = ?, `age` = ? WHERE `user_id` = ?").ExpectExec().WithArgs("kovey", 18, 2).WillReturnResult(sqlmock.NewResult(0, 1))
	mock1.ExpectCommit()
	mock2.ExpectBegin()
	mock2.ExpectPrepare("INSERT INTO `user_1` (`user_id`, `name`, `age`) VALUES (?, ?, ?)").ExpectExec().WithArgs(1, "kovey1", 19).WillReturnResult(sqlmock.NewResult(1, 1))
	mock2.ExpectPrepare("UPDATE `user_1` SET `name` = ?, `age` = ? WHERE `user_id` = ?").ExpectExec().WithArgs("kovey1", 19, 1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock2.ExpectCommit()

	err = database.Clone().Transaction(context.Background(), []any{1, 2}, func(ctx context.Context, conn ConnectionInterface) error {
		in1 := db.NewInsert()
		in1.Table("user_0").Add("user_id", 2).Add("name", "kovey").Add("age", 18)
		in2 := db.NewInsert()
		in2.Table("user_1").Add("user_id", 1).Add("name", "kovey1").Add("age", 19)
		_, err := database.Insert(2, context.Background(), in1)
		if err != nil {
			return err
		}
		assert.Nil(t, err)
		_, err = database.Insert(1, context.Background(), in2)
		up1 := db.NewUpdate()
		w1 := db.NewWhere().Where("user_id", "=", 2)
		up1.Table("user_0").Set("name", "kovey").Set("age", 18).Where(w1)
		up2 := db.NewUpdate()
		w2 := db.NewWhere().Where("user_id", "=", 1)
		up2.Table("user_1").Set("name", "kovey1").Set("age", 19).Where(w2)
		_, err = database.Update(2, context.Background(), up1)
		if err != nil {
			return err
		}
		_, err = database.Update(1, context.Background(), up2)
		return err
	})
	assert.Nil(t, err)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}
