package sharding

import (
	"context"
	"database/sql"
	"fmt"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
)

var baseConns []ksql.ConnectionInterface
var connsCount = 0

func Init(configs []db.Config) error {
	baseConns = make([]ksql.ConnectionInterface, len(configs))
	for index, conf := range configs {
		dbConn, err := sql.Open(conf.DriverName, conf.DataSourceName)
		if err != nil {
			return err
		}

		dbConn.SetConnMaxIdleTime(conf.MaxIdleTime)
		dbConn.SetConnMaxLifetime(conf.MaxLifeTime)
		dbConn.SetMaxIdleConns(conf.MaxIdleConns)
		dbConn.SetMaxOpenConns(conf.MaxIdleConns)
		conn, err := db.Open(dbConn, conf.DriverName)
		if err != nil {
			return err
		}

		baseConns[index] = conn
	}

	connsCount = len(baseConns)
	return nil
}

func Database(key any) *sql.DB {
	if connsCount == 0 {
		return nil
	}

	return baseConns[node(key, connsCount)].Database()
}

func Get(key any) (ksql.ConnectionInterface, error) {
	if connsCount == 0 {
		return nil, db.Err_Database_Not_Initialized
	}

	return baseConns[node(key, connsCount)].Clone(), nil
}

func Close() error {
	var err error
	for _, conn := range baseConns {
		err = conn.Database().Close()
	}

	return err
}

func _getConn(key any) ksql.ConnectionInterface {
	if connsCount == 0 {
		return nil
	}

	return baseConns[node(key, connsCount)]
}

func Insert(key any, ctx context.Context, table string, data *db.Data) (int64, error) {
	return db.InsertBy(ctx, _getConn(key), table, data)
}

func Update(key any, ctx context.Context, table string, data *db.Data, where ksql.WhereInterface) (int64, error) {
	return db.UpdateBy(ctx, _getConn(key), table, data, where)
}

func Delete(key any, ctx context.Context, table string, where ksql.WhereInterface) (int64, error) {
	return db.DeleteBy(ctx, _getConn(key), table, where)
}

func Exec(key any, ctx context.Context, op ksql.SqlInterface) (int64, error) {
	return db.ExecBy(ctx, _getConn(key), op)
}

func Query[T ksql.RowInterface](key any, ctx context.Context, op ksql.QueryInterface, models *[]T) error {
	if err := db.QueryBy(ctx, _getConn(key), op, models); err != nil {
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
	var tmp any = model
	if t, ok := tmp.(ShardingInterface); ok {
		t.WithKey(key)
	}

	if conn := model.Conn(); conn != nil {
		return db.QueryRowBy(ctx, conn, op, model)
	}

	return db.QueryRowBy(ctx, _getConn(key), op, model)
}

func Find[T db.FindType](ctx context.Context, model ModelInterface, id T) error {
	query := db.NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).Where(model.PrimaryId(), "=", id)
	return db.QueryRowBy(ctx, _getConn(model.Key()), query, model)
}

func FindBy(ctx context.Context, model ModelInterface, call func(query ksql.QueryInterface)) error {
	query := db.NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...)
	call(query)
	return db.QueryRowBy(ctx, _getConn(model.Key()), query, model)
}

func Lock[T db.FindType](ctx context.Context, model ModelInterface, id T) error {
	conn := _getConn(model.Key())
	if conn == nil || !conn.InTransaction() {
		return db.Err_Not_In_Transaction
	}

	query := db.NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).Where(model.PrimaryId(), "=", id).For().Update()
	return db.QueryRowBy(ctx, conn, query, model)
}

func LockBy(ctx context.Context, model ModelInterface, call func(query ksql.QueryInterface)) error {
	conn := _getConn(model.Key())
	if conn == nil || !conn.InTransaction() {
		return db.Err_Not_In_Transaction
	}

	query := db.NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).For().Update()
	call(query)
	return db.QueryRowBy(ctx, conn, query, model)
}

func Table(ctx context.Context, table string, call func(table ksql.TableInterface)) error {
	for index := 0; index < connsCount; index++ {
		ta := db.NewTable().Table(fmt.Sprintf("%s_%d", table, index)).WithConn(baseConns[index])
		call(ta)
		if err := ta.Exec(ctx); err != nil {
			return fmt.Errorf("alter table: %s error: %s", fmt.Sprintf("%s_%d", table, index), err)
		}
	}

	return nil
}

func Schema(ctx context.Context, schema string, call func(schema ksql.SchemaInterface)) error {
	for index := 0; index < connsCount; index++ {
		sc := db.NewSchema().Schema(schema)
		call(sc)
		if _, err := db.ExecBy(ctx, baseConns[index], sc); err != nil {
			return err
		}
	}

	return nil
}

func DropTable(ctx context.Context, table string) error {
	for index := 0; index < connsCount; index++ {
		if err := db.DropTableBy(ctx, baseConns[index], table); err != nil {
			return err
		}
	}

	return nil
}

func ShowDDL(ctx context.Context, table string) (string, error) {
	return db.ShowDDLBy(ctx, baseConns[0], table)
}
