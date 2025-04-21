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

func (c *Connection) BeginTo(ctx context.Context, point string) error {
	if !driver.SupportSavePoint(c.driverName) {
		return Err_Un_Support_Save_Point
	}

	_, err := c.ExecRaw(ctx, ks.Raw("SAVEPOINT ?", point))
	return err
}

func (c *Connection) beginTo(ctx context.Context) error {
	c.transCount++
	if err := c.BeginTo(ctx, fmt.Sprintf("trans_%d", c.transCount)); err != nil {
		c.transCount--
		if err == Err_Un_Support_Save_Point {
			return nil
		}

		return err
	}

	return nil
}

func (c *Connection) RollbackTo(ctx context.Context, point string) error {
	if !driver.SupportSavePoint(c.driverName) {
		return Err_Un_Support_Save_Point
	}

	_, err := c.ExecRaw(ctx, ks.Raw("ROLLBACK TO SAVEPOINT ?", point))
	return err
}

func (c *Connection) rollbackTo(ctx context.Context) error {
	if err := c.RollbackTo(ctx, fmt.Sprintf("trans_%d", c.transCount)); err != nil {
		if err == Err_Un_Support_Save_Point {
			c.transCount--
		}
		return err
	}

	c.transCount--
	return nil
}

func (c *Connection) CommitTo(ctx context.Context, point string) error {
	if !driver.SupportSavePoint(c.driverName) {
		return Err_Un_Support_Save_Point
	}

	_, err := c.ExecRaw(ctx, ks.Raw("RELEASE SAVEPOINT ?", point))
	return err
}

func (c *Connection) commitTo(ctx context.Context) error {
	if err := c.CommitTo(ctx, fmt.Sprintf("trans_%d", c.transCount)); err != nil {
		if err == Err_Un_Support_Save_Point {
			c.transCount--
		}
		return err
	}

	c.transCount--
	return nil
}

func (c *Connection) Begin(ctx context.Context, options *sql.TxOptions) error {
	if c.tx != nil {
		return c.beginTo(ctx)
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
		return c.rollbackTo(ctx)
	}

	defer c.reset()
	return c.tx.Rollback()
}

func (c *Connection) Commit(ctx context.Context) error {
	if c.tx == nil {
		return Err_Not_In_Transaction
	}

	if c.transCount > 0 {
		return c.commitTo(ctx)
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
	cc := NewContext(ctx)
	cc.SqlLogStart(op)
	defer cc.SqlLogEnd()

	stmt, err := c.Prepare(cc, op)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(cc, op.Binds()...)
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
	cc := NewContext(ctx)
	cc.SqlLogStart(op)
	defer cc.SqlLogEnd()

	stmt, err := c.Prepare(cc, op)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(cc, op.Binds()...)
	if row.Err() != nil {
		return _err(err, op)
	}

	if err := model.Scan(row, model); err != nil {
		if err == sql.ErrNoRows {
			model.WithConn(c)
			return nil
		}

		return _err(err, op)
	}

	model.Sharding(op.GetSharding())
	model.WithConn(c)
	return nil
}

func (c *Connection) QueryRowRaw(ctx context.Context, raw ksql.ExpressInterface, model ksql.RowInterface) error {
	if raw.IsExec() {
		return Err_Sql_Not_Query
	}

	cc := NewContext(ctx)
	cc.RawSqlLogStart(raw)
	defer cc.SqlLogEnd()

	stmt, err := c.PrepareRaw(cc, raw)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(cc, raw.Binds()...)
	if row.Err() != nil {
		return _errRaw(err, raw)
	}

	if err := model.Scan(row, model); err != nil {
		if err == sql.ErrNoRows {
			model.WithConn(c)
			return nil
		}

		return _errRaw(err, raw)
	}

	model.WithConn(c)
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
	if !raw.IsExec() {
		return nil, Err_Sql_Not_Exec
	}

	cc := NewContext(ctx)
	cc.RawSqlLogStart(raw)
	defer cc.SqlLogEnd()

	stmt, err := c.PrepareRaw(cc, raw)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(cc, raw.Binds()...)
	return result, _errRaw(err, raw)
}

func (c *Connection) InTransaction() bool {
	return c.tx != nil
}

func (c *Connection) ScanRaw(ctx context.Context, raw ksql.ExpressInterface, data ...any) error {
	if raw.IsExec() {
		return Err_Sql_Not_Query
	}

	cc := NewContext(ctx)
	cc.RawSqlLogStart(raw)
	defer cc.SqlLogEnd()

	stmt, err := c.PrepareRaw(cc, raw)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(cc, raw.Binds()...)
	if row.Err() != nil {
		return _errRaw(err, raw)
	}

	if err := row.Scan(data...); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return _errRaw(err, raw)
	}

	return nil
}

func (c *Connection) Scan(ctx context.Context, query ksql.QueryInterface, data ...any) error {
	cc := NewContext(ctx)
	cc.SqlLogStart(query)
	defer cc.SqlLogEnd()

	stmt, err := c.Prepare(cc, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(cc, query.Binds()...)
	if row.Err() != nil {
		return _err(err, query)
	}

	if err := row.Scan(data...); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return _err(err, query)
	}

	return nil
}
