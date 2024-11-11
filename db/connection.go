package db

import (
	"context"
	"database/sql"
	"fmt"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db/driver"
	ks "github.com/kovey/db-go/v3/sql"
)

type Connection struct {
	tx         *sql.Tx
	database   *sql.DB
	driverName string
	transCount int
}

func (c *Connection) DriverName() string {
	return c.driverName
}

func (c *Connection) Database() *sql.DB {
	return c.database
}

func (c *Connection) Clone() ksql.ConnectionInterface {
	return &Connection{database: c.database, driverName: c.driverName, tx: nil}
}

func (c *Connection) BeginTo(ctx context.Context) error {
	if c.tx == nil {
		return Err_Not_In_Transaction
	}

	if !driver.SupportSavePoint(c.driverName) {
		c.transCount++
		return nil
	}

	c.transCount++
	_, err := c.ExecRaw(ctx, ks.Raw("SAVEPOINT ?", fmt.Sprintf("trans_%d", c.transCount)))
	if err != nil {
		c.transCount--
	}

	return err
}

func (c *Connection) RollbackTo(ctx context.Context) error {
	if c.tx == nil || c.transCount == 0 {
		return Err_Not_In_Transaction
	}

	if !driver.SupportSavePoint(c.driverName) {
		c.transCount--
		return nil
	}

	_, err := c.ExecRaw(ctx, ks.Raw("ROLLBACK SAVEPOINT ?", fmt.Sprintf("trans_%d", c.transCount)))
	if err == nil {
		c.transCount--
	}
	return err
}

func (c *Connection) CommitTo(ctx context.Context) error {
	if c.tx == nil || c.transCount == 0 {
		return Err_Not_In_Transaction
	}

	if !driver.SupportSavePoint(c.driverName) {
		c.transCount--
		return nil
	}

	_, err := c.ExecRaw(ctx, ks.Raw("RELEASE SAVEPOINT ?", fmt.Sprintf("trans_%d", c.transCount)))
	if err == nil {
		c.transCount--
	}
	return err
}

func (c *Connection) Begin(ctx context.Context, options *sql.TxOptions) error {
	if c.tx != nil {
		if err := c.BeginTo(ctx); err != nil {
			return err
		}

		return nil
	}

	tx, err := c.database.BeginTx(ctx, options)
	if err != nil {
		return err
	}

	c.tx = tx
	return nil
}

func (c *Connection) Rollback(ctx context.Context) error {
	if c.tx == nil {
		return Err_Not_In_Transaction
	}

	if c.transCount > 0 {
		return c.RollbackTo(ctx)
	}

	defer c.reset()
	return c.tx.Rollback()
}

func (c *Connection) Commit(ctx context.Context) error {
	if c.tx == nil {
		return Err_Not_In_Transaction
	}

	if c.transCount > 0 {
		return c.CommitTo(ctx)
	}

	defer c.reset()
	return c.tx.Commit()
}

func (c *Connection) reset() {
	c.tx = nil
}

func (c *Connection) Transaction(ctx context.Context, call func(ctx context.Context, conn ksql.ConnectionInterface) error) ksql.TxError {
	return c.TransactionBy(ctx, nil, call)
}

func (c *Connection) TransactionBy(ctx context.Context, options *sql.TxOptions, call func(ctx context.Context, conn ksql.ConnectionInterface) error) ksql.TxError {
	if err := c.Begin(ctx, options); err != nil {
		return &TxErr{beginErr: err}
	}

	callErr := call(ctx, c)
	if callErr != nil {
		txErr := &TxErr{callErr: callErr}
		if err := c.Rollback(ctx); err != nil {
			txErr.rollbackErr = err
		}

		return txErr
	}

	if err := c.Commit(ctx); err != nil {
		return &TxErr{commitErr: err}
	}

	return nil
}

func (c *Connection) Insert(ctx context.Context, op ksql.InsertInterface) (int64, error) {
	return c.Exec(ctx, op)
}

func (c *Connection) Update(ctx context.Context, op ksql.UpdateInterface) (int64, error) {
	return c.Exec(ctx, op)
}

func (c *Connection) Delete(ctx context.Context, op ksql.DeleteInterface) (int64, error) {
	return c.Exec(ctx, op)
}

func (c *Connection) Prepare(ctx context.Context, op ksql.SqlInterface) (*sql.Stmt, error) {
	if c.tx != nil {
		stmt, err := c.tx.PrepareContext(ctx, op.Prepare())
		return stmt, _err(err, op)
	}

	stmt, err := c.database.PrepareContext(ctx, op.Prepare())
	return stmt, _err(err, op)
}

func _err(err error, op ksql.SqlInterface) error {
	if err != nil {
		return &SqlErr{Sql: op.Prepare(), Binds: op.Binds(), Err: err}
	}

	return err
}

func _errRaw(err error, op ksql.ExpressInterface) error {
	if err != nil {
		return &SqlErr{Sql: op.Statement(), Binds: op.Binds(), Err: err}
	}

	return err
}

func (c *Connection) Exec(ctx context.Context, op ksql.SqlInterface) (int64, error) {
	stmt, err := c.Prepare(ctx, op)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, op.Binds()...)
	if err != nil {
		return 0, _err(err, op)
	}

	switch op.(type) {
	case ksql.InsertInterface:
		id, err := result.LastInsertId()
		return id, _err(err, op)
	default:
		id, err := result.RowsAffected()
		return id, _err(err, op)
	}
}

func (c *Connection) QueryRow(ctx context.Context, op ksql.QueryInterface, model ksql.RowInterface) error {
	stmt, err := c.Prepare(ctx, op)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, op.Binds()...)
	if row.Err() != nil {
		return _err(err, op)
	}

	if err := row.Scan(model.Values()...); err != nil {
		if err == sql.ErrNoRows {
			model.SetConn(c)
			return nil
		}

		return _err(err, op)
	}

	model.FromFetch()
	model.SetConn(c)
	return nil
}

func (c *Connection) QueryRowRaw(ctx context.Context, raw ksql.ExpressInterface, model ksql.RowInterface) error {
	stmt, err := c.PrepareRaw(ctx, raw)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, raw.Binds()...)
	if row.Err() != nil {
		return _errRaw(err, raw)
	}

	if err := row.Scan(model.Values()...); err != nil {
		if err == sql.ErrNoRows {
			model.SetConn(c)
			return nil
		}

		return _errRaw(err, raw)
	}

	model.FromFetch()
	model.SetConn(c)
	return nil
}

func (c *Connection) PrepareRaw(ctx context.Context, raw ksql.ExpressInterface) (*sql.Stmt, error) {
	if c.tx != nil {
		stmt, err := c.tx.PrepareContext(ctx, raw.Statement())
		return stmt, _errRaw(err, raw)
	}

	stmt, err := c.database.PrepareContext(ctx, raw.Statement())
	return stmt, _errRaw(err, raw)
}

func (c *Connection) ExecRaw(ctx context.Context, raw ksql.ExpressInterface) (sql.Result, error) {
	stmt, err := c.PrepareRaw(ctx, raw)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, raw.Binds()...)
	return result, _errRaw(err, raw)
}

func (c *Connection) InTransaction() bool {
	return c.tx != nil
}
