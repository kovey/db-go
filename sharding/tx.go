package sharding

import (
	"github.com/kovey/db-go/v2/db"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/pool"
	"github.com/kovey/pool/object"
)

const (
	namespace = "ko.db.sharding"
	tx_name   = "Tx"
)

func init() {
	pool.DefaultNoCtx(namespace, tx_name, func() any {
		return &Tx{txs: make(map[int]*db.Tx), ObjNoCtx: object.NewObjNoCtx(namespace, tx_name)}
	})
}

type Tx struct {
	*object.ObjNoCtx
	txs map[int]*db.Tx
}

func NewTx() *Tx {
	return &Tx{txs: make(map[int]*db.Tx)}
}

func NewTxBy(ctx object.CtxInterface) *Tx {
	return ctx.GetNoCtx(namespace, tx_name).(*Tx)
}

func (t *Tx) Reset() {
	t.txs = make(map[int]*db.Tx)
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
