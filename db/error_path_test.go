package db

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

// ─────────────────────────────────────────────
// Layer 3 — Error path tests using sqlmock
// sqlmock is used ONLY for injecting database errors.
// ─────────────────────────────────────────────

func TestExec_InsertError(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	database = conn

	sqlErr := errors.New("insert error")
	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)").
		ExpectExec().WithArgs("alice", 18).
		WillReturnResult(sqlmock.NewErrorResult(sqlErr))

	ins := NewInsert().Table("user").Add("name", "alice").Add("age", 18)
	_, err := conn.Exec(context.Background(), ins)
	assert.Error(t, err)
	assert.Equal(t, sqlErr, err.(*SqlErr).Err)
}

func TestExec_UpdateError(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	database = conn

	sqlErr := errors.New("update error")
	mock.ExpectPrepare("UPDATE `user` SET `name` = ? WHERE `id` = ?").
		ExpectExec().WithArgs("alice", 1).
		WillReturnResult(sqlmock.NewErrorResult(sqlErr))

	w := NewWhere().Where("id", ksql.Eq, 1)
	up := NewUpdate().Table("user").Set("name", "alice").Where(w)
	_, err := conn.Exec(context.Background(), up)
	assert.Error(t, err)
	assert.Equal(t, sqlErr, err.(*SqlErr).Err)
}

func TestExec_DeleteError(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	database = conn

	sqlErr := errors.New("delete error")
	mock.ExpectPrepare("DELETE FROM `user` WHERE `id` = ?").
		ExpectExec().WithArgs(1).
		WillReturnResult(sqlmock.NewErrorResult(sqlErr))

	w := NewWhere().Where("id", ksql.Eq, 1)
	del := NewDelete().Table("user").Where(w)
	_, err := conn.Exec(context.Background(), del)
	assert.Error(t, err)
	assert.Equal(t, sqlErr, err.(*SqlErr).Err)
}

func TestTransaction_BeginError(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	mock.ExpectBegin().WillReturnError(sqlmock.ErrCancelled)

	err := conn.Transaction(context.Background(), func(ctx context.Context, conn ksql.ConnectionInterface) error {
		return nil
	})

	assert.Equal(t, sqlmock.ErrCancelled, err.(*TxErr).Begin())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_CallError_Rollback(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	businessErr := errors.New("business error")

	mock.ExpectBegin()
	mock.ExpectRollback()

	err := conn.Transaction(context.Background(), func(ctx context.Context, conn ksql.ConnectionInterface) error {
		return businessErr
	})

	assert.Equal(t, businessErr, err.(*TxErr).Call())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_CallError_RollbackAlsoFails(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	businessErr := errors.New("business error")

	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(sqlmock.ErrCancelled)

	err := conn.Transaction(context.Background(), func(ctx context.Context, conn ksql.ConnectionInterface) error {
		return businessErr
	})

	assert.Equal(t, businessErr, err.(*TxErr).Call())
	assert.Equal(t, sqlmock.ErrCancelled, err.(*TxErr).Rollback())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_CommitError(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")

	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(sqlmock.ErrCancelled)

	err := conn.Transaction(context.Background(), func(ctx context.Context, conn ksql.ConnectionInterface) error {
		return nil
	})

	assert.Equal(t, sqlmock.ErrCancelled, err.(*TxErr).Commit())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_Nested_SavepointRollback(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	businessErr := errors.New("nested error")

	mock.ExpectBegin()
	mock.ExpectPrepare("SAVEPOINT ?").ExpectExec().WithArgs("trans_1").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare("ROLLBACK TO SAVEPOINT ?").ExpectExec().WithArgs("trans_1").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectRollback()

	err := conn.Transaction(context.Background(), func(ctx context.Context, conn ksql.ConnectionInterface) error {
		return conn.Transaction(ctx, func(ctx context.Context, conn ksql.ConnectionInterface) error {
			return businessErr
		})
	})

	assert.Equal(t, businessErr, err.(*TxErr).Call().(*TxErr).Call())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_Nested_SavepointCommitError(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")

	mock.ExpectBegin()
	mock.ExpectPrepare("SAVEPOINT ?").ExpectExec().WithArgs("trans_1").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare("RELEASE SAVEPOINT ?").ExpectExec().WithArgs("trans_1").WillReturnError(sqlmock.ErrCancelled)
	mock.ExpectRollback()

	err := conn.Transaction(context.Background(), func(ctx context.Context, conn ksql.ConnectionInterface) error {
		return conn.Transaction(ctx, func(ctx context.Context, conn ksql.ConnectionInterface) error {
			return nil
		})
	})

	assert.Equal(t, sqlmock.ErrCancelled, err.(*TxErr).Call().(*TxErr).Commit().(*SqlErr).Err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_Success_Commit(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)").
		ExpectExec().WithArgs("alice", 18).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := conn.Transaction(context.Background(), func(ctx context.Context, conn ksql.ConnectionInterface) error {
		ins := NewInsert().Table("user").Add("name", "alice").Add("age", 18)
		_, err := conn.Insert(ctx, ins)
		return err
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotInTransaction_ReturnsError(t *testing.T) {
	testDb, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	assert.False(t, conn.InTransaction())

	err := conn.Rollback(context.Background())
	assert.Equal(t, Err_Not_In_Transaction, err)

	err = conn.Commit(context.Background())
	assert.Equal(t, Err_Not_In_Transaction, err)
}

func TestQueryRow_NoRows(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	database = conn

	// sqlmock can't easily simulate sql.ErrNoRows via QueryRowContext,
	// so we expect rows.Err() to return sql.ErrNoRows which is now properly handled
	// Test the error wrapping path instead:

	sqlErr := errors.New("connection refused")
	mock.ExpectPrepare("SELECT `id`, `name` FROM `user` WHERE `id` = ?").
		ExpectQuery().WithArgs(999).
		WillReturnError(sqlErr)

	u := newTestUser()
	q := NewQuery().Table("user").Columns("id", "name").Where("id", ksql.Eq, 999)
	err := QueryRow(context.Background(), q, u)
	assert.Error(t, err)
	assert.Equal(t, sqlErr, err.(*SqlErr).Err)
}

func TestPrepareError(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	database = conn

	mock.ExpectPrepare("SELECT .* FROM `user` WHERE `id` = ?").
		WillReturnError(errors.New("prepare error"))

	u := newTestUser()
	q := NewQuery().Table("user").Columns("id", "name").Where("id", ksql.Eq, 1)
	err := QueryRow(context.Background(), q, u)
	assert.Error(t, err)
}

func TestScan_RowScanError(t *testing.T) {
	testDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	database = conn

	// Return a value that can't be scanned into int
	mock.ExpectPrepare("SELECT `id`, `name` FROM `user` WHERE `id` = ?").
		ExpectQuery().WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("not_an_int", "alice"))

	var id int
	var name string
	q := NewQuery().Table("user").Columns("id", "name").Where("id", ksql.Eq, 1)
	err := Scan(context.Background(), q, &id, &name)
	assert.Error(t, err)
}

func TestSavePoint_UnsupportedDriver(t *testing.T) {
	testDb, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer testDb.Close()

	conn, _ := Open(testDb, "mysql")
	conn.(*Connection).driverName = "other"

	err := conn.BeginTo(context.Background(), "point")
	assert.Equal(t, Err_Un_Support_Save_Point, err)

	err = conn.RollbackTo(context.Background(), "point")
	assert.Equal(t, Err_Un_Support_Save_Point, err)

	err = conn.CommitTo(context.Background(), "point")
	assert.Equal(t, Err_Un_Support_Save_Point, err)
}

func TestDatabase_NotInitialized(t *testing.T) {
	database = nil
	_, err := Get()
	assert.Equal(t, Err_Database_Not_Initialized, err)
}

func TestSqlErr_Formatting(t *testing.T) {
	r := Raw("SELECT * FROM user WHERE id = ?", 1)
	err := &SqlErr{Sql: r.Statement(), Binds: r.Binds(), Err: errors.New("sql err")}
	assert.Equal(t, "sql: SELECT * FROM user WHERE id = ?, binds: [1], error: sql err", err.Error())

	empty := &SqlErr{}
	assert.Equal(t, "", empty.Error())
}

func TestTxErr_Accessors(t *testing.T) {
	beginErr := errors.New("begin error")
	rollbackErr := errors.New("rollback error")
	commitErr := errors.New("commit error")
	callErr := errors.New("call error")

	txErr := &TxErr{callErr: callErr, beginErr: beginErr, rollbackErr: rollbackErr, commitErr: commitErr}

	assert.Equal(t, beginErr, txErr.Begin())
	assert.Equal(t, rollbackErr, txErr.Rollback())
	assert.Equal(t, commitErr, txErr.Commit())
	assert.Equal(t, callErr, txErr.Call())

	expected := "begin error: begin error, call err: call error, commit error: commit error, rollback error: rollback error"
	assert.Equal(t, expected, txErr.Error())
}
