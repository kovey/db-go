package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	ksql "github.com/kovey/db-go/v3"
)

var Err_Un_Support_Operate = errors.New("unsupport operate")
var Err_Not_In_Transaction = errors.New("not in transaction")
var Err_Database_Not_Initialized = errors.New("data not initialized")
var Err_Un_Support_Save_Point = errors.New("unsupport save point")

var database ksql.ConnectionInterface

type Config struct {
	DriverName     string
	DataSourceName string
	MaxIdleTime    time.Duration
	MaxLifeTime    time.Duration
	MaxIdleConns   int
	MaxOpenConns   int
}

func Database() *sql.DB {
	return database.Database()
}

func Open(conn *sql.DB, driverName string) (ksql.ConnectionInterface, error) {
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Connection{database: conn, driverName: driverName}, nil
}

func Get() (ksql.ConnectionInterface, error) {
	if database == nil {
		return nil, Err_Database_Not_Initialized
	}

	return database.Clone(), nil
}

// init global connection
func Init(conf Config) error {
	db, err := sql.Open(conf.DriverName, conf.DataSourceName)
	if err != nil {
		return err
	}

	db.SetConnMaxIdleTime(conf.MaxIdleTime)
	db.SetConnMaxLifetime(conf.MaxLifeTime)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxIdleConns)
	conn, err := Open(db, conf.DriverName)
	if err != nil {
		return err
	}

	database = conn
	return nil
}

func InsertBy(ctx context.Context, conn ksql.ConnectionInterface, table string, data *Data) (int64, error) {
	op := NewInsert()
	op.Table(table)
	data.Range(func(key string, val any) {
		op.Add(key, val)
	})

	return conn.Insert(ctx, op)
}

func Insert(ctx context.Context, table string, data *Data) (int64, error) {
	return InsertBy(ctx, database, table, data)
}

func UpdateBy(ctx context.Context, conn ksql.ConnectionInterface, table string, data *Data, where ksql.WhereInterface) (int64, error) {
	op := NewUpdate()
	op.Table(table)
	data.Range(func(key string, val any) {
		op.Set(key, val)
	})
	op.Where(where)

	return conn.Exec(ctx, op)
}

func Update(ctx context.Context, table string, data *Data, where ksql.WhereInterface) (int64, error) {
	return UpdateBy(ctx, database, table, data, where)
}

func DeleteBy(ctx context.Context, conn ksql.ConnectionInterface, table string, where ksql.WhereInterface) (int64, error) {
	op := NewDelete()
	op.Table(table).Where(where)

	return conn.Exec(ctx, op)
}

func Delete(ctx context.Context, table string, where ksql.WhereInterface) (int64, error) {
	return DeleteBy(ctx, database, table, where)
}

func ExecBy(ctx context.Context, conn ksql.ConnectionInterface, op ksql.SqlInterface) (int64, error) {
	return conn.Exec(ctx, op)
}

func Exec(ctx context.Context, op ksql.SqlInterface) (int64, error) {
	return ExecBy(ctx, database, op)
}

func QueryBy[T ksql.RowInterface](ctx context.Context, conn ksql.ConnectionInterface, op ksql.QueryInterface, models *[]T) error {
	stmt, err := conn.Prepare(ctx, op)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, op.Binds()...)
	if err != nil {
		return _err(err, op)
	}
	defer rows.Close()
	if rows.Err() != nil {
		return _err(rows.Err(), op)
	}

	var m T
	for rows.Next() {
		tmp := m.Clone()
		if err := rows.Scan(tmp.Values()...); err != nil {
			return _err(err, op)
		}

		model, ok := tmp.(T)
		if !ok {
			continue
		}

		model.SetConn(conn)
		model.FromFetch()
		*models = append(*models, model)
	}

	return nil
}

func Query[T ksql.RowInterface](ctx context.Context, op ksql.QueryInterface, models *[]T) error {
	return QueryBy(ctx, database, op, models)
}

func QueryRowBy[T ksql.RowInterface](ctx context.Context, conn ksql.ConnectionInterface, op ksql.QueryInterface, model T) error {
	stmt, err := conn.Prepare(ctx, op)
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
			return nil
		}

		return _err(err, op)
	}

	model.FromFetch()
	model.SetConn(conn)
	return nil
}

func QueryRow[T ksql.RowInterface](ctx context.Context, op ksql.QueryInterface, model T) error {
	return QueryRowBy(ctx, database, op, model)
}

func TransactionBy(ctx context.Context, options *sql.TxOptions, call func(ctx context.Context, db ksql.ConnectionInterface) error) ksql.TxError {
	conn := database.Clone()
	if err := conn.Begin(ctx, options); err != nil {
		return &TxErr{beginErr: err}
	}

	callErr := call(ctx, conn)
	if callErr != nil {
		txErr := &TxErr{callErr: callErr}
		if err := conn.Rollback(ctx); err != nil {
			txErr.rollbackErr = err
		}

		return txErr
	}

	if err := conn.Commit(ctx); err != nil {
		return &TxErr{commitErr: err}
	}

	return nil
}

func Transaction(ctx context.Context, call func(ctx context.Context, db ksql.ConnectionInterface) error) ksql.TxError {
	return TransactionBy(ctx, nil, call)
}

func Find[T FindType](ctx context.Context, model ksql.ModelInterface, id T) error {
	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).Where(model.PrimaryId(), "=", id)
	return QueryRow(ctx, query, model)
}

func FindBy(ctx context.Context, model ksql.ModelInterface, call func(query ksql.QueryInterface)) error {
	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...)
	call(query)
	return QueryRow(ctx, query, model)
}

func Lock[T FindType](ctx context.Context, conn ksql.ConnectionInterface, model ksql.ModelInterface, id T) error {
	if conn == nil || !conn.InTransaction() {
		return Err_Not_In_Transaction
	}

	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).Where(model.PrimaryId(), "=", id).ForUpdate()
	return QueryRowBy(ctx, conn, query, model)
}

func LockBy(ctx context.Context, conn ksql.ConnectionInterface, model ksql.ModelInterface, call func(query ksql.QueryInterface)) error {
	if conn == nil || !conn.InTransaction() {
		return Err_Not_In_Transaction
	}

	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).ForUpdate()
	call(query)
	return QueryRowBy(ctx, conn, query, model)
}

func Table(ctx context.Context, table string, call func(table ksql.TableInterface)) error {
	ta := NewTable().Table(table)
	call(ta)
	return ta.Exec(ctx)
}

func Schema(ctx context.Context, schema string, call func(schema ksql.SchemaInterface)) error {
	sc := NewSchema().Schema(schema)
	call(sc)
	_, err := Exec(ctx, sc)
	return err
}

func DropTableBy(ctx context.Context, conn ksql.ConnectionInterface, table string) error {
	op := NewDropTable().Table(table)
	_, err := conn.Exec(ctx, op)
	return err
}

func DropTable(ctx context.Context, table string) error {
	return DropTableBy(ctx, database, table)
}
