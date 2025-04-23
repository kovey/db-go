package db

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlError(t *testing.T) {
	raw := Raw("SELECT * FROM user WHERE id = ?", 1)
	err := &SqlErr{Sql: raw.Statement(), Binds: raw.Binds(), Err: errors.New("sql err")}
	assert.Equal(t, "sql: SELECT * FROM user WHERE id = ?, binds: [1], error: sql err", err.Error())
	nErr := &SqlErr{}
	assert.Equal(t, "", nErr.Error())
}

func TestTxError(t *testing.T) {
	raw := Raw("SELECT * FROM user WHERE id = ?", 1)
	err := &SqlErr{Sql: raw.Statement(), Binds: raw.Binds(), Err: errors.New("sql err")}
	beginErr := errors.New("begin error")
	rollbackErr := errors.New("rollback error")
	commitErr := errors.New("commit error")
	txErr := &TxErr{callErr: err, beginErr: beginErr, rollbackErr: rollbackErr, commitErr: commitErr}
	assert.Equal(t, "begin error: begin error, call err: sql: SELECT * FROM user WHERE id = ?, binds: [1], error: sql err, commit error: commit error, rollback error: rollback error", txErr.Error())
	assert.Equal(t, err, txErr.Call())
	assert.Equal(t, beginErr, txErr.Begin())
	assert.Equal(t, rollbackErr, txErr.Rollback())
	assert.Equal(t, commitErr, txErr.Commit())
}
