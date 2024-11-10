package db

import "fmt"

type SqlErr struct {
	Sql   string
	Binds []any
	Err   error
}

func (s *SqlErr) Error() string {
	if s.Err == nil {
		return ""
	}

	return fmt.Sprintf("sql: %s, binds: %v, error: %s", s.Sql, s.Binds, s.Err)
}

type TxErr struct {
	commitErr   error
	rollbackErr error
	beginErr    error
	callErr     error
}

func (t *TxErr) Commit() error {
	return t.commitErr
}

func (t *TxErr) Begin() error {
	return t.beginErr
}

func (t *TxErr) Rollback() error {
	return t.rollbackErr
}

func (t *TxErr) Call() error {
	return t.callErr
}

func (t *TxErr) Error() string {
	return fmt.Sprintf(
		"begin error: %s, call err: %s, commit error: %s, rollback error: %s",
		t.beginErr, t.callErr, t.commitErr, t.rollbackErr,
	)
}
