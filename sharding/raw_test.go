package sharding

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kovey/db-go/v3/db"
	"github.com/stretchr/testify/assert"
)

func TestRawInsert(t *testing.T) {
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
	id1, err := InsertRaw(2, context.Background(), db.Raw("INSERT INTO `user_0` (`user_id`, `name`, `age`) VALUES (?, ?, ?)", 2, "kovey", 18))
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id1)
	id2, err := InsertRaw(1, context.Background(), db.Raw("INSERT INTO `user_1` (`user_id`, `name`, `age`) VALUES (?, ?, ?)", 1, "kovey1", 19))
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestRawUpdate(t *testing.T) {
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
	id1, err := UpdateRaw(2, context.Background(), db.Raw("UPDATE `user_0` SET `name` = ?, `age` = ? WHERE `user_id` = ?", "kovey", 18, 2))
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id1)
	id2, err := UpdateRaw(1, context.Background(), db.Raw("UPDATE `user_1` SET `name` = ?, `age` = ? WHERE `user_id` = ?", "kovey1", 19, 1))
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestRawDelete(t *testing.T) {
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
	id1, err := DeleteRaw(2, context.Background(), db.Raw("DELETE FROM `user_0` WHERE `user_id` = ?", 2))
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id1)
	id2, err := DeleteRaw(1, context.Background(), db.Raw("DELETE FROM `user_1` WHERE `user_id` = ?", 1))
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id2)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}
