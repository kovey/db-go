package db

import (
	"database/sql"
	"fmt"
)

type Tx struct {
	tx *sql.Tx
}

func NewTx(tx *sql.Tx) *Tx {
	return &Tx{tx: tx}
}

func (t *Tx) Tx() *sql.Tx {
	return t.tx
}

func (t *Tx) Commit() error {
	if t.IsCompleted() {
		return fmt.Errorf("transaction is IsCompleted")
	}

	err := t.tx.Commit()
	t.tx = nil
	return err
}

func (t *Tx) Rollback() error {
	if t.IsCompleted() {
		return fmt.Errorf("transaction is IsCompleted")
	}

	err := t.tx.Rollback()
	t.tx = nil
	return err
}

func (t *Tx) IsCompleted() bool {
	return t.tx == nil
}
