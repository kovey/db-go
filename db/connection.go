package db

import (
	"context"
	"database/sql"

	"github.com/kovey/db-go/v3"
)

type Connection struct {
	Tx         *sql.Tx
	database   *sql.DB
	driverName string
}

func (c *Connection) DriverName() string {
	return c.driverName
}

func (c *Connection) Database() *sql.DB {
	return c.database
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
	if c.Tx != nil {
		stmt, err := c.Tx.PrepareContext(ctx, op.Prepare())
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
	if c.Tx != nil {
		stmt, err := c.Tx.PrepareContext(ctx, raw.Statement())
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

	result, err := stmt.ExecContext(ctx, raw.Binds()...)
	return result, _errRaw(err, raw)
}
