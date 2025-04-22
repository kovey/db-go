package sharding

import (
	"context"
	"database/sql"

	ksql "github.com/kovey/db-go/v3"
)

type baseConnection struct {
	conns      []ksql.ConnectionInterface
	count      int
	driverName string
}

func (b *baseConnection) first() ksql.ConnectionInterface {
	return b.conns[0]
}

func (b *baseConnection) conn(key any) ksql.ConnectionInterface {
	return b.conns[b.node(key)]
}

func (b *baseConnection) node(key any) int {
	return node(key, b.count)
}

func (b *baseConnection) Close() error {
	var err error
	for _, conn := range b.conns {
		err = conn.Database().Close()
	}

	return err
}

func (b *baseConnection) Range(call func(index int, conn ksql.ConnectionInterface) error) error {
	for index, conn := range b.conns {
		if err := call(index, conn); err != nil {
			return err
		}
	}

	return nil
}

type Connection struct {
	*baseConnection
	currents      map[any]ksql.ConnectionInterface
	inTransaction bool
	keys          []any
}

func (c *Connection) Clone() ConnectionInterface {
	return &Connection{baseConnection: c.baseConnection, currents: make(map[any]ksql.ConnectionInterface)}
}

func (c *Connection) _keys(keys []any) {
	for _, key := range keys {
		if _, ok := c.currents[key]; ok {
			continue
		}

		c.currents[key] = c.conn(key).Clone()
		c.keys = append(c.keys, key)
	}
}

func (c *Connection) Get(key any) ksql.ConnectionInterface {
	if conn, ok := c.currents[key]; ok {
		return conn
	}

	c.currents[key] = c.conn(key).Clone()
	c.keys = append(c.keys, key)
	return c.currents[key]
}

func (c *Connection) Exec(key any, ctx context.Context, op ksql.SqlInterface) (int64, error) {
	return c.Get(key).Exec(ctx, op)
}

func (c *Connection) QueryRow(key any, ctx context.Context, op ksql.QueryInterface, model ksql.RowInterface) error {
	return c.Get(key).QueryRow(ctx, op, model)
}

func (c *Connection) Insert(key any, ctx context.Context, op ksql.InsertInterface) (int64, error) {
	return c.Get(key).Insert(ctx, op)
}

func (c *Connection) Update(key any, ctx context.Context, op ksql.UpdateInterface) (int64, error) {
	return c.Get(key).Update(ctx, op)
}

func (c *Connection) Delete(key any, ctx context.Context, op ksql.DeleteInterface) (int64, error) {
	return c.Get(key).Delete(ctx, op)
}

func (c *Connection) Database(key any) *sql.DB {
	conn := c.Get(key)
	if conn == nil {
		return nil
	}

	return conn.Database()
}

func (c *Connection) Prepare(key any, ctx context.Context, op ksql.SqlInterface) (*sql.Stmt, error) {
	return c.Get(key).Prepare(ctx, op)
}

func (c *Connection) ExecRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (sql.Result, error) {
	return c.Get(key).ExecRaw(ctx, raw)
}

func (c *Connection) PrepareRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (*sql.Stmt, error) {
	return c.Get(key).PrepareRaw(ctx, raw)
}

func (c *Connection) QueryRowRaw(key any, ctx context.Context, raw ksql.ExpressInterface, model ksql.RowInterface) error {
	return c.Get(key).QueryRowRaw(ctx, raw, model)
}

func (c *Connection) DriverName() string {
	return c.driverName
}

func (c *Connection) InTransaction() bool {
	return c.inTransaction
}

func (c *Connection) _rollback(ctx context.Context, i int) *TxErr {
	txErr := newTxErr()
	for i >= 0 {
		if err := c.Get(c.keys[i]).Rollback(ctx); err != nil {
			txErr.AppendRollback(c.keys[i], err)
		}
		i--
	}

	return txErr
}

func (c *Connection) Begin(ctx context.Context, options *sql.TxOptions) ksql.TxError {
	for i := 0; i < len(c.keys); i++ {
		if err := c.currents[c.keys[i]].Begin(ctx, options); err != nil {
			return c._rollback(ctx, i)
		}
	}

	c.inTransaction = true
	return nil
}

func (c *Connection) Rollback(ctx context.Context) ksql.TxError {
	var txErr *TxErr
	for _, key := range c.keys {
		if err := c.currents[key].Rollback(ctx); err != nil {
			if txErr == nil {
				txErr = newTxErr()
			}

			txErr.AppendRollback(key, err)
		}
	}

	c.inTransaction = false
	return txErr
}

func (c *Connection) Commit(ctx context.Context) ksql.TxError {
	var txErr *TxErr
	for _, key := range c.keys {
		if err := c.currents[key].Commit(ctx); err != nil {
			if txErr == nil {
				txErr = newTxErr()
			}

			txErr.AppendCommit(key, err)
		}
	}

	c.inTransaction = false
	return txErr
}

func (c *Connection) Transaction(ctx context.Context, keys []any, call func(ctx context.Context, conn ConnectionInterface) error) ksql.TxError {
	return c.TransactionBy(ctx, keys, nil, call)
}

func (c *Connection) TransactionBy(ctx context.Context, keys []any, options *sql.TxOptions, call func(ctx context.Context, conn ConnectionInterface) error) ksql.TxError {
	c._keys(keys)
	if err := c.Begin(ctx, options); err != nil {
		return err
	}

	if err := call(ctx, c); err != nil {
		if txErr := c.Rollback(ctx); txErr != nil {
			tmp := txErr.(*TxErr)
			return tmp.AppendCall(err)
		}

		return newTxErr().AppendCall(err)
	}

	return c.Commit(ctx)
}

func (c *Connection) ScanRaw(key any, ctx context.Context, raw ksql.ExpressInterface, data ...any) error {
	return c.Get(key).ScanRaw(ctx, raw, data...)
}
