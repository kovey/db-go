package sharding

import (
	"github.com/kovey/db-go/v2/db"
	"github.com/kovey/debug-go/debug"
)

type Tx struct {
	txs map[int]*db.Tx
}

func NewTx() *Tx {
	return &Tx{txs: make(map[int]*db.Tx)}
}

func (t *Tx) Add(id int, tx *db.Tx) {
	t.txs[id] = tx
}

func (t *Tx) Commit() {
	for _, tx := range t.txs {
		if tx.IsCompleted() {
			continue
		}
		if err := tx.Commit(); err != nil {
			debug.Erro(err.Error())
		}
	}
}

func (t *Tx) Rollback() {
	for _, tx := range t.txs {
		if tx.IsCompleted() {
			continue
		}

		if err := tx.Rollback(); err != nil {
			debug.Erro(err.Error())
		}
	}
}

func (t *Tx) IsCompleted() bool {
	for _, tx := range t.txs {
		if !tx.IsCompleted() {
			return false
		}
	}

	return true
}
