package sharding

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type ConnectionInterface interface {
	Get(key any) ksql.ConnectionInterface
	Exec(key any, ctx context.Context, op ksql.SqlInterface) (int64, error)
	QueryRow(key any, ctx context.Context, op ksql.QueryInterface, model ksql.RowInterface) error
	Insert(key any, ctx context.Context, op ksql.InsertInterface) (int64, error)
	Update(key any, ctx context.Context, op ksql.UpdateInterface) (int64, error)
	Delete(key any, ctx context.Context, op ksql.DeleteInterface) (int64, error)
	Database(key any) *sql.DB
	Prepare(key any, ctx context.Context, op ksql.SqlInterface) (*sql.Stmt, error)
	ExecRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (sql.Result, error)
	PrepareRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (*sql.Stmt, error)
	QueryRowRaw(key any, ctx context.Context, raw ksql.ExpressInterface, model ksql.RowInterface) error
	DriverName() string
	InTransaction() bool
	Clone() ConnectionInterface
	Begin(ctx context.Context, options *sql.TxOptions) ksql.TxError
	Rollback(ctx context.Context) ksql.TxError
	Commit(ctx context.Context) ksql.TxError
	Transaction(ctx context.Context, keys []any, call func(ctx context.Context, conn ConnectionInterface) error) ksql.TxError
	TransactionBy(ctx context.Context, keys []any, options *sql.TxOptions, call func(ctx context.Context, conn ConnectionInterface) error) ksql.TxError
	ScanRaw(key any, ctx context.Context, raw ksql.ExpressInterface, data ...any) error
}

type ShardingInterface interface {
	WithKey(key any)
	Key() any
}

type ModelInterface interface {
	ksql.ModelInterface
	ShardingInterface
}

type TxErr struct {
	begins    []string
	calls     []string
	rollbacks []string
	commits   []string
}

func newTxErr() *TxErr {
	return &TxErr{}
}

func (t *TxErr) AppendBegin(key any, begin error) *TxErr {
	t.begins = append(t.begins, fmt.Sprintf("%s on key %v", begin.Error(), key))
	return t
}

func (t *TxErr) AppendCall(call error) *TxErr {
	t.calls = append(t.calls, call.Error())
	return t
}

func (t *TxErr) AppendRollback(key any, rollback error) *TxErr {
	t.rollbacks = append(t.rollbacks, fmt.Sprintf("%s on key %v", rollback.Error(), key))
	return t
}

func (t *TxErr) AppendCommit(key any, commit error) *TxErr {
	t.commits = append(t.commits, fmt.Sprintf("%s on key %v", commit.Error(), key))
	return t
}

func (t *TxErr) Begin() error {
	return fmt.Errorf(strings.Join(t.begins, ";"))
}

func (t *TxErr) Call() error {
	return fmt.Errorf(strings.Join(t.calls, ";"))
}

func (t *TxErr) Rollback() error {
	return fmt.Errorf(strings.Join(t.rollbacks, ";"))
}

func (t *TxErr) Commit() error {
	return fmt.Errorf(strings.Join(t.rollbacks, ";"))
}

func (t *TxErr) Error() string {
	return fmt.Sprintf("begin: %s, calls: %s, rollback: %s, commit: %s", t.Begin(), t.Call(), t.Rollback(), t.Commit())
}
