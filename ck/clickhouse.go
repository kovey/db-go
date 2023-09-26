package ck

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/db"
	"github.com/kovey/db-go/v2/itf"
	ds "github.com/kovey/db-go/v2/sql"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/debug-go/debug"
)

var (
	database *sql.DB
	dbName   string
)

type ClickHouse[T itf.ModelInterface] struct {
	database *sql.DB
	tx       *db.Tx
	DbName   string
}

func NewClickHouse[T itf.ModelInterface]() *ClickHouse[T] {
	return &ClickHouse[T]{database: database, tx: nil, DbName: dbName}
}

func Init(conf config.ClickHouse) error {
	db, err := OpenDB(conf)
	if err != nil {
		return err
	}

	database = db
	dbName = conf.Dbname
	return nil
}

func OpenDB(conf config.ClickHouse) (*sql.DB, error) {
	dsn := GetDSN(conf)
	debug.Info("clickhouse dsn: %s", dsn)
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			err = fmt.Errorf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}

		return nil, err
	}

	return db, nil
}

func GetDSN(conf config.ClickHouse) string {
	// tcp://host1:9000?username=user&password=qwerty&database=clicks&read_timeout=10&write_timeout=20&alt_hosts=host2:9000,host3:9000
	dsn := "tcp://%s?%s"
	host := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	configs := make([]string, 0)
	configs = append(configs, formatString("username", conf.Username))
	configs = append(configs, formatString("password", conf.Password))
	configs = append(configs, formatString("database", conf.Dbname))
	configs = append(configs, formatInt("read_timeout", conf.Timeout.Read))
	configs = append(configs, formatInt("write_timeout", conf.Timeout.Write))
	configs = append(configs, formatInt("write_timeout", conf.Timeout.Write))
	if conf.Cluster.Open == "On" {
		configs = append(configs, formatList("alt_hosts", conf.Cluster.Servers))
	}
	configs = append(configs, formatString("connection_open_strategy", conf.OpenStrategy))
	configs = append(configs, formatInt("block_size", conf.BlockSize))
	configs = append(configs, formatInt("pool_size", conf.PoolSize))
	configs = append(configs, formatBool("debug", conf.Debug))
	configs = append(configs, formatInt("compress", conf.Compress))

	return fmt.Sprintf(dsn, host, strings.Join(configs, "&"))
}

func formatBool(key string, value bool) string {
	return fmt.Sprintf("%s=%t", key, value)
}

func formatString(key string, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

func formatInt(key string, value int) string {
	return fmt.Sprintf("%s=%d", key, value)
}

func formatList(key string, servers []config.Addr) string {
	hosts := make([]string, len(servers))
	for index, server := range servers {
		hosts[index] = fmt.Sprintf("%s:%d", server.Host, server.Port)
	}

	return strings.Join(hosts, ",")
}

func (ck *ClickHouse[T]) getDb() db.ConnInterface {
	if ck.InTransaction() {
		return ck.tx.Tx()
	}

	return ck.database
}

func (ck *ClickHouse[T]) Transaction(func(*db.Tx) error) error {
	return fmt.Errorf("clickhouse unsupported transaction")
}

func (ck *ClickHouse[T]) TransactionCtx(context.Context, func(*db.Tx) error) error {
	return fmt.Errorf("clickhouse unsupported transaction")
}

func (ck *ClickHouse[T]) SetTx(tx *db.Tx) {
	ck.tx = tx
}

func (ck *ClickHouse[T]) begin() error {
	tx, err := ck.database.Begin()
	if err != nil {
		return err
	}

	ck.SetTx(db.NewTx(tx))
	return nil
}

func (ck *ClickHouse[T]) commit() error {
	if ck.tx == nil {
		return fmt.Errorf("transaction is not open or close")
	}

	err := ck.tx.Commit()
	ck.tx = nil
	return err
}

func (ck *ClickHouse[T]) InTransaction() bool {
	return ck.tx != nil && !ck.tx.IsCompleted()
}

func (ck *ClickHouse[T]) Query(query string, model T, args ...any) ([]T, error) {
	return db.Query(context.Background(), ck.getDb(), query, model)
}

func (ck *ClickHouse[T]) Exec(statement string) error {
	return db.Exec(context.Background(), ck.getDb(), statement)
}

func (ck *ClickHouse[T]) Insert(insert *ds.Insert) (int64, error) {
	return 0, errors.New("insert statement supported only in the batch mode (use begin/commit)")
}

func (ck *ClickHouse[T]) Update(update *ds.Update) (int64, error) {
	if _, err := db.Update(context.Background(), ck.getDb(), update); err != nil {
		return 0, err
	}

	return 1, nil
}

func (ck *ClickHouse[T]) Delete(del *ds.Delete) (int64, error) {
	if _, err := db.Delete(context.Background(), ck.getDb(), del); err != nil {
		return 0, err
	}

	return 1, nil
}

func (ck *ClickHouse[T]) BatchInsert(batch *ds.Batch) (int64, error) {
	ins := batch.Inserts()
	count := int64(len(ins))
	if count == 0 {
		return count, errors.New("batch is empty")
	}

	err := ck.begin()
	if err != nil {
		return 0, err
	}

	first := ins[0]
	smt, e := ck.getDb().Prepare(first.Prepare())
	if e != nil {
		if err := ck.RollBack(); err != nil {
			return 0, err
		}

		return 0, e
	}

	for _, insert := range ins {
		insert.ParseValue(first.Fields())
		_, err = smt.Exec(insert.Args()...)
		if err != nil {
			debug.Erro("insert fail, error: %s", err)
		}
	}

	err = ck.commit()

	return count, err
}

func (ck *ClickHouse[T]) Select(sel *ds.Select, model T) ([]T, error) {
	return db.Select(context.Background(), ck.getDb(), sel, model)
}

func (ck *ClickHouse[T]) Desc(desc *ds.Desc, model T) ([]T, error) {
	return db.Desc(context.Background(), ck.database, desc, model)
}

func (ck *ClickHouse[T]) ShowTables(show *ds.ShowTables, model T) ([]T, error) {
	return db.ShowTables(context.Background(), ck.database, show, model)
}

func (ck *ClickHouse[T]) FetchRow(table string, where meta.Where, model T) error {
	return db.FetchRow(context.Background(), ck.getDb(), table, where, model)
}

func (ck *ClickHouse[T]) LockRow(table string, where meta.Where, model T) error {
	return db.FetchRow(context.Background(), ck.getDb(), table, where, model)
}

func (ck *ClickHouse[T]) FetchAll(table string, where meta.Where, model T) ([]T, error) {
	return db.FetchAll(context.Background(), ck.getDb(), table, where, model)
}

func (ck *ClickHouse[T]) FetchAllByWhere(table string, where *ds.Where, model T) ([]T, error) {
	return db.FetchAllByWhere(context.Background(), ck.getDb(), table, where, model)
}

func (ck *ClickHouse[T]) RollBack() error {
	return nil
}

func (ck *ClickHouse[T]) FetchPage(table string, where meta.Where, model T, page int, pageSize int) (*meta.Page[T], error) {
	return db.FetchPage(context.Background(), ck.getDb(), table, where, model, page, pageSize)
}

func (ck *ClickHouse[T]) FetchPageByWhere(table string, where *ds.Where, model T, page int, pageSize int) (*meta.Page[T], error) {
	return db.FetchPageByWhere(context.Background(), ck.getDb(), table, where, model, page, pageSize)
}

func (ck *ClickHouse[T]) Count(table string, where *ds.Where) (int64, error) {
	return db.Count(context.Background(), ck.getDb(), table, where)
}

func (ck *ClickHouse[T]) FetchPageBySelect(sel *ds.Select, model T) (*meta.Page[T], error) {
	return db.FetchPageBySelect(context.Background(), ck.getDb(), sel, model)
}

func (ck *ClickHouse[T]) QueryCtx(ctx context.Context, query string, model T, args ...any) ([]T, error) {
	return db.Query(ctx, ck.getDb(), query, model)
}

func (ck *ClickHouse[T]) ExecCtx(ctx context.Context, statement string) error {
	return db.Exec(ctx, ck.getDb(), statement)
}

func (ck *ClickHouse[T]) InsertCtx(ctx context.Context, insert *ds.Insert) (int64, error) {
	return 0, errors.New("insert statement supported only in the batch mode (use begin/commit)")
}

func (ck *ClickHouse[T]) UpdateCtx(ctx context.Context, update *ds.Update) (int64, error) {
	if _, err := db.Update(ctx, ck.getDb(), update); err != nil {
		return 0, err
	}

	return 1, nil
}

func (ck *ClickHouse[T]) DeleteCtx(ctx context.Context, del *ds.Delete) (int64, error) {
	if _, err := db.Delete(ctx, ck.getDb(), del); err != nil {
		return 0, err
	}

	return 1, nil
}

func (ck *ClickHouse[T]) BatchInsertCtx(ctx context.Context, batch *ds.Batch) (int64, error) {
	ins := batch.Inserts()
	count := int64(len(ins))
	if count == 0 {
		return count, errors.New("batch is empty")
	}

	err := ck.begin()
	if err != nil {
		return 0, err
	}

	first := ins[0]
	smt, e := ck.getDb().PrepareContext(ctx, first.Prepare())
	if e != nil {
		if err := ck.RollBack(); err != nil {
			return 0, err
		}

		return 0, e
	}

	for _, insert := range ins {
		insert.ParseValue(first.Fields())
		_, err = smt.ExecContext(ctx, insert.Args()...)
		if err != nil {
			debug.Erro("insert fail, error: %s", err)
		}
	}

	err = ck.commit()

	return count, err
}

func (ck *ClickHouse[T]) SelectCtx(ctx context.Context, sel *ds.Select, model T) ([]T, error) {
	return db.Select(ctx, ck.getDb(), sel, model)
}

func (ck *ClickHouse[T]) DescCtx(ctx context.Context, desc *ds.Desc, model T) ([]T, error) {
	return db.Desc(ctx, ck.database, desc, model)
}

func (ck *ClickHouse[T]) ShowTablesCtx(ctx context.Context, show *ds.ShowTables, model T) ([]T, error) {
	return db.ShowTables(ctx, ck.database, show, model)
}

func (ck *ClickHouse[T]) FetchRowCtx(ctx context.Context, table string, where meta.Where, model T) error {
	return db.FetchRow(ctx, ck.getDb(), table, where, model)
}

func (ck *ClickHouse[T]) LockRowCtx(ctx context.Context, table string, where meta.Where, model T) error {
	return db.FetchRow(ctx, ck.getDb(), table, where, model)
}

func (ck *ClickHouse[T]) FetchAllCtx(ctx context.Context, table string, where meta.Where, model T) ([]T, error) {
	return db.FetchAll(ctx, ck.getDb(), table, where, model)
}

func (ck *ClickHouse[T]) FetchAllByWhereCtx(ctx context.Context, table string, where *ds.Where, model T) ([]T, error) {
	return db.FetchAllByWhere(ctx, ck.getDb(), table, where, model)
}

func (ck *ClickHouse[T]) FetchPageCtx(ctx context.Context, table string, where meta.Where, model T, page int, pageSize int) (*meta.Page[T], error) {
	return db.FetchPage(ctx, ck.getDb(), table, where, model, page, pageSize)
}

func (ck *ClickHouse[T]) FetchPageByWhereCtx(ctx context.Context, table string, where *ds.Where, model T, page int, pageSize int) (*meta.Page[T], error) {
	return db.FetchPageByWhere(ctx, ck.getDb(), table, where, model, page, pageSize)
}

func (ck *ClickHouse[T]) CountCtx(ctx context.Context, table string, where *ds.Where) (int64, error) {
	return db.Count(ctx, ck.getDb(), table, where)
}

func (ck *ClickHouse[T]) FetchPageBySelectCtx(ctx context.Context, sel *ds.Select, model T) (*meta.Page[T], error) {
	return db.FetchPageBySelect(ctx, ck.getDb(), sel, model)
}
