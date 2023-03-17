package ck

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/kovey/config-go/config"
	"github.com/kovey/db-go/db"
	ds "github.com/kovey/db-go/sql"
	"github.com/kovey/debug-go/debug"
)

var (
	database *sql.DB
	dev      string
)

type ClickHouse[T any] struct {
	database        *sql.DB
	tx              *sql.Tx
	isInTransaction bool
}

func NewClickHouse[T any]() *ClickHouse[T] {
	return &ClickHouse[T]{database: database, tx: nil, isInTransaction: false}
}

func Init(conf config.ClickHouse) error {
	db, err := OpenDB(conf)
	if err != nil {
		return err
	}

	database = db
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
	configs = append(configs)
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
	if ck.isInTransaction {
		return ck.tx
	}

	return ck.database
}

func (ck *ClickHouse[T]) Begin() error {
	tx, err := ck.database.Begin()
	if err != nil {
		return err
	}

	ck.tx = tx
	ck.isInTransaction = true
	return nil
}

func (ck *ClickHouse[T]) Commit() error {
	if ck.tx == nil {
		return fmt.Errorf("transaction is not open or close")
	}

	ck.isInTransaction = false
	err := ck.tx.Commit()
	ck.tx = nil
	return err
}

func (ck *ClickHouse[T]) InTransaction() bool {
	return ck.isInTransaction
}

func (ck *ClickHouse[T]) Query(query string, t any, args ...any) ([]any, error) {
	return db.Query(ck.getDb(), query, t)
}

func (ck *ClickHouse[T]) Exec(statement string) error {
	return db.Exec(ck.getDb(), statement)
}

func (ck *ClickHouse[T]) Insert(insert *ds.Insert) (int64, error) {
	return 0, errors.New("insert statement supported only in the batch mode (use begin/commit)")
}

func (ck *ClickHouse[T]) Update(update *ds.Update) (int64, error) {
	db.Update(ck.getDb(), update)
	return 1, nil
}

func (ck *ClickHouse[T]) Delete(del *ds.Delete) (int64, error) {
	db.Delete(ck.getDb(), del)
	return 1, nil
}

func (ck *ClickHouse[T]) BatchInsert(batch *ds.Batch) (int64, error) {
	ins := batch.Inserts()
	count := int64(len(ins))
	if count == 0 {
		return count, errors.New("batch is empty")
	}

	err := ck.Begin()
	if err != nil {
		return 0, err
	}

	first := ins[0]
	smt, e := ck.getDb().Prepare(first.Prepare())
	if e != nil {
		ck.RollBack()
		return 0, e
	}

	for _, insert := range ins {
		insert.ParseValue(first.Fields())
		_, err = smt.Exec(insert.Args()...)
		if err != nil {
			debug.Erro("insert fail, error: %s", err)
		}
	}

	err = ck.Commit()

	return count, err
}

func (ck *ClickHouse[T]) Select(sel *ds.Select, model T) ([]T, error) {
	return db.Select(ck.getDb(), sel, model)
}

func (ck *ClickHouse[T]) FetchRow(table string, where map[string]any, model T) (any, error) {
	return db.FetchRow(ck.getDb(), table, where, model)
}

func (ck *ClickHouse[T]) FetchAll(table string, where map[string]any, model T) ([]T, error) {
	return db.FetchAll(ck.getDb(), table, where, model)
}

func (ck *ClickHouse[T]) FetchAllByWhere(table string, where *ds.Where, model T) ([]T, error) {
	return db.FetchAllByWhere(ck.getDb(), table, where, model)
}

func (ck *ClickHouse[T]) RollBack() error {
	ck.isInTransaction = false
	return nil
}

func (ck *ClickHouse[T]) FetchPage(table string, where map[string]any, model T, page int, pageSize int) ([]T, error) {
	return db.FetchPage(ck.getDb(), table, where, model, page, pageSize)
}

func (ck *ClickHouse[T]) FetchPageByWhere(table string, where *ds.Where, model T, page int, pageSize int) ([]T, error) {
	return db.FetchPageByWhere(ck.getDb(), table, where, model, page, pageSize)
}
