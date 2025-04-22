package sharding

import (
	"context"
	"database/sql"
	"fmt"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
)

var database *Connection

func Init(configs []db.Config) error {
	conn := &Connection{baseConnection: &baseConnection{conns: make([]ksql.ConnectionInterface, len(configs))}, currents: make(map[any]ksql.ConnectionInterface)}
	for index, conf := range configs {
		dbConn, err := sql.Open(conf.DriverName, conf.DataSourceName)
		if err != nil {
			return err
		}

		dbConn.SetConnMaxIdleTime(conf.MaxIdleTime)
		dbConn.SetConnMaxLifetime(conf.MaxLifeTime)
		dbConn.SetMaxIdleConns(conf.MaxIdleConns)
		dbConn.SetMaxOpenConns(conf.MaxIdleConns)
		c, err := db.Open(dbConn, conf.DriverName)
		if err != nil {
			return err
		}

		conn.conns[index] = c
		conn.driverName = conf.DriverName
	}

	conn.count = len(conn.conns)
	database = conn
	return nil
}

func InitBy(driverName string, conns []*sql.DB) error {
	conn := &Connection{baseConnection: &baseConnection{conns: make([]ksql.ConnectionInterface, len(conns)), driverName: driverName}, currents: make(map[any]ksql.ConnectionInterface)}
	for index, co := range conns {
		c, err := db.Open(co, driverName)
		if err != nil {
			return err
		}

		conn.conns[index] = c
	}

	conn.count = len(conn.conns)
	database = conn
	return nil
}

func Database(key any) *sql.DB {
	return database.Database(key)
}

func Get() ConnectionInterface {
	return database.Clone()
}

func Close() error {
	if database == nil {
		return nil
	}

	return database.Close()
}

func Insert(key any, ctx context.Context, table string, data *db.Data) (int64, error) {
	return InsertBy(key, ctx, database, table, data)
}

func Update(key any, ctx context.Context, table string, data *db.Data, where ksql.WhereInterface) (int64, error) {
	return UpdateBy(key, ctx, database, table, data, where)
}

func Delete(key any, ctx context.Context, table string, where ksql.WhereInterface) (int64, error) {
	return DeleteBy(key, ctx, database, table, where)
}

func Exec(key any, ctx context.Context, op ksql.SqlInterface) (int64, error) {
	return ExecBy(key, ctx, database, op)
}

func InsertBy(key any, ctx context.Context, conn ConnectionInterface, table string, data *db.Data) (int64, error) {
	return db.InsertBy(ctx, conn.Get(key), table, data)
}

func UpdateBy(key any, ctx context.Context, conn ConnectionInterface, table string, data *db.Data, where ksql.WhereInterface) (int64, error) {
	return db.UpdateBy(ctx, conn.Get(key), table, data, where)
}

func DeleteBy(key any, ctx context.Context, conn ConnectionInterface, table string, where ksql.WhereInterface) (int64, error) {
	return db.DeleteBy(ctx, conn.Get(key), table, where)
}

func ExecBy(key any, ctx context.Context, conn ConnectionInterface, op ksql.SqlInterface) (int64, error) {
	return db.ExecBy(ctx, conn.Get(key), op)
}

func Query[T ksql.RowInterface](key any, ctx context.Context, op ksql.QueryInterface, models *[]T) error {
	return QueryBy(key, ctx, database, op, models)
}

func QueryBy[T ksql.RowInterface](key any, ctx context.Context, conn ConnectionInterface, op ksql.QueryInterface, models *[]T) error {
	if err := db.QueryBy(ctx, conn.Get(key), op, models); err != nil {
		return err
	}

	for _, model := range *models {
		var tmp any = model
		if t, ok := tmp.(ShardingInterface); ok {
			t.WithKey(key)
		}
	}

	return nil
}

func QueryRow[T ksql.RowInterface](key any, ctx context.Context, op ksql.QueryInterface, model T) error {
	return QueryRowBy(key, ctx, database, op, model)
}

func QueryRowBy[T ksql.RowInterface](key any, ctx context.Context, conn ConnectionInterface, op ksql.QueryInterface, model T) error {
	var tmp any = model
	if t, ok := tmp.(ShardingInterface); ok {
		t.WithKey(key)
	}

	if conn := model.Conn(); conn != nil {
		return db.QueryRowBy(ctx, conn, op, model)
	}

	return db.QueryRowBy(ctx, conn.Get(key), op, model)
}

func Find[T db.FindType](ctx context.Context, model ModelInterface, id T) error {
	return FindWith(ctx, database, model, id)
}

func FindWith[T db.FindType](ctx context.Context, conn ConnectionInterface, model ModelInterface, id T) error {
	query := db.NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).Where(model.PrimaryId(), "=", id)
	return db.QueryRowBy(ctx, conn.Get(model.Key()), query, model)
}

func FindBy(ctx context.Context, model ModelInterface, call func(query ksql.QueryInterface)) error {
	return FindByWith(ctx, database, model, call)
}

func FindByWith(ctx context.Context, conn ConnectionInterface, model ModelInterface, call func(query ksql.QueryInterface)) error {
	query := db.NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...)
	call(query)
	return db.QueryRowBy(ctx, conn.Get(model.Key()), query, model)
}

func Transaction(ctx context.Context, keys []any, call func(ctx context.Context, conn ConnectionInterface) error) ksql.TxError {
	return database.Clone().TransactionBy(ctx, keys, nil, call)
}

func TransactionBy(ctx context.Context, keys []any, options *sql.TxOptions, call func(ctx context.Context, conn ConnectionInterface) error) ksql.TxError {
	return database.Clone().TransactionBy(ctx, keys, options, call)
}

func Lock[T db.FindType](ctx context.Context, conn ConnectionInterface, model ModelInterface, id T) error {
	if conn == nil || !conn.InTransaction() {
		return fmt.Errorf("%s on key %v", db.Err_Not_In_Transaction, model.Key())
	}

	query := db.NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).Where(model.PrimaryId(), "=", id).For().Update()
	return db.QueryRowBy(ctx, conn.Get(model.Key()), query, model)
}

func LockBy(ctx context.Context, conn ConnectionInterface, model ModelInterface, call func(query ksql.QueryInterface)) error {
	if conn == nil || !conn.InTransaction() {
		return fmt.Errorf("%s on key %v", db.Err_Not_In_Transaction, model.Key())
	}

	query := db.NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).For().Update()
	call(query)
	return db.QueryRowBy(ctx, conn.Get(model.Key()), query, model)
}

func Table(ctx context.Context, table string, call func(table ksql.TableInterface)) error {
	return database.Range(func(index int, conn ksql.ConnectionInterface) error {
		ta := db.NewTable().Table(fmt.Sprintf("%s_%d", table, index)).WithConn(conn).IfNotExists()
		call(ta)
		if err := ta.Exec(ctx); err != nil {
			return fmt.Errorf("alter table: %s error: %s", fmt.Sprintf("%s_%d", table, index), err)
		}

		return nil
	})
}

func Schema(ctx context.Context, schema string, call func(schema ksql.SchemaInterface)) error {
	return database.Range(func(index int, conn ksql.ConnectionInterface) error {
		sc := db.NewSchema().Schema(schema).IfNotExists()
		call(sc)
		_, err := db.ExecBy(ctx, conn, sc)
		return err
	})
}

func DropTable(ctx context.Context, table string) error {
	return database.Range(func(index int, conn ksql.ConnectionInterface) error {
		return db.DropTableIfExistsBy(ctx, conn, table)
	})
}

func ShowDDL(ctx context.Context, table string) (string, error) {
	return db.ShowDDLBy(ctx, database.first(), table)
}
